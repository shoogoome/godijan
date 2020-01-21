package godijan

import (
	"bufio"
	"net"
)

func NewGoDijanConnection(host string, mapping map[string]string) GoDijan {
	c, e := net.Dial("tcp", host)
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
	}
	dijan.setCircle()
	return &dijan
}
