// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/user_server/controller/report"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
	)
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	r := report.Create(helper.Context{Uid: userInfo.Id}, report.CreateParams{
		Title:   title,
		Content: content,
		Type:    reportType,
	})

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	reportInfo := schema.Report{}

	assert.Nil(t, r.Decode(&reportInfo))

	defer report.DeleteReportById(reportInfo.Id)

	assert.Equal(t, title, reportInfo.Title)
	assert.Equal(t, content, reportInfo.Content)
	assert.Equal(t, reportType, reportInfo.Type)
}

func TestCreateRouter(t *testing.T) {
	var (
		title      = "title"
		content    = "content"
		reportType = model.ReportTypeBug
	)
	userInfo, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	body, _ := json.Marshal(&report.CreateParams{
		Title:   title,
		Content: content,
		Type:    reportType,
	})

	res := tester.HttpUser.Post("/v1/report", body, &header)
	r := schema.Response{}

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Nil(t, json.Unmarshal([]byte(res.Body.String()), &r))

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	reportInfo := schema.Report{}

	assert.Nil(t, r.Decode(&reportInfo))

	defer report.DeleteReportById(reportInfo.Id)

	assert.Equal(t, title, reportInfo.Title)
	assert.Equal(t, content, reportInfo.Content)
	assert.Equal(t, reportType, reportInfo.Type)

}
