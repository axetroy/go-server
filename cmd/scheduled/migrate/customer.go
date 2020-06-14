// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package migrate

import (
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

type Customer struct {
}

func (c Customer) GetTableName() string {
	m := model.CustomerSession{}

	return m.TableName()
}

func (c Customer) GetModel() interface{} {
	return model.CustomerSession{}
}

func (c Customer) GetTimeInterval(_ time.Time) time.Duration {
	return time.Hour * 24 * 7
}

func (c Customer) getLatest(db *gorm.DB) (model.CustomerSession, error) {
	oldestData := model.CustomerSession{}

	if err := db.Table(c.GetTableName()).Order("created_at ASC").First(&oldestData).Error; err != nil {
		return oldestData, err
	}

	return oldestData, nil
}

func (c Customer) Next(nows ...time.Time) (shouldGoNext bool, err error) {
	var (
		tx  = database.Db.Begin()
		now = time.Now()
	)

	if len(nows) > 0 {
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

	dataList := make([]model.CustomerSession, 0)

	limit := 100 // 一次迁移一百条数据

	// 查找在这个时间段的数据，移动到新表中
	// 最早的数据排到前面
	if err = tx.Model(model.CustomerSession{}).Where("created_at >= ?", startAt).Where("created_at <= ?", endAt).Limit(limit).Order("created_at ASC").Find(&dataList).Error; err != nil {
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

			// 移动它下面的聊天记录
			{
				sessionItems := make([]model.CustomerSessionItem, 0)
				if err = tx.Model(model.CustomerSessionItem{}).Where("session_id = ?", data.Id).Find(&sessionItems).Error; err != nil {
					return false, err
				}

				for _, sessionItem := range sessionItems {
					sessionItemId := sessionItem.Id

					newItemTableName := generateTableName(sessionItem.TableName(), data.CreatedAt)

					// 如果表不存在，那么创建表
					if !tx.HasTable(newItemTableName) {
						if err = tx.Table(newItemTableName).CreateTable(model.CustomerSessionItem{}).Error; err != nil {
							return false, err
						}
					}

					// 在新的表中创建
					if err = tx.Table(newItemTableName).Create(&sessionItem).Error; err != nil {
						return false, err
					} else {
						// 更新表信息 - 还原 ID/创建时间/更新时间 信息
						if err = tx.Table(newItemTableName).Where("id = ?", sessionItem.Id).Updates(map[string]interface{}{
							"id":         sessionItemId,
							"created_at": sessionItem.CreatedAt,
							"updated_at": sessionItem.UpdatedAt,
						}).Error; err != nil {
							return false, err
						}

						// 删除旧数据
						if err = tx.Unscoped().Delete(model.CustomerSessionItem{Id: sessionItemId}).Error; err != nil {
							return false, err
						}
					}
				}
			}

			// 删除旧数据
			if err = tx.Unscoped().Table(c.GetTableName()).Delete(model.CustomerSession{Id: dataID}).Error; err != nil {
				return false, err
			}
		}
	}

	// 如果获取的数据已经不够，那么我们就认为它已经是最后一页了
	shouldGoNext = len(dataList) > limit

	return shouldGoNext, nil
}

func (c Customer) Do() error {
	for {
		if shouldGoNext, err := c.Next(); err != nil {
			return err
		} else if !shouldGoNext {
			return nil
		}
	}
}
