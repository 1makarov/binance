package binance

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
	"net/url"
	"strconv"
)

type Client struct {
	Bkey       string // Binance Key
	Bsecretkey string // Binance Secret Key
	Timestamp  int64  // time
	RecvWindow int64  // max delay time
	Balance    FuturesBalances
}

const (
	CreateOrder      = "/fapi/v1/order"
	SymbolPrice      = "/fapi/v2/positionRisk"
	ChangeMarginType = "/fapi/v1/marginType"
	ChangeLeverage   = "/fapi/v1/leverage"
	GetBracket       = "/fapi/v1/leverageBracket"
	AccountInfo      = "/fapi/v2/account"

	ISOLATED   = "ISOLATED"
	CROSSED    = "CROSSED"
	SideBuy    = "BUY"
	SideSell   = "SELL"
	TypeMarket = "MARKET"
)

type RequestSuccess struct {
	Success bool `json:"success"`
}

type Order struct {
	Symbol   string
	Side     string // покупка/продажа
	Price    float64
	Quantity uint64 // количество
	Setting  OrderSetting
}

type OrderSetting struct {
	MarginType         string // ISOLATED, CROSSED
	AvailableAmountBuy uint64 // доступный баланс для покупки
	Leverage           uint64 // плечо
	Type               string // продажа лимитным ордером/маркетом
}

// open new order in futures
func (data *Client) CreateOrder(o *Order) error {
	body := &url.Values{
		"symbol":   {o.Symbol},
		"type":     {o.Setting.Type},
		"side":     {o.Side},
		"quantity": {strconv.FormatUint(o.Quantity, 10)},
	}

	response, err := data.post(CreateOrder, body)
	if err != nil || response.StatusCode() != fasthttp.StatusOK {
		return err
	}

	return nil
}

type FuturesBalances struct {
	USDT float64
	BUSD float64
	BNB  float64
}

// take account balance
func FuturesBalance(body *fasthttp.Response) (*FuturesBalances, error) {
	v, err := fastjson.ParseBytes(body.Body())
	if err != nil {
		return nil, err
	}
	array := v.GetArray("assets")

	AvailableBalance := &FuturesBalances{}
	for _, o := range array {
		t := string(o.GetStringBytes("asset"))
		b := o.GetStringBytes("availableBalance")
		balance, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			continue
		}
		switch t {
		case "BNB":
			AvailableBalance.BNB = balance
		case "USDT":
			AvailableBalance.USDT = balance
		case "BUSD":
			AvailableBalance.BUSD = balance
		}
	}
	return AvailableBalance, nil
}

// take price pair
func (data *Client) MarkPrice(o *Order) error {
	body := &url.Values{
		"symbol": {o.Symbol},
	}
	responce, err := data.get(SymbolPrice, body)
	if err != nil || responce.StatusCode() != fasthttp.StatusOK {
		return err
	}
	v, err := fastjson.ParseBytes(responce.Body())
	if err != nil {
		return err
	}
	vArray := v.GetArray()
	p := string(vArray[0].GetStringBytes("markPrice"))
	price, err := strconv.ParseFloat(p, 64)
	if err != nil {
		return err
	}
	o.Price = price
	return nil
}

type ChangeMargin struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// change margin type
func (data *Client) ChangeMarginType(o *Order) error {
	body := &url.Values{
		"symbol":     {o.Symbol},
		"marginType": {o.Setting.MarginType},
	}
	response, err := data.post(ChangeMarginType, body)
	if err != nil {
		return err
	}
	var r ChangeMargin
	if err = json.Unmarshal(response.Body(), &r); err != nil {
		return err
	}
	if !(r.Code == fasthttp.StatusOK || r.Code == -4046) {
		return fmt.Errorf("error change margin type %d, %s\n", r.Code, response.Body())
	}
	return nil
}

// change initial leverage
func (data *Client) ChangeLeverage(o *Order) error {
	body := &url.Values{
		"symbol":   {o.Symbol},
		"leverage": {strconv.FormatUint(o.Setting.Leverage, 10)},
	}
	response, err := data.post(ChangeLeverage, body)
	if err != nil || response.StatusCode() != fasthttp.StatusOK {
		return err
	}
	return nil
}

func (data *Client) AccountInfo() (*fasthttp.Response, error) {
	response, err := data.get(AccountInfo, &url.Values{})
	if err != nil || response.StatusCode() != fasthttp.StatusOK {
		return nil, err
	}
	return response, nil
}

func (data *Client) Bracket(o *Order) error {
	response, err := data.get(GetBracket, &url.Values{
		"symbol": {o.Symbol},
	})
	if err != nil || response.StatusCode() != fasthttp.StatusOK {
		return err
	}

	v, err := fastjson.ParseBytes(response.Body())
	if err != nil {
		return err
	}

	brackets := v.GetArray()[0].GetArray("brackets")
	for index := range brackets {
		leverage := brackets[index].GetUint64("initialLeverage")
		if index == 0 {
			if leverage < o.Setting.Leverage {
				o.Setting.AvailableAmountBuy = brackets[index].GetUint64("notionalCap")
				o.Setting.Leverage = leverage
				break
			}
		}
		if leverage == o.Setting.Leverage {
			o.Setting.AvailableAmountBuy = brackets[index].GetUint64("notionalCap")
			break
		}
		if leverage < o.Setting.Leverage {
			o.Setting.AvailableAmountBuy = brackets[index].GetUint64("notionalCap")
			break
		}
	}
	return nil
}

func (data *Client) FastSettingOrder(o *Order) error {
	accountInfo, err := data.AccountInfo()
	if err != nil {
		return err
	}
	balances, err := FuturesBalance(accountInfo)
	if err != nil {
		return err
	}
	data.Balance = *balances
	v, err := fastjson.ParseBytes(accountInfo.Body())
	if err != nil {
		return err
	}
	for _, p := range v.GetArray("positions") {
		symbol := string(p.GetStringBytes("symbol"))
		if symbol != o.Symbol {
			continue
		}
		leverage, err := strconv.ParseUint(string(p.GetStringBytes("leverage")), 10, 64)
		if err != nil {
			return err
		}
		if leverage != o.Setting.Leverage {
			if err = data.ChangeLeverage(o); err != nil {
				return err
			}
		}
		switch p.GetBool("isolated") {
		case true:
			if o.Setting.MarginType != ISOLATED {
				if err = data.ChangeMarginType(o); err != nil {
					return err
				}
			}
		case false:
			if o.Setting.MarginType != CROSSED {
				if err = data.ChangeMarginType(o); err != nil {
					return err
				}
			}
		}
		if err = data.MarkPrice(o); err != nil {
			return err
		}
		data.QuantitySearch(o, 0.95)
		return nil
	}
	return fmt.Errorf("symbol not found")
}

func (data *Client) QuantitySearch(o *Order, percent float64) {
	buyBalance := data.Balance.USDT * float64(o.Setting.Leverage)
	if buyBalance > float64(o.Setting.AvailableAmountBuy) {
		buyBalance = float64(o.Setting.AvailableAmountBuy)
	}
	o.Quantity = uint64(buyBalance / o.Price * percent)
}
