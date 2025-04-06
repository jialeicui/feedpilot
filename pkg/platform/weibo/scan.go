package weibo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (s *Session) renewWithQrCode() (image []byte, err error) {
	tm := time.Now().Unix() * 10000
	url := fmt.Sprintf("https://login.sina.com.cn/sso/qrcode/image?entry=sso&size=180&service_id=pc_protection&callback=STK_%d", tm)

	cli := s.httpCli
	cli.SetHeader("referer", "https://passport.weibo.com/")

	resp, err := cli.Get(url)
	if err != nil {
		return
	}

	_, content, err := extractJson(resp)
	if err != nil {
		return
	}

	type qrResp struct {
		Data struct {
			QRid  string `json:"qrid"`
			Image string `json:"image"`
		} `json:"data"`
	}
	var qr = new(qrResp)
	err = json.Unmarshal(content, qr)
	if err != nil {
		return nil, fmt.Errorf("unmarshal qr resp fail with err: %v, content %q", err, string(content))
	}

	image, err = cli.Get(qr.Data.Image)
	if err != nil {
		return nil, err
	}
	s.qrId = qr.Data.QRid

	return image, nil
}

func (s *Session) checkScanState(ctx context.Context) (alt string, err error) {
	tm := time.Now().Unix() * 10000
	url := fmt.Sprintf("https://login.sina.com.cn/sso/qrcode/check?entry=sso&qrid=%s&callback=STK_%d", s.qrId, tm)

	cli := s.httpCli
	cli.SetHeader("referer", "https://weibo.com/")

	type checkResp struct {
		RetCode int `json:"retcode"`
		Data    struct {
			Alt string `json:"alt"`
		} `json:"data"`
	}

	for {
		resp, err := cli.Get(url)
		if err != nil {
			return "", err
		}
		_, content, err := extractJson(resp)
		if err != nil {
			return "", err
		}
		var check = new(checkResp)
		err = json.Unmarshal(content, check)
		if err != nil {
			return "", fmt.Errorf("unmarshal check resp fail with err: %v, content: %q", err, string(content))
		}
		if check.RetCode == 20000000 {
			return check.Data.Alt, nil
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}

		time.Sleep(time.Second)
	}
}

func (s *Session) fetchCookies(alt string) (err error) {
	url := "https://login.sina.com.cn/sso/login.php"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	qs := map[string]string{
		"entry":       "weibo",
		"returntype":  "TEXT",
		"crossdomain": "1",
		"cdult":       "3",
		"domain":      "weibo.com",
		"alt":         alt,
		"savestate":   "30",
		"callback":    fmt.Sprintf("STK_%d", time.Now().Unix()*10000),
	}
	q := req.URL.Query()
	for k, v := range qs {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	s.updateCookies(resp.Cookies())

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return s.getWeiboCookie(content)
}

func (s *Session) getWeiboCookie(content []byte) (err error) {
	_, content, err = extractJson(content)
	if err != nil {
		return
	}

	type loginResp struct {
		RetCode            string   `json:"retcode"`
		CrossDomainUrlList []string `json:"crossDomainUrlList"`
		Nick               string   `json:"nick"`
		Uid                string   `json:"uid"`
	}

	var login = new(loginResp)
	err = json.Unmarshal(content, login)
	if err != nil {
		return
	}

	var url string
	for _, u := range login.CrossDomainUrlList {
		if strings.Contains(u, "wbsso") {
			url = u
		}
	}
	tm := time.Now().Unix() * 10000
	url = fmt.Sprintf("%s&action=login&callback=%d", url, tm)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Add("referer", "https://weibo.com/")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	s.updateCookies(resp.Cookies())

	defer resp.Body.Close()
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type userInfo struct {
		Result   bool `json:"result"`
		Userinfo struct {
			UniqueId    string `json:"uniqueid"`
			DisplayName string `json:"displayname"`
		} `json:"userinfo"`
	}

	_, content, err = extractJson(content)
	if err != nil {
		return err
	}

	var info = new(userInfo)
	err = json.Unmarshal(content, info)
	if err != nil {
		return err
	}
	s.DisplayName = info.Userinfo.DisplayName
	s.Uid = info.Userinfo.UniqueId

	return nil
}
