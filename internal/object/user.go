package object

import (
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type User struct {
	tx  *gorm.DB
	ctx helper.Context
}

func NewUser(db *gorm.DB, ctx helper.Context) User {
	return User{
		tx:  db,
		ctx: ctx,
	}
}

type ParamsCreateWithUserName struct {
	Username   string  `json:"username" validate:"required,max=32" comment:"用户名"`    // 用户名
	Password   string  `json:"password" validate:"required,max=32" comment:"密码"`     // 密码
	InviteCode *string `json:"invite_code" validate:"omitempty,len=8" comment:"邀请码"` // 邀请码
}

// 内部方法
func (u User) create(userInfo *model.User, inviterCode *string) (err error) {
	var tx = u.tx
	if err = tx.Create(userInfo).Error; err != nil {
		return err
	}

	if inviterCode != nil && len(*inviterCode) > 0 {

		inviter := model.User{
			InviteCode: *inviterCode,
		}

		if err := tx.Where(&inviter).Find(&inviter).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.InvalidInviteCode
			}
			return err
		}

		// 如果存在邀请者的话，写入邀请列表中
		if inviter.Id != "" {
			inviteHistory := model.InviteHistory{
				Inviter:       inviter.Id,
				Invitee:       userInfo.Id,
				Status:        model.StatusInviteRegistered,
				RewardSettled: false,
			}

			// 创建邀请记录
			if err = tx.Create(&inviteHistory).Error; err != nil {
				return err
			}
		}
	}

	// 创建用户对应的钱包账号
	for _, walletName := range model.Wallets {
		if err = tx.Table(wallet.GetTableName(walletName)).Create(&model.Wallet{
			Id:       userInfo.Id,
			Currency: walletName,
			Balance:  0,
			Frozen:   0,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (u User) Create(input ParamsCreateWithUserName) (data schema.Profile, err error) {
	var tx = u.tx

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if err = validator.ValidateUsername(input.Username); err != nil {
		return
	}

	// 检查用户是否存在
	{
		userInfo := model.User{Username: input.Username}

		if err = tx.Model(userInfo).Where("username = ?", input.Username).Find(&userInfo).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				return
			}
		}

		if userInfo.Id != "" {
			err = exception.UserExist
			return
		}
	}

	userInfo := model.User{
		Username: input.Username,
		Nickname: &input.Username,
		Password: util.GeneratePassword(input.Password),
		Status:   model.UserStatusInit,
		Role:     pq.StringArray{model.DefaultUser.Name},
		Phone:    nil,
		Email:    nil,
		Gender:   model.GenderUnknown,
	}

	if err = u.create(&userInfo, input.InviteCode); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

type ParamsUpdateProfile struct {
	Username *string                    `json:"username" validate:"omitempty,max=32" comment:"用户名"` // 用户名，部分用户有机会修改自己的用户名，比如微信注册的帐号
	Nickname *string                    `json:"nickname" validate:"omitempty,max=32" comment:"昵称"`
	Gender   *model.Gender              `json:"gender" validate:"omitempty,number,oneof=0 1 2" comment:"性别"`
	Avatar   *string                    `json:"avatar" validate:"omitempty,url,max=255" comment:"头像"`
	Wechat   *UpdateWechatProfileParams `json:"wechat" validate:"omitempty" comment:"微信绑定信息"` // 更新微信绑定的帐号相关
}

// 绑定的微信信息帐号相关
type UpdateWechatProfileParams struct {
	Nickname  *string `json:"nickname" validate:"omitempty,max=32" comment:"微信昵称"`        // 用户昵称
	AvatarUrl *string `json:"avatar_url" validate:"omitempty,url,max=255" comment:"微信头像"` // 用户头像
	Gender    *int    `json:"gender" validate:"omitempty,number" comment:"性别"`            // 性别
	Country   *string `json:"country" validate:"omitempty,max=32" comment:"国家"`           // 国家
	Province  *string `json:"province" validate:"omitempty,max=32" comment:"省份"`          // 省份
	City      *string `json:"city" validate:"omitempty,max=32" comment:"城市"`              // 城市
	Language  *string `json:"language" validate:"omitempty,max=32" comment:"语言"`          // 语言
}

func (u User) Update(input ParamsUpdateProfile) (data schema.Profile, err error) {
	var (
		tx           = u.tx
		shouldUpdate bool
	)

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	updated := model.User{}

	if input.Username != nil {
		shouldUpdate = true

		if err = validator.ValidateUsername(*input.Username); err != nil {
			return
		}

		u := model.User{Id: u.ctx.Uid}

		if err = tx.Where(&u).First(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.UserNotExist
			}
			return
		}

		// 如果没有剩余的重命名次数的话
		if u.UsernameRenameRemaining <= 0 {
			err = exception.RenameUserNameFail
			return
		}

		updated.Username = *input.Username
		updated.UsernameRenameRemaining = u.UsernameRenameRemaining - 1
	}

	if input.Nickname != nil {
		updated.Nickname = input.Nickname
		shouldUpdate = true
	}

	if input.Avatar != nil {
		updated.Avatar = *input.Avatar
		shouldUpdate = true
	}

	if input.Gender != nil {
		updated.Gender = *input.Gender
		shouldUpdate = true
	}

	if shouldUpdate {
		if err = tx.Table(updated.TableName()).Where(model.User{Id: u.ctx.Uid}).Updates(updated).Error; err != nil {
			return
		}
	}

	userInfo := model.User{
		Id: u.ctx.Uid,
	}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if input.Wechat != nil {
		wechatInfo := model.WechatOpenID{
			Uid: userInfo.Id,
		}
		// 判断该用户是否绑定了微信帐号
		if err = tx.Where(&wechatInfo).First(&wechatInfo).Error; err != nil {
			// 如果没有找到，说明帐号没有绑定微信，抛出异常
			if err == gorm.ErrRecordNotFound {
				err = exception.InvalidParams
			}
			return
		}

		// 更新对应的字段
		wechatUpdated := model.WechatOpenID{}
		shouldUpdateWechat := false

		if input.Wechat.Nickname != nil {
			wechatUpdated.Nickname = input.Wechat.Nickname
			shouldUpdateWechat = true
		}

		if input.Wechat.AvatarUrl != nil {
			wechatUpdated.AvatarUrl = input.Wechat.AvatarUrl
			shouldUpdateWechat = true
		}

		if input.Wechat.Gender != nil {
			wechatUpdated.Gender = input.Wechat.Gender
			shouldUpdateWechat = true
		}

		if input.Wechat.Country != nil {
			wechatUpdated.Country = input.Wechat.Country
			shouldUpdateWechat = true
		}

		if input.Wechat.Province != nil {
			wechatUpdated.Province = input.Wechat.Province
			shouldUpdateWechat = true
		}

		if input.Wechat.City != nil {
			wechatUpdated.City = input.Wechat.City
			shouldUpdateWechat = true
		}

		if input.Wechat.Language != nil {
			wechatUpdated.Language = input.Wechat.Language
			shouldUpdateWechat = true
		}

		if shouldUpdateWechat {
			info := model.WechatOpenID{Id: wechatInfo.Id}
			if err = tx.Where(&info).Updates(wechatUpdated).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = exception.InvalidParams
				}
				return
			}

			wechat := schema.WechatBindingInfo{}

			if err = mapstructure.Decode(info, &wechat); err != nil {
				return
			}

			data.Wechat = &wechat
		}
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (u User) GetProfile(uid string) {

}
