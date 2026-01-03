package model

import "time"

type ShortURL struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	OriginalURL string    `json:"original_url" bson:"original_url"`
	ShortCode   string    `json:"short_code" bson:"short_code"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}
