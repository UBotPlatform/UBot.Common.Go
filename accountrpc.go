package ubot

import (
	"github.com/1354092549/wsrpc"
)

type AccountEventEmitter struct {
	OnReceiveChatMessage     func(msgType MsgType, source string, sender string, message string, info MsgInfo) error
	OnMemberJoined           func(source string, sender string, inviter string) error
	OnMemberLeft             func(source string, sender string) error
	ProcessGroupInvitation   func(sender string, target string, reason string) (EventResultType, *string, error)
	ProcessFriendRequest     func(sender string, reason string) (EventResultType, *string, error)
	ProcessMembershipRequest func(source string, sender string, inviter string, reason string) (EventResultType, *string, error)
}

func (a *AccountEventEmitter) Get(rpcConn *wsrpc.WebsocketRPCConn) {
	rpcConn.MakeNotify("on_receive_chat_message", &a.OnReceiveChatMessage, nil)
	rpcConn.MakeNotify("on_member_joined", &a.OnMemberJoined, nil)
	rpcConn.MakeNotify("on_member_left", &a.OnMemberLeft, nil)
	rpcConn.MakeCall("process_group_invitation", &a.ProcessGroupInvitation, nil, []string{"type", "reason"})
	rpcConn.MakeCall("process_friend_request", &a.ProcessFriendRequest, nil, []string{"type", "reason"})
	rpcConn.MakeCall("process_membership_request", &a.ProcessMembershipRequest, nil, []string{"type", "reason"})
}

type Account struct {
	GetGroupName    func(id string) (string, error)
	GetUserName     func(id string) (string, error)
	SendChatMessage func(msgType MsgType, source string, target string, message string) error
	RemoveMember    func(source string, target string) error
	ShutupMember    func(source string, target string, duration int) error
	ShutupAllMember func(source string, shutupSwitch bool) error
	GetMemberName   func(source string, target string) (string, error)
	GetUserAvatar   func(id string) (string, error)
	GetSelfID       func() (string, error)
	GetPlatformID   func() (string, error)
	GetGroupList    func() ([]string, error)
	GetMemberList   func(id string) ([]string, error)
}

func (a *Account) Register(rpc *wsrpc.WebsocketRPC) {
	rpc.Register("get_group_name",
		a.GetGroupName,
		[]string{"id"},
		nil)
	rpc.Register("get_user_name",
		a.GetUserName,
		[]string{"id"},
		nil)
	rpc.Register("send_chat_message",
		a.SendChatMessage,
		[]string{"type", "source", "target", "message"},
		nil)
	rpc.Register("remove_member",
		a.RemoveMember,
		[]string{"source", "target"},
		nil)
	rpc.Register("shutup_member",
		a.ShutupMember,
		[]string{"source", "target", "duration"},
		nil)
	rpc.Register("shutup_all_member",
		a.ShutupAllMember,
		[]string{"source", "switch"},
		nil)
	rpc.Register("get_member_name",
		a.GetMemberName,
		[]string{"source", "target"},
		nil)
	rpc.Register("get_user_avatar",
		a.GetUserAvatar,
		[]string{"id"},
		nil)
	rpc.Register("get_self_id",
		a.GetSelfID,
		nil,
		nil)
	rpc.Register("get_platform_id",
		a.GetPlatformID,
		nil,
		nil)
	rpc.Register("get_group_list",
		a.GetGroupList,
		nil,
		nil)
	rpc.Register("get_member_list",
		a.GetMemberList,
		[]string{"id"},
		nil)
}
