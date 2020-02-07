package godijan

type goDijanPool struct {
	conn chan *goDijan
}

func (p *goDijanPool) Get(key string) ([]byte, error) {
	conn := <- p.conn
	defer func() {
		p.conn <- conn
	}()
	return conn.Get(key)
}

func (p *goDijanPool) Set(key string, value []byte, ttl ...int) error {
	conn := <- p.conn
	defer func() {
		p.conn <- conn
	}()
	return conn.Set(key, value, ttl...)
}

func (p *goDijanPool) Del(key string) error {
	conn := <- p.conn
	defer func() {
		p.conn <- conn
	}()
	return conn.Del(key)
}

func (p *goDijanPool) Run(cmd *Cmd) {
	c := <- p.conn
	defer func() {
		p.conn <- c
	}()
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

func (p *goDijanPool) PipelinedRun(cmds []*Cmd) {
	//if len(cmds) == 0 {
	//	return
	//}
	//for _, cmd := range cmds {
	//	if cmd.Name == "get" {
	//		c.sendGet(cmd.Key)
	//	}
	//	if cmd.Name == "set" {
	//		c.sendSet(cmd.Key, cmd.Value, cmd.TTL)
	//	}
	//	if cmd.Name == "del" {
	//		c.sendDel(cmd.Key)
	//	}
	//}
	//for _, cmd := range cmds {
	//	cmd.Value, cmd.Error = c.recvResponse()
	//}
}