// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package schema

import "github.com/axetroy/go-server/internal/model"

type NewsPure struct {
	Id      string           `json:"id"`
	Author  string           `json:"author"`
	Title   string           `json:"title"`
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
