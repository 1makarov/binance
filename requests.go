package binance

import (
	"github.com/valyala/fasthttp"
	"net/url"
)

const binancefapi = "https://fapi.binance.com"

func (data *Client) post(url string, body *url.Values) (*fasthttp.Response, error) {
	data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(binancefapi + url)
	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.SetContentType("application/x-www-form-urlencoded")
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	request.SetBody([]byte(body.Encode()))
	responce := fasthttp.AcquireResponse()

	if err := fasthttp.Do(request, responce); err != nil {
		return responce, err
	}

	return responce, nil
}

func (data *Client) get(url string, body *url.Values) (*fasthttp.Response, error) {
	requestValue := data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(binancefapi + url + "?" + requestValue)
	request.Header.SetMethod(fasthttp.MethodGet)
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	responce := fasthttp.AcquireResponse()

	if err := fasthttp.Do(request, responce); err != nil {
		return nil, err
	}

	return responce, nil
}

func (data *Client) delete(url string, body *url.Values) (*fasthttp.Response, error) {
	requestValue := data.sign(body)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(binancefapi + url + "?" + requestValue)
	request.Header.SetMethod(fasthttp.MethodDelete)
	request.Header.Set("X-MBX-APIKEY", data.Bkey)
	responce := fasthttp.AcquireResponse()

	if err := fasthttp.Do(request, responce); err != nil {
		return nil, err
	}

	return responce, nil
}
