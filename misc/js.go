// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package main
import (
        "fmt"
        "encoding/json"
)

type Config struct {
        MarketRefreshRate int
        QuotesRefreshRate int
        Tickers           []string
        SortBy            string
        SortOrder         string
}

func main() {
        var cfg Config
        cfg.MarketRefreshRate = 1
        cfg.QuotesRefreshRate = 1
        cfg.Tickers = []string{ "AAPL", "ALU", "HPQ", "IBM" }
        cfg.SortBy = "Ticker"
        cfg.SortOrder = "Desc"
        fmt.Printf("%+v\n", cfg)
        blob, err := json.Marshal(cfg)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%q\n", blob)

        var cfg2 Config
        err = json.Unmarshal(blob, &cfg2)
        if err != nil {
                panic(err)
        }
        fmt.Printf("%+v\n", cfg2)
}
