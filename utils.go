package godijan

import (
	"bufio"
	"log"
	"strconv"
	"strings"
)

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	tmp = strings.Trim(tmp, "\x00")
	if e != nil {
		log.Println(e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(tmp, e)
		return 0
	}
	return l
}
