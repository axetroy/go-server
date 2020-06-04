// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package accession

var (
	// 用户类
	CustomerServiceConnect = New("customer_service::connect", "有权限连接客服")

	// 用户的所有的权限
	//WaiterList = []*Accession{
	//	CustomerServiceConnect,
	//}

	WaiterMap = map[string]*Accession{}
)

func init() {
	for _, a := range List {
		WaiterMap[a.Name] = a
	}
}
