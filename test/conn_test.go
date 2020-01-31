package test

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestConn(t *testing.T) {

	a := bufio.NewReader(bytes.NewReader([]byte("-17 asdasdasdasd")))
	tmp, e := a.ReadString(' ')
	if e != nil {
		log.Println(e)
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(tmp, e)
	}
}
