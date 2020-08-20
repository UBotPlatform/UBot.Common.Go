package ubot

func defaultAppOnReceiveChatMessage(bot string, msgType MsgType, source string, sender string, message string, info MsgInfo) (EventResultType, error) {
	return IgnoreEvent, nil
}
func defaultAppOnMemberJoined(bot string, source string, sender string, inviter string) (EventResultType, error) {
	return IgnoreEvent, nil
}
func defaultAppOnMemberLeft(bot string, source string, sender string) (EventResultType, error) {
	return IgnoreEvent, nil
}
func defaultAppProcessGroupInvitation(bot string, sender string, target string, reason string) (EventResultType, *string, error) {
	return IgnoreEvent, nil, nil
}
func defaultAppProcessFriendRequest(bot string, sender string, reason string) (EventResultType, *string, error) {
	return IgnoreEvent, nil, nil
}
func defaultAppProcessMembershipRequest(bot string, source string, sender string, inviter string, reason string) (EventResultType, *string, error) {
	return IgnoreEvent, nil, nil
}
