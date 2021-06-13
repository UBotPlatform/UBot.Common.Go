package ubot_test

import (
	"reflect"
	"testing"
	"unicode"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

func TestMsgParser(t *testing.T) {
	entities := ubot.ParseMsg(`asd\[hh\][at:123][!~!~!~]`)
	if !reflect.DeepEqual(entities, []ubot.MsgEntity{
		{Type: "text", Args: []string{"asd[hh]"}},
		{Type: "at", Args: []string{"123"}},
		{Type: "text", Args: []string{"[!~!~!~]"}},
	}) {
		t.Fail()
	}

	//UTF-16 Escape
	entities = ubot.ParseMsg(`\u3064|asd`)
	if entities[0].Type != "text" || entities[0].Args[0] != "\u3064|asd" {
		t.Fail()
	}
	//UTF-16 Surrogate Pair
	entities = ubot.ParseMsg(`\ud83d\ude06|asd`)
	if entities[0].Type != "text" || entities[0].Args[0] != "\U0001F606|asd" {
		t.Fail()
	}
	//Invaild Surrogate Pair
	entities = ubot.ParseMsg(`\ud83d|asd`)
	if entities[0].Type != "text" || entities[0].Args[0] != string(unicode.ReplacementChar)+"|asd" {
		t.Fail()
	}
	//Unicode point escape
	entities = ubot.ParseMsg(`\u{3064}|asd`)
	if entities[0].Type != "text" || entities[0].Args[0] != "\u3064|asd" {
		t.Fail()
	}
	entities = ubot.ParseMsg(`\u{1f606}|asd`)
	if entities[0].Type != "text" || entities[0].Args[0] != "\U0001F606|asd" {
		t.Fail()
	}

	//Args Test
	entities = ubot.ParseMsg(`[image:<url>1\,2\=3,md5=xxx][file:<xxx>,a=b]`)
	if !reflect.DeepEqual(entities, []ubot.MsgEntity{
		{Type: "image", Args: []string{"<url>1,2=3"}, NamedArgs: map[string]string{"md5": "xxx"}},
		{Type: "file", Args: []string{"<xxx>"}, NamedArgs: map[string]string{"a": "b"}},
	}) {
		t.Fail()
	}
}
