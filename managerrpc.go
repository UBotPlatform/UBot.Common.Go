package ubot

import "github.com/1354092549/wsrpc"

type Manager struct {
	RegisterApp     func(id string) (string, error)
	RegisterAccount func(id string) (string, error)
}

func (a *Manager) Get(rpcConn *wsrpc.WebsocketRPCConn) {
	rpcConn.MakeCall("register_app", &a.RegisterApp, nil, nil)
	rpcConn.MakeCall("register_account", &a.RegisterAccount, nil, nil)
}
