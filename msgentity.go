package ubot

import "regexp"

type MsgEntity struct {
	Type string
	Data string
}

var MsgTypePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

func IsValidMsgType(msgType string) bool {
	return MsgTypePattern.MatchString(msgType)
}
