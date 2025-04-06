package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	http.Client

	header http.Header
}

type Cookies = []*http.Cookie

func NewWithHeader(header map[string]string) (cli *Client) {
	h := http.Header{}
	for k, v := range header {
		h.Set(k, v)
	}

	cli = &Client{
		Client: http.Client{},
		header: h,
	}
	return
}

func (h *Client) Post(url string, data string) ([]byte, error) {
	content, _, err := h.do("POST", url, []byte(data))
	return content, err
}

func (h *Client) do(method, url string, data []byte) (content []byte, cookies Cookies, err error) {
	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header = h.header

	resp, err := h.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("fail with status code %d, url %q", resp.StatusCode, url)
	}
	defer resp.Body.Close()
	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	cookies = resp.Cookies()
	return
}

func (h *Client) Get(url string) ([]byte, error) {
	content, _, err := h.do("GET", url, nil)
	return content, err
}

func (h *Client) GetWithCookie(url string) ([]byte, Cookies, error) {
	return h.do("GET", url, nil)
}

func (h *Client) SetHeader(k, v string) {
	h.header.Set(k, v)
}
