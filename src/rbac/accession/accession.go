package accession

type Accession struct {
	Name        string `json:"name"`        // 权限标识符
	Description string `json:"description"` // 权限描述
}

var (
	// 用户类
	ProfileUpdate  = New("profile.update", "有权限修改用户资料")
	PasswordUpdate = New("password.update", "有权限更改自己的密码")
	DoTransfer     = New("transfer.create", "有权限发起转账交易")

	// 所有的权限
	List = []*Accession{
		ProfileUpdate,
		PasswordUpdate,
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
