// Copyright 2019 Axetroy. All rights reserved. MIT license.
package accession

type Accession struct {
	Name        string `json:"name"`        // 权限标识符
	Description string `json:"description"` // 权限描述
}

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

	// 管理员类
	AdminNewsGet    = New("news::get", "有权限获取新闻")
	AdminNewsCreate = New("news::create", "有权限创建新闻")
	AdminNewsUpdate = New("news::update", "有权限修改新闻")
	AdminNewsDelete = New("news::delete", "有权限删除新闻")

	adminNotificationGet    = New("notification::get", "有权限获取公告")
	adminNotificationCreate = New("notification::create", "有权限创建公告")
	adminNotificationUpdate = New("notification::update", "有权限修改公告")
	adminNotificationDelete = New("notification::delete", "有权限删除公告")

	AdminUserGet    = New("user::get", "有权限获取用户信息")
	AdminUserCreate = New("user::create", "有权限创建新用户")
	AdminUserUpdate = New("user::update", "有权限修改用户信息")
	AdminUserDelete = New("user::delete", "有权限删除用户")
	AdminUserExport = New("user::export", "有权限导出用户到CSV等")

	AdminMenuGet    = New("menu::get", "有权限获取菜单信息")
	AdminMenuCreate = New("menu::create", "有权限创建新菜单")
	AdminMenuUpdate = New("menu::update", "有权限修改菜单信息")
	AdminMenuDelete = New("menu::delete", "有权限删除菜单")

	// 管理员的所有权限
	AdminList = []*Accession{
		AdminNewsGet,
		AdminNewsCreate,
		AdminNewsUpdate,
		AdminNewsDelete,

		adminNotificationGet,
		adminNotificationUpdate,
		adminNotificationDelete,
		adminNotificationCreate,

		AdminUserGet,
		AdminUserCreate,
		AdminUserUpdate,
		AdminUserDelete,
		AdminUserExport,

		AdminMenuGet,
		AdminMenuCreate,
		AdminMenuUpdate,
		AdminMenuDelete,
	}

	AdminMap = map[string]*Accession{}
)

func init() {
	for _, a := range List {
		Map[a.Name] = a
	}
	for _, a := range AdminList {
		AdminMap[a.Name] = a
	}
}

// 校验一个权限是否是合法的字符串
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

// 筛选出有效的管理员权限
func FilterAdminAccession(AccessionStr []string) (accession []string) {
	for _, v := range AccessionStr {
		if _, ok := AdminMap[v]; ok == true {
			accession = append(accession, v)
		}
	}
	if accession == nil {
		accession = []string{}
	}
	return
}
