package binance

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/valyala/fasthttp"
	"net/url"
)

const (
	fapi  = "https://fapi.binance.com"
	wsapi = "wss://fstream.binance.com/ws/"
)

func (data *Client) post(url string, body *url.Values) (*fasthttp.Response, error) {
	data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(fapi + url)
	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.SetContentType("application/x-www-form-urlencoded")
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	request.SetBody([]byte(body.Encode()))
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(request, response)

	return response, err
}

func (data *Client) get(url string, body *url.Values) (*fasthttp.Response, error) {
	requestValue := data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(fapi + url + "?" + requestValue)
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(request, response)
	fmt.Println(string(response.Body()))

	return response, err
}

func (data *Client) delete(url string, body *url.Values) (*fasthttp.Response, error) {
	requestValue := data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(fapi + url + "?" + requestValue)
	request.Header.SetMethod(fasthttp.MethodDelete)
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	response := fasthttp.AcquireResponse()

	err := fasthttp.Do(request, response)

	return response, err
}

func (w *wsClient) wss(url string) error {
	dialer := new(websocket.Dialer)

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return err
	}

	w.session = conn

	return nil
}