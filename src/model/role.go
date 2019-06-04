package model

import (
	"github.com/axetroy/go-server/src/rbac/accession"
	"github.com/axetroy/go-server/src/rbac/role"
	"github.com/lib/pq"
	"time"
)

var (
	DefaultUser = role.New("user", "普通用户", []accession.Accession{
		*accession.ProfileUpdate,
		*accession.Password2Set,
		*accession.Password2Reset,
		*accession.Password2Update,
		*accession.PasswordUpdate,
		*accession.DoTransfer,
	})
)

type Role struct {
	Name        string         `gorm:"primary_key;unique;not null;index;type:varchar(64)" json:"name"` // 角色名, 作为主建而且唯一
	Description string         `gorm:"not null;index;type:varchar(64)" json:"description"`             // 角色描述
	Accession   pq.StringArray `gorm:"not null;index;type:varchar(64)[]" json:"accession"`             // 改角色拥有的权限
	BuildIn     bool           `gorm:"not null;index;" json:"build_in"`                                // 是否是内建的角色，该角色通常是不可改的
	Note        *string        `gorm:"null;index;type:varchar(64)" json:"note"`                        // 备注
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (news *Role) TableName() string {
	return "role"
}
