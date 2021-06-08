package ubot

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseMsg(content string) []MsgEntity {
	r := make([]MsgEntity, 0, 4)
	var curData strings.Builder
	inBracket := false
	var curType string
	start, count := 0, 0
	var i int
	var args []string
	var namedArgs map[string]string
	var isNamedArg bool
	var curArgName string
	flushBuf := func() {
		curData.WriteString(content[start : start+count])
		start = i + 1
		count = 0
	}
	finishArg := func() {
		flushBuf()
		if isNamedArg {
			if namedArgs == nil {
				namedArgs = make(map[string]string)
			}
			namedArgs[curArgName] = curData.String()
		} else {
			args = append(args, curData.String())
		}
		curData.Reset()
		curArgName = ""
		isNamedArg = false
	}
	endText := func() {
		flushBuf()
		if curData.Len() != 0 {
			r = append(r, MsgEntity{Type: "text", Args: []string{curData.String()}})
		}
		curData.Reset()
	}
	beginEntity := func(newType string) {
		endText()
		args = make([]string, 0, 1)
		namedArgs = nil
		curType = newType
	}
	endEntity := func() {
		finishArg()
		r = append(r, MsgEntity{Type: curType, Args: args, NamedArgs: namedArgs})
		args = make([]string, 0, 1)
	}
	for i = 0; i < len(content); i++ {
		switch content[i] {
		case '[':
			if inBracket {
				count++
				break
			}
			j := i
			for j < len(content) && content[j] != ':' && content[j] != ']' {
				j++
			}
			if j >= len(content) || content[j] == ']' {
				count++
				break
			}
			newType := content[i+1 : j]
			if !IsValidMsgType(newType) {
				count++
				break
			}
			i = j
			beginEntity(newType)
			inBracket = true
		case ']':
			if inBracket {
				endEntity()
				inBracket = false
			} else {
				count++
			}
		case ',':
			if inBracket {
				finishArg()
			} else {
				count++
			}
		case '=':
			if inBracket && !isNamedArg {
				flushBuf()
				curArgName = curData.String()
				isNamedArg = true
				curData.Reset()
			} else {
				count++
			}
		case '\\':
			if i+1 >= len(content) {
				count++
				break
			}
			switch content[i+1] {
			case '[':
				fallthrough
			case ']':
				fallthrough
			case ',':
				fallthrough
			case '=':
				fallthrough
			case '\\':
				i++
				flushBuf()
				curData.WriteByte(content[i])
			case 'n':
				i++
				flushBuf()
				curData.WriteByte('\n')
			case 'r':
				i++
				flushBuf()
				curData.WriteByte('\r')
			case 't':
				i++
				flushBuf()
				curData.WriteByte('\t')
			case 'u':
				if i+2 >= len(content) {
					count++
					break
				}
				if content[i+2] == '{' {
					pEnd := strings.IndexByte(content[i+3:], '}')
					if pEnd == -1 {
						count++
						break
					}
					pEnd += i + 3
					code, err := strconv.ParseUint(content[i+3:pEnd], 16, 32)
					if err != nil {
						count++
						break
					}
					i = pEnd
					flushBuf()
					curData.WriteRune(rune(code))
					break
				}
				if i+5 >= len(content) {
					count++
					break
				}
				code1, err := strconv.ParseUint(content[i+2:i+6], 16, 16)
				if err != nil {
					count++
					break
				}
				i += 5
				if code1 >= 0xD800 && code1 <= 0xDBFF {
					// UTF-16 Surrogate Pair
					if i+6 >= len(content) {
						flushBuf()
						curData.WriteRune(rune(unicode.ReplacementChar))
						break
					}
					code2, _ := strconv.ParseUint(content[i+3:i+7], 16, 16)
					if content[i+1] != '\\' || content[i+2] != 'u' || code2 < 0xDC00 || code2 > 0xDFFF {
						flushBuf()
						curData.WriteRune(rune(unicode.ReplacementChar))
						break
					}
					i += 6
					flushBuf()
					curData.WriteRune(rune((code1-0xD800)<<10 | (code2 - 0xDC00) + 0x10000))
				} else {
					flushBuf()
					curData.WriteRune(rune(code1))
				}
			default:
				count++
			}
		default:
			count++
		}
	}
	endText()
	return r
}
