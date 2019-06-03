package exception

var (
	RoleNotExist     = New("角色不存在")
	RoleCannotUpdate = New("无法更新角色")
	RoleHadBeenUsed  = New("角色正在被使用，无法删除")
)
