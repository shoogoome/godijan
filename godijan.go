package godijan

import (
	"fmt"
	"stathat.com/c/consistent"
	"sync"
)

var lock sync.RWMutex
type connMapping map[string]dijanConn

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
}

type goDijan struct {
	conn connMapping
	hostnameMapping map[string]string
	circle *consistent.Consistent
}

func (c *goDijan) Get(key string) (string, error) {
	conn := c.getConn(key)
	if err := conn.sendGet(key); err != nil {
		c.setCircle()
		return c.Get(key)
	}
	return conn.recvResponse()
}

func (c *goDijan) Set(key, value string, ttl ...int) error {
	conn := c.getConn(key)
	fmt.Println("?")
	if err := conn.sendSet(key, value, ttl...); err != nil {
		fmt.Println("!", err)
		c.setCircle()
		//return c.Set(key, value, ttl...)
	}
	_, err := conn.recvResponse()
	return err
}

func (c *goDijan) Del(key string) error {
	conn := c.getConn(key)
	if err := conn.sendDel(key); err != nil {
		c.setCircle()
		return c.Del(key)
	}
	_, err := conn.recvResponse()
	return err
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


// 获取key对应主机tcp连接
func (c *goDijan) getConn(key string) dijanConn {
	count := 5
GET:
	hostname, err := c.circle.Get(key)
	fmt.Println("key", hostname)
	if err != nil {
		if count <= 0 {
			panic("[!] consistent setting fail")
		}
		c.setCircle()
		count -=  1
		goto GET
	}
	if c.hostnameMapping != nil {
		if nHostname, ok := c.hostnameMapping[hostname]; ok {
			hostname = nHostname
		}
	}
	lock.RLock()
	conn, ok := c.conn[hostname]
	lock.RUnlock()
	if !ok {
		if count <= 0 {
			panic("[!] consistent setting fail")
		}
		c.setCircle()
		count -=  1
		goto GET
	}
	//fmt.Println("!!")
	return conn
}



