package model

import "time"

type Item struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Category      string    `json:"category"`
	Chapter       string    `json:"chapter"`
	File          string    `json:"file"`
	FileType      string    `json:"fileType"`
	CreatedBy     string    `json:"createdBy"`
	CreatedByName string    `json:"createdByName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
