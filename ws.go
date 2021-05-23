package binance

import (
	"github.com/gorilla/websocket"
	"log"
	"strconv"
)

type wsClient struct {
	session *websocket.Conn
}

func GetSingleMarkPrice(o *Order) {
	client := wsClient{}
	wssURL := wsapi + o.Symbol + markPriceSingle

	for {
		if err := client.wss(wssURL); err != nil {
			log.Fatal(err)
		}

		if err := client.singleMarkPrice(&o.Price); err != nil {
			log.Printf("error in getting the cost: %s\n", err.Error())
			continue
		}
	}

}

type singleMarkPriceR struct {
	Ee   string `json:"e"`
	E    int64  `json:"E"`
	S    string `json:"s"`
	Mark string `json:"p"`
	P    string `json:"P"`
	I    string `json:"i"`
	R    string `json:"r"`
	Time int64  `json:"T"`
}

func (w *wsClient) singleMarkPrice(price *float64) error {
	for {
		var amount singleMarkPriceR
		if err := w.session.ReadJSON(&amount); err != nil {
			return err
		}
		mPrice, err := strconv.ParseFloat(amount.Mark, 64)
		if err != nil {
			return err
		}
		*price = mPrice
	}
}
