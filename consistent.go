package godijan

import (
	"bufio"
	"fmt"
	"net"
	"stathat.com/c/consistent"
	"strings"
)

func (c *goDijan) setCircle() {

	reset := false
	lock.Lock()
	defer lock.Unlock()

	var nodesList []string
	for _, conn := range c.conn {
		if reset {
			conn.Close()
			continue
		}
		if err := conn.sendMemberSignal(); err != nil {
			conn.Close()
			continue
		}
		nodes, err := conn.recvResponse()
		if err != nil {
			conn.Close()
			continue
		}
		nodesList = strings.Split(nodes, " ")
		circle := consistent.New()
		circle.NumberOfReplicas = len(nodesList)
		circle.Set(nodesList)
		c.circle = circle
		conn.Close()
		reset = true
	}

	connList := make(map[string]dijanConn, len(nodesList))
	for _, node := range nodesList {
		if c.hostnameMapping != nil {
			for i := 0; i < len(nodesList); i++ {
				if hostname, ok := c.hostnameMapping[nodesList[i]]; ok {
					node = hostname
				}
			}
		}
		co, e := net.Dial("tcp", fmt.Sprintf("%s:%d", node, c.port))
		if e != nil {
			panic(e)
		}
		r := bufio.NewReader(co)
		connList[node] = dijanConn {
			Conn: co,
			r: r,
		}
	}
	c.conn = connList
}


