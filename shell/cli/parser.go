package cli

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

var ParseEnv bool = false
var ParseBacktick bool = false

func getEnv(_ string) string {
	return ""
}

func replaceEnv(env func(string) string, s string) string {
	if env == nil {
		env = getEnv
	}

	var buf bytes.Buffer
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		r := rs[i]
		if r == '\\' {
			i++
			if i == len(rs) {
				break
			}
			buf.WriteRune(rs[i])
			continue
		} else if r == '$' {
			i++
			if i == len(rs) {
				buf.WriteRune(r)
				break
			}
			if rs[i] == 0x7b {
				i++
				p := i
				for ; i < len(rs); i++ {
					r = rs[i]
					if r == '\\' {
						i++
						if i == len(rs) {
							return s
						}
						continue
					}
					if r == 0x7d || (!unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r)) {
						break
					}
				}
				if r != 0x7d {
					return s
				}
				if i > p {
					buf.WriteString(env(s[p:i]))
				}
			} else {
				p := i
				for ; i < len(rs); i++ {
					r := rs[i]
					if r == '\\' {
						i++
						if i == len(rs) {
							return s
						}
						continue
					}
					if !unicode.IsLetter(r) && r != '_' && !unicode.IsDigit(r) {
						break
					}
				}
				if i > p {
					buf.WriteString(env(s[p:i]))
					i--
				} else {
					buf.WriteString(s[p:])
				}
			}
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

type Parser struct {
	ParseEnv      bool
	ParseBacktick bool
	Position      int
	Dir           string
	GetEnv        func(string) string
}

func NewParser() *Parser {
	return &Parser{
		ParseEnv:      ParseEnv,
		ParseBacktick: ParseBacktick,
		Position:      0,
		Dir:           "",
	}
}

type argType int

const (
	argNo argType = iota
	argSingle
	argQuoted
)

func (p *Parser) Parse(line string) ([]string, error) {
	var args []string
	buf := ""
	var escaped, doubleQuoted, singleQuoted, backQuote, dollarQuote bool
	backtick := ""

	pos := -1
	got := argNo

	i := -1
loop:
	for _, r := range line {
		i++
		if escaped {
			if r == 't' {
				r = '\t'
			}
			if r == 'n' {
				r = '\n'
			}
			buf += string(r)
			escaped = false
			got = argSingle
			continue
		}

		if r == '\\' {
			if singleQuoted {
				buf += string(r)
			} else {
				escaped = true
			}
			continue
		}

		if unicode.IsSpace(r) {
			if singleQuoted || doubleQuoted || backQuote || dollarQuote {
				buf += string(r)
				backtick += string(r)
			} else if got != argNo {
				if p.ParseEnv {
					if got == argSingle {
						parser := &Parser{ParseEnv: false, ParseBacktick: false, Position: 0, Dir: p.Dir}
						strs, err := parser.Parse(replaceEnv(p.GetEnv, buf))
						if err != nil {
							return nil, err
						}
						args = append(args, strs...)
					} else {
						args = append(args, replaceEnv(p.GetEnv, buf))
					}
				} else {
					args = append(args, buf)
				}
				buf = ""
				got = argNo
			}
			continue
		}

		switch r {
		case '`':
			if !singleQuoted && !doubleQuoted && !dollarQuote {
				if p.ParseBacktick {
					if backQuote {
						buf = buf[:len(buf)-len(backtick)]
					}
					backtick = ""
					backQuote = !backQuote
					continue
				}
				backtick = ""
				backQuote = !backQuote
			}
		case ')':
			if !singleQuoted && !doubleQuoted && !backQuote {
				if p.ParseBacktick {
					if dollarQuote {
						buf = buf[:len(buf)-len(backtick)-2]
					}
					backtick = ""
					dollarQuote = !dollarQuote
					continue
				}
				backtick = ""
				dollarQuote = !dollarQuote
			}
		case '(':
			if !singleQuoted && !doubleQuoted && !backQuote {
				if !dollarQuote && strings.HasSuffix(buf, "$") {
					dollarQuote = true
					buf += "("
					continue
				} else {
					return nil, errors.New("invalid command line string")
				}
			}
		case '"':
			if !singleQuoted && !dollarQuote {
				if doubleQuoted {
					got = argQuoted
				}
				doubleQuoted = !doubleQuoted
				continue
			}
		case '\'':
			if !doubleQuoted && !dollarQuote {
				if singleQuoted {
					got = argQuoted
				}
				singleQuoted = !singleQuoted
				continue
			}
		case ';', '&', '|', '<', '>':
			if !(escaped || singleQuoted || doubleQuoted || backQuote || dollarQuote) {
				if r == '>' && len(buf) > 0 {
					if c := buf[0]; '0' <= c && c <= '9' {
						i -= 1
						got = argNo
					}
				}
				pos = i
				break loop
			}
		}

		got = argSingle
		buf += string(r)
		if backQuote || dollarQuote {
			backtick += string(r)
		}
	}

	if got != argNo {
		if p.ParseEnv {
			if got == argSingle {
				parser := &Parser{ParseEnv: false, ParseBacktick: false, Position: 0, Dir: p.Dir}
				strs, err := parser.Parse(replaceEnv(p.GetEnv, buf))
				if err != nil {
					return nil, err
				}
				args = append(args, strs...)
			} else {
				args = append(args, replaceEnv(p.GetEnv, buf))
			}
		} else {
			args = append(args, buf)
		}
	}

	if escaped || singleQuoted || doubleQuoted || backQuote || dollarQuote {
		return nil, errors.New("invalid command line string")
	}

	p.Position = pos

	return args, nil
}

func (p *Parser) ParseWithEnvs(line string) (envs []string, args []string, err error) {
	_args, err := p.Parse(line)
	if err != nil {
		return nil, nil, err
	}
	envs = []string{}
	args = []string{}
	parsingEnv := true
	for _, arg := range _args {
		if parsingEnv && isEnv(arg) {
			envs = append(envs, arg)
		} else {
			if parsingEnv {
				parsingEnv = false
			}
			args = append(args, arg)
		}
	}
	return envs, args, nil
}

func isEnv(arg string) bool {
	return len(strings.Split(arg, "=")) == 2
}

func Parse(line string) ([]string, error) {
	return NewParser().Parse(line)
}

/*
func Parse(line string) ([]string, bool) {
	args := strings.Split(line, " ")
	if len(args) <= 0 {
		return nil, false
	}
	return args, true
}
*/
