# godijan--dijan客户端

## 简介
dijan的客户端sdk，自动进行集群一致性哈希适配，容器外部调用hostname映射  
dijan: ```https://github.com/shoogoome/dijan```

## 使用说明

**前置条件: 系统配置docker环境**

使用示例
```
package test

import (
	"encoding/json"
	"fmt"
	"github.com/shoogoome/godijan"
	"testing"
)

func TestConn(t *testing.T) {

	account := map[string]interface{} {
		"name": "121",
		"role": 99,
		"motto": "个性签名",
	}
	
	accountString, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	
	dijan := godijan.NewGoDijanConnection("localhost:2375", nil)
	if err := dijan.Set("models:account:1", accountString); err != nil {
		panic(err)
	}
	
	fmt.Println(dijan.Get("models:account:1"))
}
```
若在容器系统外访问则配置主机地址映射
```
	dijan := godijan.NewGoDijanConnection("localhost:2375", map[string]string {
		"dijan-0.dijan-service": "localhost:2375",
		"dijan-1.dijan-service": "localhost:2333",
	})
```
容器内则第二个参数可以为nil
