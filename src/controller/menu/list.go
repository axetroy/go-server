// Copyright 2019 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	//Platform *model.BannerPlatform `json:"platform"` // 根据平台筛选
	//Active   *bool                 `json:"active"`   // 是否激活
}

func findChild(allMenus []schema.Menu, parentId string, myAccession []string, isSuperAdmin bool) (children []schema.Menu) {
	for _, v := range allMenus {
		if v.ParentId == parentId {
			if isSuperAdmin == false && len(v.Accession) > 0 {
				if matchAccession(v.Accession, myAccession) == true {
					children = append(children, v)
				}
			} else {
				children = append(children, v)
			}
		}
	}
	return
}

// 把子集的菜单挂载到父级菜单上
func set(menus []schema.Menu, allMenus []schema.Menu, myAccession []string, isSuperAdmin bool) (result []schema.Menu) {
	for _, v := range menus {
		isRequireAccession := len(v.Accession) > 0

		if !isRequireAccession {
			v.Accession = []string{}
		}
		// 获取该父级菜单下的所有子菜单
		children := findChild(allMenus, v.Id, myAccession, isSuperAdmin)
		for _, c := range children {
			v.Children = append(v.Children, c)
		}
		// 再查找子菜单的子菜单
		if len(v.Children) > 0 {
			v.Children = set(v.Children, allMenus, myAccession, isSuperAdmin)
		} else {
			v.Children = []schema.Menu{}
		}
		result = append(result, v)
	}
	return
}

func findStrInSlice(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func matchAccession(accessionRequire []string, myAccession []string) bool {
	for _, a := range accessionRequire {
		if findStrInSlice(myAccession, a) == false {
			return false
		}
	}
	return true
}

func transform(menus []schema.Menu, myAccession []string, isSuperAdmin bool) (data []schema.Menu) {
	// 提取一级目录
	for _, v := range menus {
		// 如果是子菜单, 则跳过这个
		if v.ParentId != "" {
			continue
		}

		isRequireAccession := len(v.Accession) > 0

		if !isRequireAccession {
			v.Accession = []string{}
		}

		// 如果这个菜单需要权限验证的话
		if isSuperAdmin == false && isRequireAccession {
			if matchAccession(v.Accession, myAccession) == true {
				data = append(data, v)
			}
		} else {
			data = append(data, v)
		}
	}

	data = set(data, menus, myAccession, isSuperAdmin)

	return
}

func GetList(context controller.Context, input Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.Menu, 0)
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			if len(data) > 0 {
				res.Data = data
			} else {
				res.Data = []string{}
			}
			res.Status = schema.StatusSuccess
		}
	}()

	adminInfo := model.Admin{
		Id: context.Uid,
	}

	if err = database.Db.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	list := make([]model.Menu, 0)

	// 获取所有菜单
	if err = database.Db.Order("sort desc").Order("name").Find(&list).Error; err != nil {
		return
	}

	// 解构
	for _, v := range list {
		d := schema.Menu{}
		if er := mapstructure.Decode(v, &d.MenuPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	data = transform(data, adminInfo.Accession, adminInfo.IsSuper)

	return
}

func GetListRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetList(controller.Context{
		Uid: context.GetString(middleware.ContextUidField),
	}, input)
}
