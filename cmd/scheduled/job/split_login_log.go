package job

import (
	"fmt"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

func generateNewLoginTableName(date time.Time) string {
	year := fmt.Sprintf("%d", date.Year())
	month := fmt.Sprintf("%d", date.Month())
	if len(month) == 1 {
		month = "0" + month
	}

	loginLog := model.LoginLog{}

	tableName := loginLog.TableName() + "_" + year + month

	return tableName
}

// 迁移数据
func moveLoginLog(startAt time.Time, endAt time.Time) (bool, error) {
	var (
		tx  = database.Db.Begin()
		err error
	)

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	newTableName := generateNewLoginTableName(startAt)

	if err = ensureLoginLogTableExist(newTableName, tx); err != nil {
		return true, err
	}

	logs := make([]model.LoginLog, 0)

	limit := 100 // 一次迁移一百条数据

	// 查找在这个时间段的数据，移动到新表中
	// 最早的数据排到前面
	if err = tx.Model(model.LoginLog{}).Where("created_at >= ?", startAt).Where("created_at <= ?", endAt).Limit(limit).Order("created_at ASC").Find(&logs).Error; err != nil {
		return true, err
	}

	for _, loginLog := range logs {
		dataID := loginLog.Id
		if err = tx.Table(newTableName).Create(&loginLog).Error; err != nil {
			return true, err
		} else {
			// 更新表信息 - 还原 ID/创建时间/更新时间 信息
			if err = tx.Table(newTableName).Where("id = ?", loginLog.Id).UpdateColumn("id", dataID).UpdateColumn("created_at", loginLog.CreatedAt).UpdateColumn("updated_at", loginLog.UpdatedAt).Error; err != nil {
				return true, err
			}

			// 删除旧数据
			if err = tx.Unscoped().Table(loginLog.TableName()).Delete(model.LoginLog{Id: dataID}).Error; err != nil {
				return true, err
			}
		}
	}

	// 如果获取的数据已经不够，那么我们就认为它已经是最后一页了
	eol := len(logs) < limit

	return eol, nil
}

func ensureLoginLogTableExist(tableName string, db *gorm.DB) error {
	// 如果表不存在，那么创建表
	if !db.HasTable(tableName) {
		if err := db.Table(tableName).CreateTable(model.LoginLog{}).Error; err != nil {
			return err
		}
	}

	return nil
}

// 定时切割用户登录记录
// 因为这个表的内容是在是太大了
func SplitLoginLog() {
	var err error
	now := time.Now()

	oldestLoginLog := model.LoginLog{}

	if err = database.Db.Model(oldestLoginLog).Order("created_at ASC").First(&oldestLoginLog).Error; err != nil {
		return
	}

	startAt := time.Date(oldestLoginLog.CreatedAt.Year(), oldestLoginLog.CreatedAt.Month(), 1, 0, 0, 0, 0, now.Location())
	endAt := startAt.AddDate(0, 1, 0)

	// 如果最旧的数据是在本月或者上个月产生的，那么跳过任务
	if endAt.After(now.AddDate(0, -1, 0)) {
		return
	}

	for {
		if eol, err := moveLoginLog(startAt, endAt); err != nil {
			return
		} else if eol {
			// 开始时间往后推一个月，继续遍历
			startAt = startAt.AddDate(0, 1, 0)
			endAt = startAt.AddDate(0, 1, 0)

			// 如果最旧的数据是在本月或者上个月产生的，那么跳过任务
			if endAt.After(now.AddDate(0, -1, 0)) {
				break
			}
		}
	}
}
