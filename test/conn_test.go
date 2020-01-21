package test

import (
	"fmt"
	"godijan"
	"testing"
)

func TestConn(t *testing.T) {


	dijan := godijan.NewGoDijanConnection("localhost:2375", map[string]string {
		"server0:2375": "localhost:2375",
		"server1:2333": "localhost:2333",
		"server0:2333": "localhost:2375",
		"server1:2375": "localhost:2333",
	})
	err := dijan.Set("models:account:describe:1", "ewfwefwefwefwe")
	fmt.Println("zuihou", err)
	//dijan.Get("models:account:1")


}