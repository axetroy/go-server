// Copyright 2019 Axetroy. All rights reserved. MIT license.
package accession

type Accession struct {
	Name        string `json:"name"`        // 权限标识符
	Description string `json:"description"` // 权限描述
}

var (
	// 用户类
	ProfileUpdate   = New("profile.update", "有权限修改用户资料")
	PasswordUpdate  = New("password.update", "有权限更改自己的密码")
	Password2Set    = New("password2.set", "有权限设置二级密码")
	Password2Reset  = New("password2.reset", "有权限重置二级密码")
	Password2Update = New("password2.update", "有权限修改二级密码")
	DoTransfer      = New("transfer.create", "有权限发起转账交易")

	// 所有的权限
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

func Valid(s []string) bool {
	for _, v := range s {
		if _, ok := Map[v]; ok == false {
			return false
		}
	}
	return true
}

// 把权限转化成字符串
func Stringify(a ...*Accession) (list []string) {
	for _, v := range a {
		list = append(list, v.Name)
	}
	return
}

// 把权限字符串转化成权限模型
func Normalize(AccessionStr []string) (list []Accession) {
	for _, v := range AccessionStr {
		list = append(list, *New(v, ""))
	}
	return
}

// 生成一个新的实例
func New(name string, description string) *Accession {
	return &Accession{
		Name:        name,
		Description: description,
	}
}
