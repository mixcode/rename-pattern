package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidPattern  = errors.New("invalid match pattern")
	ErrInvalidPosition = errors.New("invalid match position")
	ErrInvalidNumber   = errors.New("match is not a number")
)

type Pattern struct {
	Type int
	Word string
	Pos  int
}

// separate filename matching pattern with ? and * into string and ?/*
// func parseSearchPattern(pattern string) []string {
func parseSearchPattern(pattern string) []Pattern {
	if len(pattern) == 0 {
		return nil
	}

	p := make([]Pattern, 0)
	word := ""
	addWord := func(t int) {
		if word != "" {
			p = append(p, Pattern{Type: t, Word: word})
			word = ""
		}
	}

	prevc := rune(0)
	for _, c := range pattern {
		if prevc == '\\' {
			word = word + string(c)
			prevc = c
			continue
		}

		switch c {

		case '|': // pattern separator
			addWord(0)
			c = 0

		case ':': // digits
			addWord(0)
			word = "d"
			addWord(':')

		case '*':
			//if prevc != '*' { // ignore multiple asterisks
			addWord(0)
			word = "*"
			addWord('*')
			//}

		case '?':
			if word == "" && prevc == '?' {
				// join multiple ?'s into one entry
				p[len(p)-1].Word += "?"
			} else {
				addWord(0)
				word = "?"
				addWord('?')
			}

		case 0:
			// do nothing

		default:
			word = word + string(c)
		}
		prevc = c
	}
	addWord(0)
	return p
}

func parseReplacePattern(pattern string) ([]Pattern, error) {
	if len(pattern) == 0 {
		return nil, nil
	}

	p := make([]Pattern, 0)
	word := ""
	addWordW := func(t int, w string, pos int) int {
		if w != "" {
			p = append(p, Pattern{Type: t, Word: w, Pos: pos})
			word = ""
			return len(p)
		}
		return -1
	}
	addWord := func(t int) int {
		return addWordW(t, word, 0)
	}

	pt := []rune(pattern)
	idx := 0
	read := func() rune {
		if idx >= len(pt) {
			if idx == len(pt) {
				// read beyond the input
				// lock the idex to prevent unread()
				idx++
			}
			return 0
		}
		var r rune
		r, idx = pt[idx], idx+1
		return r
	}
	unread := func() {
		if 0 < idx && idx <= len(pt) {
			idx = idx - 1
		}
	}

	prevc := rune(0)
	for {
		c := read()
		if c == 0 {
			break
		}

		// escaped string
		if prevc == '\\' {
			word = word + string(c)
			prevc = c
			continue
		}

		// multi-length word
		switch c {

		case '%': // printf number formatting. %00d
			addWord(0)
			prevc = 0

			pos := 0 // index position
			c := read()
			if c == 0 {
				return nil, ErrInvalidPattern
			}
			if c == '[' {
				// read index
				s := ""
				for {
					c := read()
					if c == ']' {
						if s == "" {
							return nil, ErrInvalidNumber
						}
						pos, _ = strconv.Atoi(s)
						break
					} else if c >= '0' && c <= '9' {
						s += string(c)
					} else {
						return nil, ErrInvalidPattern
					}
				}
			} else {
				unread()
			}
			for {
				c := read()
				if c == 0 {
					return nil, ErrInvalidPattern
				}
				word = word + string(c)
				if c >= 'a' && c <= 'z' {
					break
				}
			}
			addWordW('%', word, pos)
			continue

		case '$': // word pointer. $1 or ${1} for the first matching word
			addWord(0)
			prevc = 0
			braceOpen := false
			c = read()
			if c == 0 {
				return nil, ErrInvalidPattern
			}
			if c == '{' {
				braceOpen = true
			} else {
				unread()
			}
			for {
				c = read()
				if c == 0 || c < '0' || c > '9' {
					break
				}
				word = word + string(c)
			}
			if word == "" {
				return nil, ErrInvalidNumber
			}
			pos, _ := strconv.Atoi(word)
			addWordW('$', word, pos)
			if braceOpen {
				if c != '}' {
					return nil, ErrInvalidPattern
				}
			} else {
				unread()
			}
			continue

		// single-length word
		case '*':
			//if prevc != '*' {
			addWord(0)
			addWordW('*', "*", 0)
			//}

		case '?':
			if word == "" && prevc == '?' {
				// join multiple ?'s into one entry
				p[len(p)-1].Word += "?"
			} else {
				addWord(0)
				addWordW('?', "?", 0)
			}
		case 0:
			// do nothing
		default:
			word = word + string(c)
		}
		prevc = c
	}
	addWord(0)
	return p, nil
}

/*
// rename simple search pattern to another search pattern
func Rename_OLD(filename, pattern, replace string) (string, error) {
	pSrc, pDest := parseSearchPattern(pattern), parseSearchPattern(replace)

	rSrc := ""
	srcLen := 0
	if len(pSrc) == 0 {
		rSrc = "^.*$"
	} else {
		rSrc = "^(.*?)"
		srcLen++
		for _, c := range pSrc {
			switch c.Type {
			case '*':
				rSrc += "(.*?)" // minimal length
			case '?':
				// Word may has multiple '?' chars, which represents 1 char
				rSrc += "(" + strings.Replace(c.Word, "?", ".", -1) + ")"
			case ':': // digits
				rSrc += "(\\d+)"
			case 0:
				rSrc += "(" + regexp.QuoteMeta(c.Word) + ")"
			default:
				return "", ErrInvalidPattern
			}
			srcLen++
		}
		rSrc += "(.*?)$"
		srcLen++
	}
	re, err := regexp.Compile(rSrc)
	if err != nil {
		return "", err
	}

	sDest := ""
	if len(pDest) == 0 {
		sDest = "${0}"
	} else {
		sDest = "${1}"
		i := 2
		for _, c := range pDest {
			if (c.Type == '*' || c.Type == '?') && i <= srcLen {
				sDest += fmt.Sprintf("${%d}", i)
			} else {
				//sDest += regexp.QuoteMeta(c)
				sDest += c.Word
			}
			i++
		}
		for i <= srcLen {
			sDest += fmt.Sprintf("${%d}", i)
			i++
		}
	}

	// aa*bb -> xx // `xx`
	// aa*bb -> zz* // `zz${1}

	return re.ReplaceAllString(filename, sDest), nil
}
*/

// compile search pattern into a regexp
func compilePattern(patternStr string) (re *regexp.Regexp, err error) {
	pattern := parseSearchPattern(patternStr)
	rSrc := ""
	srcLen := 0
	if len(pattern) == 0 {
		rSrc = "^.*$"
	} else {
		//rSrc = "^(.*?)" // rSrc[0] is left-side all-matching
		rSrc = "^"
		srcLen++
		for _, c := range pattern {
			switch c.Type {

			case '*':
				rSrc += "(.*?)"

			case '?':
				rSrc += "(" + strings.Replace(c.Word, "?", ".", -1) + ")"

			case ':': // digits
				rSrc += "(\\d+)"

			case 0:
				rSrc += "(" + regexp.QuoteMeta(c.Word) + ")"
			}
			srcLen++
		}
		rSrc += "$"
	}
	return regexp.Compile(rSrc)
}

// rename simple search pattern to replace pattern
func ReplaceName(filename string, pattern *regexp.Regexp, replace []Pattern, leftmostOffset int) (newname string, err error) {
	mt := pattern.FindStringSubmatch(filename)
	if mt == nil {
		return filename, nil
	}
	srcLen := len(mt)
	//log.Printf("matching: %v", mt)
	//log.Printf("replace: %v", replace)

	buf := ""
	emitLastMatch := false
	if len(replace) == 0 {
		return filename, nil
	} else {
		//s = mt[1]
		i := 1 // index of regexp submatch
		for _, p := range replace {
			switch p.Type {
			case '?', '*': // copy the matching entry
				if i >= srcLen {
					if !emitLastMatch {
						// special case: if the last match is not printed, then print it
						buf += mt[srcLen-1]
						emitLastMatch = true
					} else {
						// no entry in that position
						//log.Printf("a")
						return "", ErrInvalidPattern
					}
				} else {
					// TODO: check length for '?' match
					buf += mt[i]
					if i == srcLen-1 {
						emitLastMatch = true
					}
				}

			case 0:
				// copy the word
				buf += p.Word

			case '$': // position match
				/*
					n, err := strconv.Atoi(p.Word)
					if err != nil {
						//log.Printf("b")
						return "", ErrInvalidPattern
					}
				*/
				n := p.Pos + leftmostOffset
				if n < 1 || n >= srcLen {
					return "", ErrInvalidPosition
				}
				buf += mt[n]

			case '%': // printf conversion %...d or %...s
				l := len(p.Word)
				if l == 0 {
					return "", ErrInvalidPattern
				}
				k := p.Word[l-1]
				pos := p.Pos
				if pos == 0 {
					pos = i
				}
				pos += leftmostOffset
				if pos < 1 || len(mt) <= pos {
					return "", ErrInvalidPosition
				}
				switch k {
				case 'd': //
					n, err := strconv.Atoi(mt[pos])
					if err != nil {
						log.Printf("HI [%v]%v, offset: %d", pos, mt[pos], leftmostOffset)
						return "", ErrInvalidNumber
					}
					buf += fmt.Sprintf("%"+p.Word, n)
				case 's':
					buf += fmt.Sprintf("%"+p.Word, mt[pos])
				default:
					//log.Printf("d")
					return "", ErrInvalidPattern
				}

			default:
				//log.Printf("e")
				return "", ErrInvalidPattern

			}
			i++
		}
		//for i < srcLen {
		//	s += mt[i]
		//	i++
		//}
	}

	return buf, nil
}
