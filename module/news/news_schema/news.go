// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news_schema

import (
	"github.com/axetroy/go-server/module/news/news_model"
)

type NewsPure struct {
	Id      string                `json:"id"`
	Author  string                `json:"author"`
	Title   string                `json:"title"`
	Content string                `json:"content"`
	Type    news_model.NewsType   `json:"type"`
	Tags    []string              `json:"tags"`
	Status  news_model.NewsStatus `json:"status"`
}

type News struct {
	NewsPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
