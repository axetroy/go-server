package address

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateAddressParams struct {
	Name         string `json:"name" valid:"required~请填写收货人"`         // 收货人
	Phone        string `json:"phone" valid:"required~请输入收货人电话号码"`    // 收货人手机号
	ProvinceCode string `json:"province_code" valid:"required~请选择省份"` // 省份代码
	CityCode     string `json:"city_code" valid:"required~请选择城市"`     // 城市代码
	AreaCode     string `json:"area_code" valid:"required~请选择区域"`     // 区域代码
	Address      string `json:"address" valid:"required~请输入详细地址"`     // 详细的地址
	IsDefault    *bool  `json:"is_default"`                           // 是否是默认地址
}

func Create(context controller.Context, input CreateAddressParams) (res schema.Response) {
	var (
		err          error
		data         schema.Address
		tx           *gorm.DB
		isDefault    bool = false
		isValidInput bool
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = service.Db.Begin()

	userInfo := model.User{
		Id: context.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if input.IsDefault != nil {
		isDefault = *input.IsDefault

		defaultAddress := model.Address{
			Uid:       context.Uid,
			IsDefault: true,
		}
		if err = tx.First(&defaultAddress).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			} else {
				return
			}
		} else {
			// 如果存在了默认地址，则取消它的默认属性
			if err = tx.Model(&defaultAddress).UpdateColumn(model.Address{
				IsDefault: false,
			}).Error; err != nil {
				return
			}
		}

	} else {
		firstAddress := model.Address{
			Uid: context.Uid,
		}
		if err = tx.Where(&firstAddress).First(&firstAddress).Error; err != nil {
			// 如果还没有设置过地址，那么这次设置就是默认地址
			if err == gorm.ErrRecordNotFound {
				err = nil
				isDefault = true
			} else {
				return
			}
		}
	}

	AddressInfo := model.Address{
		Uid:          context.Uid,
		Name:         input.Name,
		Phone:        input.Phone,
		ProvinceCode: input.ProvinceCode,
		CityCode:     input.CityCode,
		AreaCode:     input.AreaCode,
		Address:      input.Address,
		IsDefault:    isDefault,
	}

	if err = tx.Create(&AddressInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(AddressInfo, &data.AddressPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = AddressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = AddressInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func CreateRouter(context *gin.Context) {
	var (
		input CreateAddressParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
