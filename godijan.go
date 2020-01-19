package godijan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type Cmd struct {
	Name  string
	Key   string
	Value string
	TTL   int
	Error error
}

type GoDijan interface {
	Get(string) (string, error)
	Set(string, string, ...int) error
	Del(string) error
	Run(*Cmd)
	PipelinedRun([]*Cmd)
	Close()
}

type goDijan struct {
	net.Conn
	r *bufio.Reader
	mapping map[string]string
}

func NewGoDijanConnection(host string, mapping map[string]string) GoDijan {
	c, e := net.Dial("tcp", host)
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)
	return &goDijan{c, r, mapping}
}

func (c *goDijan) Close() {
	c.Conn.Close()
}

func (c *goDijan) sendGet(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("G%d %s", klen, key)))
}

func (c *goDijan) sendSet(key, value string, ttl ...int) {
	klen := len(key)
	vlen := len(value)
	if len(ttl) > 0 && ttl[0] != 0 {
		ttlString := strconv.Itoa(ttl[0])
		c.Write([]byte(fmt.Sprintf("S%d %d %d %s%s%s", klen, vlen, len(ttlString), key, value, ttlString)))
	} else {
		c.Write([]byte(fmt.Sprintf("S%d %d 0 %s%s", klen, vlen, key, value)))
	}
}

func (c *goDijan) sendDel(key string) {
	klen := len(key)
	c.Write([]byte(fmt.Sprintf("D%d %s", klen, key)))
}

func (c *goDijan) Get(key string) (string, error) {
	c.sendGet(key)
	value, err := c.recvResponse()
	if err != nil {
		addr := fmt.Sprintf("%s", err)
		if c.mapping != nil {
			addr = c.mapping[addr]
		}
		conn := NewGoDijanConnection(addr, nil)
		defer conn.Close()
		return conn.Get(key)
	}
	return value, err
}

func (c *goDijan) Set(key, value string, ttl ...int) error {
	c.sendSet(key, value, ttl...)
	_, err := c.recvResponse()
	if err != nil {
		addr := fmt.Sprintf("%s", err)
		if c.mapping != nil {
			addr = c.mapping[addr]
		}
		conn := NewGoDijanConnection(addr, nil)
		defer conn.Close()
		return conn.Set(key, value, ttl...)
	}
	return err
}

func (c *goDijan) Del(key string) error {
	c.sendDel(key)
	_, err := c.recvResponse()
	return err
}

func (c *goDijan) recvResponse() (string, error) {
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

func (c *goDijan) Run(cmd *Cmd) {
	switch cmd.Name {
	case "get":
		cmd.Value, cmd.Error = c.Get(cmd.Key)
		return
	case "set":
		cmd.Error = c.Set(cmd.Key, cmd.Value, cmd.TTL)
		return
	case "del":
		cmd.Error = c.Del(cmd.Key)
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

/*
 pipe方法只允许在单机模式下使用
 集群模式将会有大部分键未能命中缓存的情况
*/
func (c *goDijan) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		if cmd.Name == "get" {
			c.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			c.sendSet(cmd.Key, cmd.Value, cmd.TTL)
		}
		if cmd.Name == "del" {
			c.sendDel(cmd.Key)
		}
	}
	for _, cmd := range cmds {
		cmd.Value, cmd.Error = c.recvResponse()
	}
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
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

