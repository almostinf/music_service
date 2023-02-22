package models

import (
	"container/list"
	"time"
)

type Playlist struct {
	ID        int64         `json:"id"`
	Title     string        `json:"title"`
	Creator   int64         `json:"creator"`
	Songs     *list.List    `json:"songs"`
	CurSong   *list.Element `json:"current_song"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
