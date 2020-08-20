package ubot

import (
	"github.com/1354092549/wsrpc"
)

type AppApi struct {
	GetGroupName    func(bot string, id string) (string, error)
	GetUserName     func(bot string, id string) (string, error)
	SendChatMessage func(bot string, msgType MsgType, source string, target string, message string) error
	RemoveMember    func(bot string, source string, target string) error
	ShutupMember    func(bot string, source string, target string, duration int) error
	ShutupAllMember func(bot string, source string, shutupSwitch bool) error
	GetMemberName   func(bot string, source string, target string) (string, error)
	GetUserAvatar   func(bot string, id string) (string, error)
	GetSelfID       func(bot string) (string, error)
}

func (a *AppApi) Get(rpcConn *wsrpc.WebsocketRPCConn) {
	rpcConn.MakeCall("get_group_name", &a.GetGroupName, nil, nil)
	rpcConn.MakeCall("get_user_name", &a.GetUserName, nil, nil)
	rpcConn.MakeCall("send_chat_message", &a.SendChatMessage, nil, nil)
	rpcConn.MakeCall("remove_member", &a.RemoveMember, nil, nil)
	rpcConn.MakeCall("shutup_member", &a.ShutupMember, nil, nil)
	rpcConn.MakeCall("shutup_all_member", &a.ShutupAllMember, nil, nil)
	rpcConn.MakeCall("get_member_name", &a.GetMemberName, nil, nil)
	rpcConn.MakeCall("get_user_avatar", &a.GetUserAvatar, nil, nil)
	rpcConn.MakeCall("get_self_id", &a.GetSelfID, nil, nil)
}

type App struct {
	OnReceiveChatMessage     func(bot string, msgType MsgType, source string, sender string, message string, info MsgInfo) (EventResultType, error)
	OnMemberJoined           func(bot string, source string, sender string, inviter string) (EventResultType, error)
	OnMemberLeft             func(bot string, source string, sender string) (EventResultType, error)
	ProcessGroupInvitation   func(bot string, sender string, target string, reason string) (EventResultType, *string, error)
	ProcessFriendRequest     func(bot string, sender string, reason string) (EventResultType, *string, error)
	ProcessMembershipRequest func(bot string, source string, sender string, reason string) (EventResultType, *string, error)
}

func (a *App) Register(rpc *wsrpc.WebsocketRPC) {
	rpc.Register("on_receive_chat_message",
		withDefault(a.OnReceiveChatMessage, defaultAppOnReceiveChatMessage),
		[]string{"bot", "type", "source", "sender", "message", "info"},
		[]string{"type"})
	rpc.Register("on_member_joined",
		withDefault(a.OnMemberJoined, defaultAppOnMemberJoined),
		[]string{"bot", "source", "sender", "inviter"},
		[]string{"type"})
	rpc.Register("on_member_left",
		withDefault(a.OnMemberLeft, defaultAppOnMemberLeft),
		[]string{"bot", "source", "sender"},
		[]string{"type"})
	rpc.Register("process_group_invitation",
		withDefault(a.ProcessGroupInvitation, defaultAppProcessGroupInvitation),
		[]string{"bot", "sender", "target", "reason"},
		[]string{"type", "reason"})
	rpc.Register("process_friend_request",
		withDefault(a.ProcessFriendRequest, defaultAppProcessFriendRequest),
		[]string{"bot", "sender", "reason"},
		[]string{"type", "reason"})
	rpc.Register("process_membership_request",
		withDefault(a.ProcessMembershipRequest, defaultAppProcessMembershipRequest),
		[]string{"bot", "source", "sender", "reason"},
		[]string{"type", "reason"})
}
