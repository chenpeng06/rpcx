package server

import (
	"net"

	"github.com/smallnest/rpcx/protocol"
)

type PluginContainer interface {
	Add(plugin Plugin)
	Remove(plugin Plugin)
	All() []Plugin

	DoRegister(name string, rcvr interface{}, metadata string) error
	DoRegisterFunction(serviceName, fname string, fn interface{}, metadata string) error
	DoUnregister(name string) error

	DoPostConnAccept(net.Conn) (net.Conn, bool)
	DoPostConnClose(net.Conn)

	DoPreReadRequest(*Context) error
	DoPostReadRequest(*Context, *protocol.Message) error
}

type Plugin interface{}

type pluginContainer struct {
	plugins []Plugin
}
