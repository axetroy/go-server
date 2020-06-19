// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var now = time.Now()

var dates = []time.Time{
	time.Date(2000, 6, 1, 12, 0, 0, 0, now.Location()),
	time.Date(2000, 7, 2, 13, 0, 0, 0, now.Location()),
	time.Date(2000, 8, 3, 14, 0, 0, 0, now.Location()),
	time.Date(2000, 9, 4, 15, 0, 0, 0, now.Location()),
	time.Date(2000, 10, 5, 16, 0, 0, 0, now.Location()),
}

func createTestData() (err error) {
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	times := dates

	userInfo, _ := tester.CreateUser()
	waiterInfo, _ := tester.CreateWaiter()

	// 创建多个会话
	for index, t := range times {
		session := model.CustomerSession{
			Id:        util.MD5(fmt.Sprintf("%s%s%d", userInfo.Id, waiterInfo.Id, index)),
			Uid:       userInfo.Id,
			WaiterID:  waiterInfo.Id,
			CreatedAt: t,
			UpdatedAt: t,
		}

		if err = tx.Create(&session).Error; err != nil {
			return err
		}

		// 这个会话下创建多条

		i := 0

	children:
		for {
			if i > 3 {
				break children
			}

			b, _ := json.Marshal(ws.MessageTextPayload{
				Text: fmt.Sprintf("Hello world %d", i),
			})

			sessionItem := model.CustomerSessionItem{
				SessionID:  session.Id,
				SenderID:   userInfo.Id,
				ReceiverID: waiterInfo.Id,
				Payload:    string(b),
				CreatedAt:  t,
				UpdatedAt:  t,
			}

			if err = tx.Create(&sessionItem).Error; err != nil {
				return err
			}

			i++
		}
	}

	return
}

func TestCustomer_Next(t *testing.T) {
	assert.Nil(t, createTestData())

	c := Customer{}
	gotShouldGoNext, err := c.Next(time.Now())

	assert.Nil(t, err)
	assert.False(t, gotShouldGoNext)

	session := model.CustomerSession{}

	for _, d := range dates {
		tableName := generateTableName(session.TableName(), d)

		hasExist := database.Db.HasTable(tableName)

		assert.True(t, hasExist)

		{
			sessions := make([]model.CustomerSession, 0)
			// 确保这些表下的数据都符合这个时间范围
			assert.Nil(t, database.Db.Table(tableName).Limit(100).Find(&sessions).Error)

			// 查得到数据
			assert.True(t, len(sessions) > 0)

			for _, session := range sessions {
				assert.Equal(t, d.Year(), session.CreatedAt.Year())
				assert.Equal(t, d.Month(), session.CreatedAt.Month())
			}
		}

		{
			sessionItem := model.CustomerSessionItem{}
			itemTableName := generateTableName(sessionItem.TableName(), d)
			sessionItems := make([]model.CustomerSessionItem, 0)
			// 确保这些表下的数据都符合这个时间范围
			assert.Nil(t, database.Db.Table(itemTableName).Limit(100).Find(&sessionItems).Error)

			// 查得到数据
			assert.True(t, len(sessionItems) > 0)

			for _, sessionItem := range sessionItems {
				assert.Equal(t, d.Year(), sessionItem.CreatedAt.Year())
				assert.Equal(t, d.Month(), sessionItem.CreatedAt.Month())
			}
		}
	}

}
