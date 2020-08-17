package ubot_test

import (
	"testing"
	"unicode"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

func TestMsgParser(t *testing.T) {
	entities := ubot.ParseMsg(`asd\[hh\][at:123][!~!~!~]`)
	if entities[0].Type != "text" || entities[0].Data != "asd[hh]" {
		t.Fail()
	}
	if entities[1].Type != "at" || entities[1].Data != "123" {
		t.Fail()
	}
	if entities[2].Type != "text" || entities[2].Data != "[!~!~!~]" {
		t.Fail()
	}

	//UTF-16 Escape
	entities = ubot.ParseMsg(`\u3064|asd`)
	if entities[0].Type != "text" || entities[0].Data != "\u3064|asd" {
		t.Fail()
	}
	//UTF-16 Surrogate Pair
	entities = ubot.ParseMsg(`\ud83d\ude06|asd`)
	if entities[0].Type != "text" || entities[0].Data != "\U0001F606|asd" {
		t.Fail()
	}
	//Invaild Surrogate Pair
	entities = ubot.ParseMsg(`\ud83d|asd`)
	if entities[0].Type != "text" || entities[0].Data != string(unicode.ReplacementChar)+"|asd" {
		t.Fail()
	}
	//Unicode point escape
	entities = ubot.ParseMsg(`\u{3064}|asd`)
	if entities[0].Type != "text" || entities[0].Data != "\u3064|asd" {
		t.Fail()
	}
	entities = ubot.ParseMsg(`\u{1f606}|asd`)
	if entities[0].Type != "text" || entities[0].Data != "\U0001F606|asd" {
		t.Fail()
	}
}
