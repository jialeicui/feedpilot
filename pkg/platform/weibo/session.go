package weibo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"sync"

	"github.com/jialeicui/feedpilot/pkg/utils"
)

type Session struct {
	sync.RWMutex
	qrId        string
	baseCookies map[string]string
	authCookies map[string]string

	lastCtxCancel context.CancelFunc

	persistPath string

	Uid         string
	DisplayName string

	httpCli *utils.Client
}

func NewSession(persistPath string) *Session {
	s := &Session{
		// the cookie names we need for renew the auto cookies
		baseCookies: map[string]string{
			"SRF": "",
			"SRT": "",
			"tgc": "",
		},
		authCookies: map[string]string{},
		httpCli:     utils.NewWithHeader(nil),
		persistPath: persistPath,
	}
	if s.persistPath != "" {
		err := s.load()
		if err != nil {
			slog.Warn("failed to load session persist file", "path", s.persistPath, "err", err)
		}
	}
	return s
}

// RenewWithQrCode starts a new session and returns the qrcode image
func (s *Session) RenewWithQrCode() (image []byte, err error) {
	return s.renewWithQrCode()
}

// Wait waits for the scan action and get the necessary baseCookies & user info
func (s *Session) Wait(ctx context.Context) (cookies string, err error) {
	alt, err := s.checkScanState(ctx)
	if err != nil {
		return
	}
	err = s.fetchCookies(alt)
	if err != nil {
		return "", err
	}
	_ = s.dump()
	return s.Cookies(), nil
}

// Cookies returns the cookie string for auth usage
func (s *Session) Cookies() string {
	s.RLock()
	defer s.RUnlock()
	return cookiesToString(s.authCookies)
}

func (s *Session) Renew() (cookies string, err error) {
	s.Lock()
	defer s.Unlock()
	// check necessary cookies
	for _, v := range s.baseCookies {
		if v == "" {
			return "", fmt.Errorf("no valid base cookies, please login via qr code")
		}
	}
	url := fmt.Sprintf("https://passport.weibo.com/visitor/visitor?a=restore&cb=restore_back&from=weibo&_rand=%.16f", rand.Float32())
	cli := s.httpCli
	cli.SetHeader("referer", "https://weibo.com/")
	cli.SetHeader("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	cli.SetHeader("cookie", cookiesToString(s.baseCookies))
	content, cookie, err := cli.GetWithCookie(url)
	if err != nil {
		return "", err
	}

	_, content, err = extractJson(content)
	if err != nil {
		return "", err
	}
	s.updateCookies(cookie)

	type renewResp struct {
		RetCode int `json:"retcode"`
		Data    struct {
			Alt       string `json:"alt"`
			SaveState int    `json:"savestate"`
		} `json:"data"`
	}
	var resp = new(renewResp)
	err = json.Unmarshal(content, resp)
	if err != nil {
		return "", err
	}
	if resp.RetCode != 20000000 {
		return "", fmt.Errorf("try renew passport with ret %s", string(content))
	}

	url = fmt.Sprintf("https://login.sina.com.cn/sso/login.php?entry=sso&alt=%s&returntype=META&gateway=1&savestate=%d", resp.Data.Alt, resp.Data.SaveState)
	cli.SetHeader("cookie", cookiesToString(s.authCookies))
	content, cookie, err = cli.GetWithCookie(url)
	if err != nil {
		return "", err
	}

	s.updateCookies(cookie)
	_ = s.dump()
	return cookiesToString(s.authCookies), nil
}

func (s *Session) Dumps() ([]byte, error) {
	return json.Marshal(s.baseCookies)
}

func (s *Session) Loads(content []byte) error {
	return json.Unmarshal(content, &s.baseCookies)
}

func (s *Session) dump() (err error) {
	defer func() {
		if err != nil {
			slog.Warn("dump cookie failed", "err", err)
		}
	}()
	content, err := s.Dumps()
	if err != nil {
		return err
	}
	return os.WriteFile(s.persistPath, content, 0600)
}

func (s *Session) load() error {
	content, err := os.ReadFile(s.persistPath)
	if err != nil {
		return err
	}
	return s.Loads(content)
}
