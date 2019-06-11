// Copyright 2019 Axetroy. All rights reserved. MIT license.
package accession

type Accession struct {
	Name        string `json:"name"`        // 权限标识符
	Description string `json:"description"` // 权限描述
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
