package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/almostinf/music_service/internal/entity"
	"gorm.io/gorm"
)

type Playlist struct {
	*gorm.DB
}

func NewPlaylist(db *gorm.DB) *Playlist {
	return &Playlist{db}
}

func (r *Playlist) Get() []entity.Playlist {
	var playlists []entity.Playlist
	r.Find(&playlists)
	return playlists
}

func (r *Playlist) Create(playlist *entity.Playlist) (*entity.Playlist, error) {
	if err := r.DB.Migrator().HasTable(playlist); !err {
		if err := r.DB.Debug().Migrator().CreateTable(playlist); err != nil {
			return playlist, err
		}
	}

	if err := r.DB.First(playlist, "id = ?", playlist.ID).Error; err != nil {
		if err := r.DB.Model(&entity.Playlist{}).Create(playlist).Error; err != nil {
			return playlist, err
		}

		// Find user by creator's id
		var user entity.User
		if err := r.DB.Model(&entity.User{}).First(&user, playlist.CreatorID).Error; err != nil {
			return playlist, err
		}

		// Append new playlist and save it in database
		user.Playlists = append(user.Playlists, int64(playlist.ID))
		if err := r.DB.Model(&entity.User{Model: gorm.Model{ID: playlist.CreatorID}}).Save(&user).Error; err != nil {
			return playlist, err
		}

		return playlist, nil
	} else {
		return playlist, err
	}
}

func (r *Playlist) Update(id uint, playlist *entity.Playlist) (*entity.Playlist, error) {
	var finded_playlist entity.Playlist
	if err := r.DB.First(&finded_playlist, "id = ?", id).Error; err != nil {
		return playlist, err
	}
	if err := r.DB.Save(playlist).Error; err != nil {
		return playlist, err
	}
	return playlist, nil
}

func (r *Playlist) Delete(id uint) error {
	var playlist entity.Playlist
	if err := r.DB.First(&playlist, "id = ?", id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&playlist).Error; err != nil {
		return err
	}
	return nil
}

func findPlaylist(user *entity.User, playlist_id uint) bool {
	for _, playlist := range user.Playlists {
		if playlist == int64(playlist_id) {
			return true
		}
	}
	return false
}

func (r *Playlist) findSongInPlaylist(playlist *entity.Playlist, song_id uint) bool {
	var curNode entity.SongNode
	if err := r.DB.Model(&entity.SongNode{}).First(&curNode, playlist.HeadSong).Error; err != nil {
		return false
	}
	for curNode.ID != playlist.TailSong && curNode.ID != 0 {
		if curNode.SongID == song_id {
			return true
		}
		var nextSongNode entity.SongNode
		if err := r.DB.Model(&entity.SongNode{}).First(&nextSongNode, curNode.NextSongID).Error; err != nil {
			return false
		}
		curNode = nextSongNode
	}
	return curNode.SongID == song_id
}

func (r *Playlist) getEntitiesByID(user_id uint, song_id uint, playlist_id uint) (*entity.User, *entity.Song, *entity.Playlist, error) {
	// Retrieve the playlist from the database
	playlist := &entity.Playlist{}
	if err := r.Model(&entity.Playlist{}).First(playlist, playlist_id).Error; err != nil {
		return nil, nil, nil, errors.New("failed to retrieve playlist from database")
	}

	// Retrieve the user from the database
	user := &entity.User{}
	if err := r.Model(&entity.User{}).First(user, user_id).Error; err != nil {
		return nil, nil, nil, errors.New("failed to retrieve user from database")
	}

	// Check if the user has access to the playlist
	if !findPlaylist(user, playlist_id) {
		return nil, nil, nil, errors.New("user doesn't have that playlist")
	}

	// Retrieve the song from the database
	song := &entity.Song{}
	if err := r.Model(&entity.Song{}).First(song, song_id).Error; err != nil {
		return nil, nil, nil, errors.New("failed to retrieve song from database")
	}

	return user, song, playlist, nil
}

func (r *Playlist) Play(user_id uint, song_id uint, playlist_id uint) error {
	user, song, playlist, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// Check if the song is present in the playlist
	if !r.findSongInPlaylist(playlist, song_id) {
		return errors.New("playlist doesn't contain this song")
	}

	var curSong entity.CurSongInfo
	if err := r.DB.Model(&entity.CurSongInfo{}).Find(&curSong, user.CurSongInfoID).Error; err != nil {
		return err
	}

	// Update the current song information for the user
	curSong.CurPlayingSongID = song_id
	curSong.IsPlaying = true
	curSong.PlaylistID = playlist_id

	if err := r.DB.Model(&entity.CurSongInfo{}).Where("id = ?", user.CurSongInfoID).Updates(&curSong).Error; err != nil {
		return err
	}

	// Start playing the song in a separate goroutine
	go func() {
		time.Sleep(song.Duration)
		if err := r.DB.Model(&entity.CurSongInfo{}).Where("id = ?", user.CurSongInfoID).Select("is_playing").Updates(map[string]interface{}{"is_playing": false}).Error; err != nil {
			fmt.Printf("Can't update songInfo: %v\n", err)
		}
	}()

	return nil
}

func (r *Playlist) Pause(user_id uint, song_id uint, playlist_id uint) error {
	user, _, playlist, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// Check if the song is present in the playlist
	if !r.findSongInPlaylist(playlist, song_id) {
		return errors.New("playlist doesn't contain this song")
	}

	var curSong entity.CurSongInfo
	if err := r.DB.Model(&entity.CurSongInfo{}).Find(&curSong, user.CurSongInfoID).Error; err != nil {
		return err
	}

	if err := r.DB.Model(&entity.CurSongInfo{}).Where("id = ?", user.CurSongInfoID).Select("is_playing").Updates(map[string]interface{}{"is_playing": false}).Error; err != nil {
		return err
	}

	return nil
}

func (r *Playlist) AddSong(user_id uint, song_id uint, playlist_id uint) error {
	_, _, playlist, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// create a new SongNode with the given song_id
	songNode := &entity.SongNode{SongID: song_id, NextSongID: 0, PrevSongID: 0}

	// Check if SongNode's table exist
	if err := r.DB.Model(&entity.SongNode{}).Migrator().HasTable(songNode); !err {
		if err := r.DB.Debug().Model(&entity.SongNode{}).Migrator().CreateTable(songNode); err != nil {
			return err
		}
	}

	if err := r.DB.Model(&entity.SongNode{}).Where("song_id = ?", song_id).First(songNode).Error; err != nil {
		if err := r.DB.Model(&entity.SongNode{}).Create(songNode).Error; err != nil {
			return err
		}
	}

	// Adding new Node to the tail if it's not nil
	if playlist.TailSong != 0 {
		var tailsong entity.SongNode
		if err := r.DB.Model(&entity.SongNode{}).First(&tailsong, playlist.TailSong).Error; err != nil {
			return err
		}

		tailsong.NextSongID = songNode.ID
		if err := r.DB.Model(&entity.SongNode{}).Where("id = ?", playlist.TailSong).Updates(&tailsong).Error; err != nil {
			return err
		}

		// Update PrevSongID field in new SongNode
		songNode.PrevSongID = playlist.TailSong
	}

	playlist.TailSong = songNode.ID

	// If the playlist is empty, set the new SongNode as both the HeadSong and the TailSong
	if playlist.HeadSong == 0 {
		playlist.HeadSong = songNode.ID
	}

	// Save new SongNode
	if err := r.DB.Model(&entity.SongNode{}).Where("id = ?", songNode.ID).Updates(&songNode).Error; err != nil {
		return err
	}

	// Update playlist in database
	if err := r.DB.Model(&entity.Playlist{}).Where("id = ?", playlist.ID).Updates(playlist).Error; err != nil {
		return err
	}

	return nil
}

func (r *Playlist) findNextSongInPlaylist(playlist *entity.Playlist, song_id uint) (uint, error) {
	var curSongNode entity.SongNode
	if err := r.DB.Model(&entity.SongNode{}).First(&curSongNode, playlist.HeadSong).Error; err != nil {
		return 0, err
	}

	for curSongNode.ID != playlist.TailSong && curSongNode.ID != 0 {
		if curSongNode.SongID == song_id {
			if curSongNode.NextSongID != 0 {
				var nextSongNode entity.SongNode
				if err := r.DB.Model(&entity.SongNode{}).First(&nextSongNode, curSongNode.NextSongID).Error; err != nil {
					return 0, errors.New("can't find next song in playlist")
				}
				return nextSongNode.SongID, nil
			}
			return 0, errors.New("can't find next song in playlist")
		}
		var nextSongNode entity.SongNode
		if err := r.DB.Model(&entity.SongNode{}).First(&nextSongNode, curSongNode.NextSongID).Error; err != nil {
			return 0, errors.New("can't find next song in playlist")
		}
		curSongNode = nextSongNode
	}

	if curSongNode.ID != 0 && curSongNode.SongID == song_id {
		if curSongNode.NextSongID != 0 {
			var nextSongNode entity.SongNode
			if err := r.DB.Model(&entity.SongNode{}).First(&nextSongNode, curSongNode.NextSongID).Error; err != nil {
				return 0, errors.New("can't find next song in playlist")
			}
			return nextSongNode.SongID, nil
		}
		return 0, errors.New("can't find next song in playlist")
	}
	return 0, errors.New("can't find given song in playlist")
}

func (r *Playlist) Next(user_id uint, song_id uint, playlist_id uint) error {
	_, _, playlist, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// Check if the song is present in the playlist
	if !r.findSongInPlaylist(playlist, song_id) {
		return errors.New("playlist doesn't contain this song")
	}

	next_song_id, err := r.findNextSongInPlaylist(playlist, song_id)
	if err != nil {
		return err
	}

	if err = r.Play(user_id, next_song_id, playlist_id); err != nil {
		return err
	}

	return nil
}

func (r *Playlist) findPrevSongInPlaylist(playlist *entity.Playlist, song_id uint) (uint, error) {
	var curSongNode entity.SongNode
	if err := r.DB.Model(&entity.SongNode{}).First(&curSongNode, playlist.HeadSong).Error; err != nil {
		return 0, err
	}
	for curSongNode.ID != playlist.TailSong && curSongNode.ID != 0 {
		if curSongNode.SongID == song_id {
			if curSongNode.PrevSongID != 0 {
				var prevSong entity.SongNode
				if err := r.DB.Model(&entity.SongNode{}).First(&prevSong, curSongNode.PrevSongID).Error; err != nil {
					return 0, errors.New("can't find next song in playlist")
				}
				return prevSong.SongID, nil
			}
			return 0, errors.New("can't find previous song in playlist")
		}
		var nextSongNode entity.SongNode
		if err := r.DB.Model(&entity.SongNode{}).First(&nextSongNode, curSongNode.NextSongID).Error; err != nil {
			return 0, errors.New("can't find next song in playlist")
		}
		curSongNode = nextSongNode
	}
	if curSongNode.ID != 0 && curSongNode.SongID == song_id {
		if curSongNode.PrevSongID != 0 {
			var prevSong entity.SongNode
			if err := r.DB.Model(&entity.SongNode{}).First(&prevSong, curSongNode.PrevSongID).Error; err != nil {
				return 0, errors.New("can't find next song in playlist")
			}
			return prevSong.SongID, nil
		}
		return 0, errors.New("can't find previous song in playlist")
	}
	return 0, errors.New("can't find given song in playlist")
}

func (r *Playlist) Prev(user_id uint, song_id uint, playlist_id uint) error {
	_, _, playlist, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// Check if the song is present in the playlist
	if !r.findSongInPlaylist(playlist, song_id) {
		return errors.New("playlist doesn't contain this song")
	}

	prev_song_id, err := r.findPrevSongInPlaylist(playlist, song_id)
	if err != nil {
		return err
	}

	if err = r.Play(user_id, prev_song_id, playlist_id); err != nil {
		return err
	}

	return nil
}
