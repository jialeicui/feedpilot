package weibo

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func extractJson(body []byte) (prefix []byte, obj []byte, err error) {
	p := regexp.MustCompile("(.*?)\\(({.*}).*")
	match := p.FindSubmatch(body)
	if len(match) != 3 {
		return nil, nil, fmt.Errorf("can not get valid json with resp %q", string(body))
	}

	return match[1], match[2], nil
}

func cookiesToString(cookies map[string]string) string {
	var ret []string
	for k, v := range cookies {
		ret = append(ret, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(ret, "; ")
}

func (s *Session) updateCookies(cookies []*http.Cookie) {
	speCookie := map[string]struct{}{"SUB": {}, "SUBP": {}}
	for _, c := range cookies {
		if _, ok := s.baseCookies[c.Name]; ok {
			s.baseCookies[c.Name] = c.Value
		}
		if _, ok := speCookie[c.Name]; ok {
			if c.Domain == ".weibo.com" {
				s.authCookies[c.Name] = c.Value
			}
			continue
		}
		s.authCookies[c.Name] = c.Value
	}
}
