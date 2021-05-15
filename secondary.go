package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

func timestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func (data *Client) sign(body *url.Values) string {
	body.Add("timestamp", fmt.Sprintf("%d", timestamp()))

	mac := hmac.New(sha256.New, []byte(data.Bsecretkey))
	mac.Write([]byte(body.Encode()))
	signature := hex.EncodeToString(mac.Sum(nil))

	requestValue := body.Encode() + fmt.Sprintf("&signature=%s", signature)
	body.Set("signature", signature)
	return requestValue
}