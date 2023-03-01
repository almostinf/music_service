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
	songNode := &entity.SongNode{SongID: 0, NextSongID: 0, PrevSongID: 0}

	// Check if SongNode's table exist
	if err := r.DB.Model(&entity.SongNode{}).Migrator().HasTable(songNode); !err {
		if err := r.DB.Debug().Model(&entity.SongNode{}).Migrator().CreateTable(songNode); err != nil {
			return playlist, err
		}
	}

	// Create Head and Tail SongNodes
	if err := r.DB.Model(&entity.SongNode{}).Create(&playlist.HeadSong).Error; err != nil {
		return playlist, err
	}

	if err := r.DB.Model(&entity.SongNode{}).Create(&playlist.TailSong).Error; err != nil {
		return playlist, err
	}
	
	if err := r.DB.Migrator().HasTable(playlist); !err {
		if err := r.DB.Debug().Migrator().CreateTable(playlist); err != nil {
			return playlist, err
		}
	}

	if err := r.DB.First(playlist, "id = ?", playlist.ID).Error; err != nil {
		if err := r.DB.Create(playlist).Error; err != nil {
			return playlist, err
		}

		// Find user by creator's id
		var user entity.User
		if err := r.DB.Model(&entity.User{}).First(&user, playlist.CreatorID).Error; err != nil {
			return playlist, err
		}

		// Append new playlist and save it in database
		user.Playlists = append(user.Playlists, *playlist)
		if err := r.DB.Model(&entity.User{Model: gorm.Model{ID: playlist.CreatorID}}).Save(user).Error; err != nil {
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
		if playlist.ID == playlist_id {
			return true
		}
	}
	return false
}

func findSongInPlaylist(playlist *entity.Playlist, song_id uint) bool {
	curNode := playlist.HeadSong
	for curNode != playlist.TailSong && curNode.ID != 0 {
		if curNode.ID == song_id {
			return true
		}
		curNode = entity.SongNode{Model: gorm.Model{ID: curNode.NextSongID}}
	}
	return curNode.ID == song_id
}

func (r *Playlist) getEntitiesByID(user_id uint, song_id uint, playlist_id uint) (*entity.User, *entity.Song, *entity.Playlist, error) {
	// Retrieve the playlist from the database
	playlist := &entity.Playlist{}
	if err := r.Model(&entity.Playlist{}).First(playlist, playlist_id).Error; err != nil {
		return nil, nil, nil, errors.New("failed to retrieve playlist from database")
	}

	// Check if the song is present in the playlist
	if !findSongInPlaylist(playlist, song_id) {
		return nil, nil, nil, errors.New("playlist doesn't contain this song")
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
	user, song, _, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}

	// Update the current song information for the user
	user.CurSong.CurPlayingSongID = song_id
	user.CurSong.IsPlaying = true
	user.CurSong.PlaylistID = playlist_id
	if err := r.Save(user).Error; err != nil {
		return errors.New("failed to update user in database")
	}

	// Start playing the song in a separate goroutine
	go func() {
		time.Sleep(song.Duration)
		user.CurSong.IsPlaying = false
		if err := r.Save(user).Error; err != nil {
			fmt.Printf("failed to update user in database: %v", err)
		}
	}()

	return nil
}

func (r *Playlist) Pause(user_id uint, song_id uint, playlist_id uint) error {
	user, _, _, err := r.getEntitiesByID(user_id, song_id, playlist_id)
	if err != nil {
		return err
	}
	user.CurSong.IsPlaying = false
	if err := r.Save(user).Error; err != nil {
		return errors.New("failed to update user in database")
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

	// Create new SongNode in database
	if err := r.DB.Model(&entity.SongNode{}).Create(songNode).Error; err != nil {
		return err
	}

	// Adding new Node to the tail if it's not nil
	if playlist.TailSong.ID != 0 {
		playlist.TailSong.NextSongID = songNode.ID

		// Update PrevSondID field in new SongNode
		songNode.PrevSongID = playlist.TailSong.ID
	}
	
	playlist.TailSong = *songNode

	// If the playlist is empty, set the new SongNode as both the HeadSong and the TailSong
	if playlist.HeadSong.ID == 0 {
		playlist.HeadSong = *songNode
	}

	// Save new SongNode
	if err := r.DB.Model(&entity.SongNode{Model: gorm.Model{ID: songNode.ID}}).Save(songNode).Error; err != nil {
		return err
	}

	// Update playlist in database
	if err := r.DB.Model(&entity.Playlist{Model: gorm.Model{ID: playlist_id}}).Save(playlist).Error; err != nil {
		return err
	}

	return nil
}

func (r *Playlist) findNextSongInPlaylist(playlist *entity.Playlist, song_id uint) (uint, error) {
	curSongNode := playlist.HeadSong
	for curSongNode != playlist.TailSong && curSongNode.ID != 0 {
		if curSongNode.ID == song_id {
			if curSongNode.NextSongID != 0 {
				return curSongNode.NextSongID, nil
			}
			return 0, errors.New("can't find next song in playlist")
		}
		if err := r.DB.Model(&entity.SongNode{}).First(curSongNode, curSongNode.NextSongID).Error; err != nil {
			return 0, errors.New("can't find next song in playlist")
		}
		curSongNode = entity.SongNode{Model: gorm.Model{ID: curSongNode.NextSongID}}
	}
	if curSongNode.ID != 0 && curSongNode.ID == song_id {
		if curSongNode.NextSongID != 0 {
			return curSongNode.NextSongID, nil
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
	curSongNode := playlist.HeadSong
	for curSongNode != playlist.TailSong && curSongNode.ID != 0 {
		if curSongNode.ID == song_id {
			if curSongNode.PrevSongID != 0 {
				return curSongNode.PrevSongID, nil
			}
			return 0, errors.New("can't find previous song in playlist")
		}
		if err := r.DB.Model(&entity.SongNode{}).First(curSongNode, curSongNode.PrevSongID).Error; err != nil {
			return 0, errors.New("can't find previous song in playlist")
		}
		curSongNode = entity.SongNode{Model: gorm.Model{ID: curSongNode.PrevSongID}}
	}
	if curSongNode.ID != 0 && curSongNode.ID == song_id {
		if curSongNode.PrevSongID != 0 {
			return curSongNode.PrevSongID, nil
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

	prev_song_id, err := r.findPrevSongInPlaylist(playlist, song_id)
	if err != nil {
		return err
	}

	if err = r.Play(user_id, prev_song_id, playlist_id); err != nil {
		return err
	}

	return nil
}
