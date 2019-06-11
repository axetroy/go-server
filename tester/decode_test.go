package tester_test

import (
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	source := news_model.News{
		Title:   "title",
		Content: "content",
	}
	dest := news_schema.News{}
	assert.Nil(t, tester.Decode(source, &dest))

	assert.Equal(t, "title", dest.Title)
	assert.Equal(t, "content", dest.Content)

	assert.NotNil(t, tester.Decode(source, dest), "decode: dest expect a point")
}
