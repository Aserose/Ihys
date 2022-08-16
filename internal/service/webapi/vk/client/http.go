package client

import (
	"IhysBestowal/pkg/customLogger"
	"io"
	"net/http"
)

const (
	userAgentH = `User-Agent`
	userAgentV = `VKAndroidApp/4.13.1-1206 (Android 4.4.3; SDK 19; armeabi; ; ru)","Accept": "image/gif, image/x-xbitmap, image/jpeg, image/pjpeg, */*`
)

type Client struct {
	http *http.Client
	log  customLogger.Logger
}

func New(log customLogger.Logger) Client {
	return Client{
		http: &http.Client{},
		log:  log,
	}
}

func (h Client) Send(req *http.Request) []byte {
	req.Header.Set(userAgentH, userAgentV)

	b, err := h.http.Do(req)
	if err != nil {
		h.log.Warn(h.log.CallInfoStr(), err.Error())
	}
	body, err := io.ReadAll(b.Body)
	if err != nil {
		h.log.Warn(h.log.CallInfoStr(), err.Error())
	}

	return body
}
