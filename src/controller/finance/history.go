// Copyright 2019 Axetroy. All rights reserved. MIT license.
package finance

import (
	"github.com/gin-gonic/gin"
)

func GetHistory(context *gin.Context) {
	//var (
	//	err     error
	//	session *xorm.Session
	//	tx      bool
	//	data    []Log
	//	meta    = response.Meta{}
	//	query   = request.Query{}
	//)
	//
	//defer func() {
	//	if r := recover(); r != nil {
	//		switch t := r.(type) {
	//		case string:
	//			err = errors.New(t)
	//		case error:
	//			err = t
	//		default:
	//			err = exception.Unknown
	//		}
	//	}
	//
	//	if tx {
	//		if err != nil {
	//			_ = session.Rollback()
	//		} else {
	//			err = session.Commit()
	//		}
	//	}
	//
	//	if session != nil {
	//		session.Close()
	//	}
	//
	//	if err != nil {
	//		context.JSON(http.StatusOK, response.Response{
	//			Status:  response.StatusFail,
	//			Message: err.Error(),
	//			Data:    nil,
	//		})
	//	} else {
	//		context.JSON(http.StatusOK, response.List{
	//			Response: response.Response{
	//				Status:  response.StatusSuccess,
	//				Message: "",
	//				Data:    data,
	//			},
	//			Meta: &meta,
	//		})
	//	}
	//}()
	//
	//uid := context.GetInt64("uid")
	//
	//if err = context.ShouldBindQuery(&query); err != nil {
	//	return
	//}
	//
	//query.Normalize()
	//
	//session = orm.Db.NewSession()
	//
	//if err = session.Begin(); err != nil {
	//	return
	//}
	//
	//tx = true
	//
	//keys := mapStr(structs.Names(&model.FinanceLog{}), func(s string) string {
	//	return strcase.ToSnake(s)
	//})
	//
	//seleted := strings.Join(keys[:], ",")
	//
	//sql := GenerateSql(uid, seleted)
	//
	//session = session.SQL(sql + fmt.Sprintf(" LIMIT %v", query.Limit))
	//
	//if res, er := session.QueryInterface(); er != nil {
	//	err = er
	//	return
	//} else {
	//	var (
	//		length = len(res)
	//		total  int64
	//	)
	//
	//	meta.Num = length
	//	meta.Page = query.Page
	//	meta.Limit = query.Limit
	//	meta.Sort = query.Sort
	//	meta.Platform = query.Platform
	//
	//	// 如果查出来的总数
	//	if length >= query.Limit {
	//		// 统计总数
	//		countSql := GenerateSql(uid, "COUNT(*)")
	//
	//		if total, err = session.SQL(countSql).Count(); err != nil {
	//			return
	//		}
	//
	//		meta.Total = total
	//	} else {
	//		meta.Total = int64(length)
	//	}
	//
	//	data = make([]Log, 0)
	//
	//	for _, v := range res {
	//		log := Log{}
	//		if err = mapstructure.Decode(v, &log); err != nil {
	//			return
	//		}
	//
	//		createdAt := v["created_at"].(time.Time)
	//		updatedAt := v["updated_at"].(time.Time)
	//
	//		// balance
	//		log.BeforeBalance = v["before_balance"].(float64)
	//		log.BalanceMutation = v["balance_mutation"].(float64)
	//		log.AfterBalance = v["after_balance"].(float64)
	//
	//		// frozen
	//		log.BeforeFrozen = v["before_frozen"].(float64)
	//		log.FrozenMutation = v["frozen_mutation"].(float64)
	//		log.AfterFrozen = v["after_frozen"].(float64)
	//
	//		log.CreatedAt = createdAt.Format(time.RFC3339Nano)
	//		log.UpdatedAt = updatedAt.Format(time.RFC3339Nano)
	//		data = append(data, log)
	//	}
	//}
}

//func GenerateSql(UidField int64, selected string) string {
//	suffix := `("deleted_at" IS NULL OR "deleted_at"='0001-01-01 00:00:00')`
//
//	tables := []string{"finance_log_cny", "finance_log_usd", "finance_log_coin"}
//
//	sqlList := make([]string, 0)
//
//	if selected == "*" {
//		selected = "uid"
//	}
//
//	for _, tableName := range tables {
//		sql := fmt.Sprintf(`SELECT %v FROM "%v" WHERE "uid"=%v AND %v`, selected, tableName, UidField, suffix)
//		sqlList = append(sqlList, sql)
//	}
//
//	sql := strings.Join(sqlList[:], " UNION ")
//	return sql
//}
//
//func mapStr(vs []string, f func(string) string) []string {
//	vsm := make([]string, len(vs))
//	for i, v := range vs {
//		vsm[i] = f(v)
//	}
//	return vsm
//}
