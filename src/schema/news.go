package schema

import "github.com/axetroy/go-server/src/model"

type NewsPure struct {
	Id      string           `json:"id"`
	Author  string           `json:"author"`
	Tittle  string           `json:"tittle"`
	Content string           `json:"content"`
	Type    model.NewsType   `json:"type"`
	Tags    []string         `json:"tags"`
	Status  model.NewsStatus `json:"status"`
}

type News struct {
	NewsPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
