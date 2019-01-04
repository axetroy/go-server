package news

import (
	"github.com/axetroy/go-server/model"
)

type Pure struct {
	Id      string           `json:"id"`
	Author  string           `json:"author"`
	Tittle  string           `json:"tittle"`
	Content string           `json:"content"`
	Type    model.NewsType   `json:"type"`
	Tags    []string         `json:"tags"`
	Status  model.NewsStatus `json:"status"`
}

type News struct {
	Pure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
