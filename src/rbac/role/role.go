package role

import "github.com/axetroy/go-server/src/rbac/accession"

type Role struct {
	Name        string                `json:"name"`        // 角色名
	Description string                `json:"description"` // 角色描述
	Accession   []accession.Accession `json:"accession"`   // 角色拥有的权限
}

func New(name string, description string, accessions []accession.Accession) *Role {
	return &Role{
		Name:        name,
		Description: description,
		Accession:   accessions,
	}
}

func (r *Role) AccessionArray() (list []string) {
	for _, v := range r.Accession {
		list = append(list, v.Name)
	}
	return
}
