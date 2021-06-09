package ubot

import "regexp"

type MsgEntity struct {
	Type      string
	Args      []string
	NamedArgs map[string]string
}

var MsgTypePattern = regexp.MustCompile(`^[a-z0-9_\.]+$`)

func IsValidMsgType(msgType string) bool {
	return MsgTypePattern.MatchString(msgType)
}

func (e *MsgEntity) ArgOr(i int, defaultValue string) string {
	if i < 0 || i >= len(e.Args) {
		return defaultValue
	}
	return e.Args[i]
}

func (e *MsgEntity) NamedArgOr(name string, defaultValue string) string {
	v, ok := e.NamedArgs[name]
	if !ok {
		return defaultValue
	}
	return v
}

func (e *MsgEntity) FirstArgOrEmpty() string {
	return e.ArgOr(0, "")
}
