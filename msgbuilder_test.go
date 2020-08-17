package ubot_test

import (
	"testing"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

func TestMsgBuilder(t *testing.T) {
	var builder ubot.MsgBuilder
	builder.WriteString("hello, [go]")
	builder.WriteString("\U0001F606")
	builder.WriteString("\u303D")
	builder.WriteEntity(ubot.MsgEntity{Type: "at", Data: "10000"})
	r := builder.String()
	if r != `hello, \[go\]\u{1f606}\u303d[at:10000]` {
		t.Error("got wrong result: " + r)
	}
}
