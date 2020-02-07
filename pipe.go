package godijan


func (c *goDijan) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	c.Lock()
	defer c.Unlock()

	resultChan := make(chan dijanConn, len(cmds) + 5)
	for _, cmd := range cmds {
		conn := c.getConn(cmd.Key)
		if cmd.Name == "get" {
			conn.sendGet(cmd.Key)
		}
		if cmd.Name == "set" {
			conn.sendSet(cmd.Key, cmd.Value, cmd.TTL)
		}
		if cmd.Name == "del" {
			conn.sendDel(cmd.Key)
		}
		resultChan <- conn
	}
	for _, cmd := range cmds {
		conn := <-resultChan
		cmd.Value, cmd.Error = conn.recvResponse()
	}
}


