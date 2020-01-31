package godijan

import (
	"bufio"
	"fmt"
	"net"
)

func NewGoDijanConnection(host string, port int, node int, mapping map[string]string) GoDijan {
	c, e := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if e != nil {
		panic(e)
	}
	r := bufio.NewReader(c)

	dijan := goDijan{
		conn: map[string]dijanConn{
			host: {
				c, r,
			},
		},
		hostnameMapping: mapping,
		circle:          nil,
		port: port,
		node: node,
	}
	dijan.setCircle()
	return &dijan
}
