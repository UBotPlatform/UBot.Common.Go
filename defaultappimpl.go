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
