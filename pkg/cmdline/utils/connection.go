package connection

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var Host string

func Websocket() (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "ws",
		Host:   Host,
		Path:   "/adm",
	}

	conn, res, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusSwitchingProtocols {
		return nil, errors.New(res.Status)
	}

	return conn, nil
}

func Status() (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://%s/adm/status", Host))
}

func Servers() (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://%s/adm/servers", Host))
}
