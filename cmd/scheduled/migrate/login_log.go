// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

type LoginLog struct {
}

func (c LoginLog) GetTableName() string {
	m := model.LoginLog{}

	return m.TableName()
}

func (c LoginLog) GetModel() interface{} {
	return model.LoginLog{}
}

func (c LoginLog) GetTimeInterval(now time.Time) time.Duration {
	currentDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	currentDate = currentDate.AddDate(0, 1, 0)

	// 在下一个月的第一天，减去 24 小时，那么就是当月的最后一天
	currentDate = currentDate.Add(-time.Hour * 24)

	return time.Hour * 24 * time.Duration(currentDate.Day())
}

func (c LoginLog) getLatest(db *gorm.DB) (model.LoginLog, error) {
	oldestData := model.LoginLog{}

	if err := db.Table(c.GetTableName()).Order("created_at ASC").First(&oldestData).Error; err != nil {
		return oldestData, err
	}

	return oldestData, nil
}

func (c LoginLog) Next(nows ...time.Time) (shouldGoNext bool, err error) {
	var (
		tx  = database.Db.Begin()
		now = time.Now()
	)

	if len(nows) == 0 {
		now = nows[0]
	}

	oldestData, err := c.getLatest(tx)

	if err != nil {
		return false, err
	}

	internal := c.GetTimeInterval(now)

	startAt := time.Date(oldestData.CreatedAt.Year(), oldestData.CreatedAt.Month(), 1, 0, 0, 0, 0, now.Location())
	endAt := now.Add(-internal)

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	dataList := make([]model.LoginLog, 0)

	limit := 100 // 一次迁移一百条数据

	// 查找在这个时间段的数据，移动到新表中
	// 最早的数据排到前面
	if err = tx.Model(model.LoginLog{}).Where("created_at >= ?", startAt).Where("created_at <= ?", endAt).Limit(limit).Order("created_at ASC").Find(&dataList).Error; err != nil {
		return false, err
	}

	for _, data := range dataList {
		dataID := data.Id

		newTableName := generateTableName(c.GetTableName(), data.CreatedAt)

		// 如果表不存在，那么创建表
		if !tx.HasTable(newTableName) {
			if err := tx.Table(newTableName).CreateTable(c.GetModel()).Error; err != nil {
				return false, err
			}
		}

		if err = tx.Table(newTableName).Create(&data).Error; err != nil {
			return false, err
		} else {
			// 更新表信息 - 还原 ID/创建时间/更新时间 信息
			if err = tx.Table(newTableName).Where("id = ?", data.Id).UpdateColumn("id", dataID).UpdateColumn("created_at", data.CreatedAt).UpdateColumn("updated_at", data.UpdatedAt).Error; err != nil {
				return false, err
			}

			// 删除旧数据
			if err = tx.Unscoped().Table(c.GetTableName()).Delete(model.LoginLog{Id: dataID}).Error; err != nil {
				return false, err
			}
		}
	}

	// 如果获取的数据已经不够，那么我们就认为它已经是最后一页了
	shouldGoNext = len(dataList) > limit

	return shouldGoNext, nil
}

func (c LoginLog) Do() error {
	for {
		if shouldGoNext, err := c.Next(); err != nil {
			return err
		} else if !shouldGoNext {
			return nil
		}
	}
}
