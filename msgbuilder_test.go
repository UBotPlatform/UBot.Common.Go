package ubot_test

import (
	"testing"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

func TestMsgBuilder(t *testing.T) {
	var builder ubot.MsgBuilder
	builder.WriteString("hello, [go]1,2=3")
	builder.WriteString("\U0001F606")
	builder.WriteString("\u25B6")
	builder.WriteEntity(ubot.MsgEntity{Type: "at", Args: []string{"10000"}})
	builder.WriteEntity(ubot.MsgEntity{Type: "image", Args: []string{"<url>1,2=3"}, NamedArgs: map[string]string{"md5": "xxx"}})
	r := builder.String()
	if r != `hello, \[go\]1,2=3\u{1f606}\u25b6[at:10000][image:<url>1\,2\=3,md5=xxx]` {
		t.Error("got wrong result: " + r)
	}
}
