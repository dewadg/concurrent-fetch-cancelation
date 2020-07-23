package models

type Photo struct {
	ID           uint   `json:"id"`
	AlbumID      uint   `json:"albumId"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}
