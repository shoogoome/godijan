package godijan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type dijanConn struct {
	net.Conn
	r *bufio.Reader
}

func (c *dijanConn) sendGet(key string) error {
	kLen := len(key)
	_, err := c.Write([]byte(fmt.Sprintf("G%d %s", kLen, key)))
	return err
}

func (c *dijanConn) sendSet(key, value string, ttl ...int) error {
	kLen := len(key)
	vlen := len(value)
	var err error
	if len(ttl) > 0 && ttl[0] != 0 {
		ttlString := strconv.Itoa(ttl[0])
		_, err = c.Write([]byte(fmt.Sprintf("S%d %d %d %s%s%s", kLen, vlen, len(ttlString), key, value, ttlString)))
	} else {
		_, err = c.Write([]byte(fmt.Sprintf("S%d %d 0 %s%s", kLen, vlen, key, value)))
	}
	return err
}

func (c *dijanConn) sendDel(key string) error {
	kLen := len(key)
	_, err := c.Write([]byte(fmt.Sprintf("D%d %s", kLen, key)))
	return err
}

func (c *dijanConn) sendMemberSignal() error {
	_, err := c.Write([]byte("M"))
	return err
}

func (c *dijanConn) recvResponse() (string, error) {
	vlen := readLen(c.r)
	if vlen == 0 {
		return "", nil
	}
	if vlen < 0 {
		err := make([]byte, -vlen)
		_, e := io.ReadFull(c.r, err)
		if e != nil {
			return "", e
		}
		return "", errors.New(string(err))
	}
	value := make([]byte, vlen)
	_, e := io.ReadFull(c.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

