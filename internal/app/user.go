package app

import (
	"errors"
	"github.com/axetroy/go-server/core/controller/wallet"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/axetroy/go-server/core/util"
	"github.com/axetroy/go-server/core/validator"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

var (
	User *UserService
)

type Model struct {
}

type SignUpWithUsernameParams struct {
	Username   string  `json:"username" valid:"required~请输入用户名"` // 用户名
	Password   string  `json:"password" valid:"required~请输入密码"`  // 密码
	InviteCode *string `json:"invite_code"`                      // 邀请码
}

type UserService struct {
}

func New() *UserService {
	return &UserService{}
}

func (u *UserService) createUserTx(tx *gorm.DB, userInfo *model.User, inviterCode *string) (err error) {
	var (
		newTx bool
	)
	if tx == nil {
		tx = database.Db.Begin()
		newTx = true
	}

	defer func() {
		if newTx {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}
	}()

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

func (u *UserService) SignUp(params SignUpWithUsernameParams) (data *schema.Profile, err error) {
	var (
		tx *gorm.DB
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
	}()

	// 参数校验
	if err = validator.ValidateStruct(params); err != nil {
		return
	}

	if err = validator.ValidateUsername(params.Username); err != nil {
		return
	}

	tx = database.Db.Begin()

	u1 := model.User{Username: params.Username}

	if err = tx.Where("username = ?", params.Username).Find(&u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	}

	if u1.Id != "" {
		err = exception.UserExist
		return
	}

	userInfo := model.User{
		Username: params.Username,
		Nickname: &params.Username,
		Password: util.GeneratePassword(params.Password),
		Status:   model.UserStatusInit,
		Role:     pq.StringArray{model.DefaultUser.Name},
		Phone:    nil,
		Email:    nil,
		Gender:   model.GenderUnknown,
	}

	if err = u.createUserTx(tx, &userInfo, params.InviteCode); err != nil {
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

func (u *UserService) GetDetail() {

}
