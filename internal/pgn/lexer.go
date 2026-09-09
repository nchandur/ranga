package pgn

import (
	"fmt"
	"regexp"
	"strings"
)

type moveComment struct {
	san     string
	comment string
}

// walks PGN movetext and returns each SAN move paired with the comment that immediately follows it
func tokenizeWithComments(movetext string) ([]moveComment, error) {
	var pairs []moveComment
	lastIdx := -1
	n := len(movetext)
	i := 0

	for i < n {
		c := movetext[i]

		switch c {
		case ' ', '\n', '\t', '\r':
			i++

		case '(':
			depth := 1
			i++
			for i < n && depth > 0 {
				switch movetext[i] {
				case '(':
					depth++
				case ')':
					depth--
				}
				i++
			}

		case '{':
			j := i + 1
			for j < n && movetext[j] != '}' {
				j++
			}
			if j >= n {
				return nil, fmt.Errorf("unterminated comment starting at byte %d", i)
			}
			if lastIdx >= 0 {
				pairs[lastIdx].comment = movetext[i+1 : j]
			}
			i = j + 1

		default:
			j := i
			for j < n {
				ch := movetext[j]
				if ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r' || ch == '{' || ch == '(' {
					break
				}
				j++
			}
			word := movetext[i:j]
			i = j

			switch {
			case isResultToken(word), isMoveNumber(word), strings.HasPrefix(word, "$"):
			default:
				pairs = append(pairs, moveComment{san: word})
				lastIdx = len(pairs) - 1
			}
		}
	}

	return pairs, nil
}

func isResultToken(s string) bool {
	switch s {
	case "1-0", "0-1", "1/2-1/2", "*":
		return true
	}
	return false
}

func isMoveNumber(s string) bool {
	i := strings.IndexFunc(s, func(r rune) bool { return r < '0' || r > '9' })
	if i == -1 {
		return false // pure digits
	}
	return strings.Trim(s[i:], ".") == ""
}

// pulls the leading signed eval number out of an engine comment,
var scoreRe = regexp.MustCompile(`^([+-]?\d+(?:\.\d+)?)`)

func extractScore(comment string) string {
	comment = strings.TrimSpace(comment)
	m := scoreRe.FindStringSubmatch(comment)
	if m == nil {
		return ""
	}
	return m[1]
}
