// Copyright 2019 Axetroy. All rights reserved. MIT license.
package accession

var (
	// 用户类
	ProfileUpdate   = New("profile::update", "有权限修改用户资料")
	PasswordUpdate  = New("password::update", "有权限更改自己的密码")
	Password2Set    = New("password2.set", "有权限设置二级密码")
	Password2Reset  = New("password2.reset", "有权限重置二级密码")
	Password2Update = New("password2::update", "有权限修改二级密码")
	DoTransfer      = New("transfer::create", "有权限发起转账交易")

	// 用户的所有的权限
	List = []*Accession{
		ProfileUpdate,
		PasswordUpdate,
		Password2Set,
		Password2Update,
		DoTransfer,
	}

	Map = map[string]*Accession{}
)

func init() {
	for _, a := range List {
		Map[a.Name] = a
	}
}
