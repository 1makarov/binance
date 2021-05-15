```go
func Bot() {
    client := &BinanceClient{
        Bkey:       setting.CFG.BKey,
        Bsecretkey: setting.CFG.BSecretKey,
        Timestamp:  5000,
    }

    order := &Order{
        Symbol:  "DOGEUSDT",
        Side:    SideBuy,
        Setting: OrderSetting {
             MarginType: ISOLATED,
             Leverage:   20,
             Type:       TypeMarket,
        },
    }

    if err := client.Bracket(order); err != nil {
        log.Fatal(err)
    }
    if err := client.FastSettingOrder(order); err != nil {
        log.Fatal(err)
    }
    if err := client.CreateOrder(order); err != nil {
        log.Fatal(err)
    }

    order.Side = SideSell

    if err := client.CreateOrder(order); err != nil {
        log.Fatal(err)
    }
}
```