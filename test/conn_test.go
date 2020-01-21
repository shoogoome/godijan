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
	// 容器外访问
	dijan := godijan.NewGoDijanConnection("localhost:2375", map[string]string {
		"dijan-0.dijan-service": "localhost:2375",
		"dijan-1.dijan-service": "localhost:2333",
	})
	if err := dijan.Set("models:account:1", accountString); err != nil {
		panic(err)
	}

	fmt.Println(dijan.Get("models:account:1"))
}