package ubot_test

import (
	"testing"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

func TestMsgTypeValidator(t *testing.T) {
	if !ubot.IsValidMsgType("at") {
		t.Fail()
	}
	if ubot.IsValidMsgType("!") {
		t.Fail()
	}
}
