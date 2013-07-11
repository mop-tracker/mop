// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
//=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
package mop
import (
        "strings"
        "os/user"
        "io/ioutil"
        "encoding/json"
)

const rcfile = "/.moprc"

type Profile struct {
        MarketRefreshRate int
        QuotesRefreshRate int
        Tickers           []string
        SortBy            string
        SortOrder         string
}

var profile Profile

//-----------------------------------------------------------------------------
func LoadProfile() string {
        data, err := ioutil.ReadFile(defaultProfile())
        if err != nil {
                // Set default values.
                profile.MarketRefreshRate = 12
                profile.QuotesRefreshRate = 5
                profile.Tickers = []string{ "AAPL", "C", "GOOG", "IBM", "KO", "ORCL", "V" }
                profile.SortBy = "Ticker"
                profile.SortOrder = "Desc"
                profile.Save()
        } else {
                json.Unmarshal(data, &profile)
        }
	return strings.Join(profile.Tickers, "+")
}

//-----------------------------------------------------------------------------
func (profile *Profile) Save() error {
        if data, err := json.Marshal(profile); err != nil {
                return err
        } else {
                return ioutil.WriteFile(defaultProfile(), data, 0644)
        }
}

//-----------------------------------------------------------------------------
func defaultProfile() string {
        usr, err := user.Current()
        if err != nil {
                panic(err)
        }
        return usr.HomeDir + rcfile
}
