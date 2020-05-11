package schema

import (
	"github.com/axetroy/go-server/internal/model"
)

type HelpPure struct {
	Id       string           `json:"id"`        // 帮助文章ID
	Title    string           `json:"title"`     // 帮助文章标题
	Content  string           `json:"content"`   // 帮助文章内容
	Tags     []string         `json:"tags"`      // 帮助文章的标签
	Status   model.HelpStatus `json:"status"`    // 帮助文章状态
	Type     model.HelpType   `json:"type"`      // 帮助文章的类型
	ParentId *string          `json:"parent_id"` // 父级 ID，如果有的话
}

type Help struct {
	HelpPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
