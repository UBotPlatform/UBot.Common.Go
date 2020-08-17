package ubot

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

const lowerhex = "0123456789abcdef"

type MsgBuilder struct {
	data strings.Builder
}

func (b *MsgBuilder) String() string {
	return b.data.String()
}

func (b *MsgBuilder) WriteRaw(s []byte) *MsgBuilder {
	b.data.Write(s)
	return b
}

func (b *MsgBuilder) WriteRawString(s string) *MsgBuilder {
	b.data.WriteString(s)
	return b
}

func (b *MsgBuilder) WriteString(s string) *MsgBuilder {
	start := 0
	i := 0
	var r rune
	var width int
	for i = 0; i < len(s); i += width {
		r, width = rune(s[i]), 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s[i:])
		}
		switch r {
		case '\\':
			fallthrough
		case '[':
			fallthrough
		case ']':
			b.data.WriteString(s[start:i])
			start = i + 1
			b.data.WriteByte('\\')
			b.data.WriteByte(s[i])
		default:
			if r >= 0x10000 {
				b.data.WriteString(s[start:i])
				start = i + width
				b.data.WriteString(`\u{`)
				b.data.WriteString(strconv.FormatUint(uint64(r), 16))
				b.data.WriteString(`}`)
			} else if lowCodePointNeedEscape(r) {
				b.data.WriteString(s[start:i])
				start = i + width
				b.data.WriteString(`\u`)
				for s := 12; s >= 0; s -= 4 {
					b.data.WriteByte(lowerhex[r>>uint(s)&0xF])
				}
			}
			continue
		}
	}
	b.data.WriteString(s[start:i])
	return b
}

func lowCodePointNeedEscape(codePoint rune) bool {
	return (codePoint >= 0xA9 && codePoint <= 0xAE) || (codePoint >= 0x200D && codePoint <= 0x3299)
}

func (b *MsgBuilder) WriteEntity(e MsgEntity) *MsgBuilder {
	b.data.WriteByte('[')
	b.WriteString(e.Type)
	b.data.WriteByte(':')
	b.WriteString(e.Data)
	b.data.WriteByte(']')
	return b
}
