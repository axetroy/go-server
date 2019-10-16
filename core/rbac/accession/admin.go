// Copyright 2019 Axetroy. All rights reserved. MIT license.
package accession

var (
	AdminAdminGet    = New("admin::get", "有权限获取横幅信息")
	AdminAdminCreate = New("admin::create", "有权限创建新横幅")
	AdminAdminUpdate = New("admin::update", "有权限修改横幅信息")
	AdminAdminDelete = New("admin::delete", "有权限删除横幅")

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

	AdminBannerGet    = New("banner::get", "有权限获取横幅信息")
	AdminBannerCreate = New("banner::create", "有权限创建新横幅")
	AdminBannerUpdate = New("banner::update", "有权限修改横幅信息")
	AdminBannerDelete = New("banner::delete", "有权限删除横幅")

	AdminReportGet    = New("report::get", "有权限获取反馈信息")
	AdminReportUpdate = New("report::update", "有权限修改反馈信息")
	AdminReportDelete = New("report::delete", "有权限删除反馈信息")

	// 管理员的所有权限
	AdminList = []*Accession{
		AdminAdminGet,
		AdminAdminCreate,
		AdminAdminUpdate,
		AdminAdminDelete,

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

		AdminBannerGet,
		AdminBannerCreate,
		AdminBannerUpdate,
		AdminBannerDelete,

		AdminReportGet,
		AdminReportUpdate,
		AdminReportDelete,
	}

	AdminMap = map[string]*Accession{}
)

func init() {
	for _, a := range AdminList {
		AdminMap[a.Name] = a
	}
}
