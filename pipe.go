package godijan


/*
 pipe方法只允许在单机模式下使用
 集群模式将会有大部分键未能命中缓存的情况
*/
func (c *goDijan) PipelinedRun(cmds []*Cmd) {
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
