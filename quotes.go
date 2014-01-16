package mop

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
)

type Quoter interface {
	Fetch()
	//parse([]byte) []Stock
	//Initialize(*Market, *Profile) Quoter
	Ok() (bool, string)
	AddTickers([]string) (int, error)
	RemoveTickers([]string) (int, error)
	isReady() bool
	GetProfile() *Profile
	GetStocks() []Stock
	SetStocks([]Stock)
}

type Stocker interface {
	SetAdvance(bool)
	GetAdvance() bool
}

type Stock struct {
	Ticker     string
	LastTrade  string // l1: last trade.
	Change     string // c6: change real time.
	ChangePct  string // k2: percent change real time.
	Open       string // o: market open price.
	Low        string // g: day's low.
	High       string // h: day's high.
	Low52      string // j: 52-weeks low.
	High52     string // k: 52-weeks high.
	Volume     string // v: volume.
	AvgVolume  string // a2: average volume.
	PeRatio    string // r2: P/E ration real time.
	PeRatioX   string // r: P/E ration (fallback when real time is N/A).
	Dividend   string // d: dividend.
	Yield      string // y: dividend yield.
	MarketCap  string // j3: market cap real time.
	MarketCapX string // j1: market cap (fallback when real time is N/A).
	Advancing  bool   // True when change is >= $0.
}

type Quotes struct {
	market  *Market  // Pointer to Market.
	profile *Profile // Pointer to Profile.
	stocks  []Stock  // Slice of stock quote data.
	errors  string   // Error string if any.
}

// Ok returns two values: 1) boolean indicating whether the error has occured,
// and 2) the error text itself.
func (quotes Quotes) Ok() (bool, string) {
	ok, err := quotes.errors == ``, quotes.errors
	v := reflect.ValueOf(quotes).FieldByName("profile")
	if ok && len(quotes.GetStocks()) == 0 && !v.IsNil() {
		ok, err = false, "No stocks found"
	}
	return ok, err
}

// AddTickers saves the list of tickers and refreshes the stock
// data if new tickers have been added. The function gets called
// from the line editor when user adds new stock tickers.
func (quotes Quotes) AddTickers(tickers []string) (added int, err error) {
	if added, err = quotes.GetProfile().AddTickers(tickers); err == nil && added > 0 {
		quotes.stocks = nil // Force fetch.
	}
	return
}

// RemoveTickers saves the list of tickers and refreshes the
// stock data if some tickers have been removed. The function
// gets called from the line editor when user removes existing
// stock tickers.
func (quotes Quotes) RemoveTickers(tickers []string) (removed int, err error) {
	if removed, err = quotes.GetProfile().RemoveTickers(tickers); err == nil && removed > 0 {
		quotes.stocks = nil // Force fetch.
	}
	return
}

// isReady returns true if we haven't fetched the quotes yet *or* the stock
// market is still open and we might want to grab the latest
// quotes. In both cases we make sure the list of requested
// tickers is not empty.
func (quotes Quotes) isReady() bool {
	return ((quotes.stocks != nil && len(quotes.stocks) == 0) || !quotes.market.IsClosed) // && len(quotes.GetProfile().Tickers) > 0
}

// Initialize ensures sane initial values for a Quotes struct.
// It returns a pointer to the new struct created.
func NewQuotes(market *Market, profile *Profile) *Quotes {
	if reflect.ValueOf(market).IsNil() {
		panic("Nil market found")
	}
	if reflect.ValueOf(profile).IsNil() {
		panic("Nil market found")
	}
	return &Quotes{market, profile, []Stock{}, ""}
}

func (quotes *Quotes) Fetch() {
	fmt.Println("You must specify a Fetch method")
	os.Exit(1)
}

/*func (quotes Quotes) parse(body []byte) []Stock {
	return []Stock{}
}*/
func (quotes *Quotes) GetProfile() *Profile {
	return quotes.profile
}
func (quotes Quotes) GetStocks() []Stock {
	return quotes.stocks
}
func (quotes Quotes) SetStocks([]Stock) {
	return
}
func (quotes Stock) SetAdvance(b bool) {
	quotes.Advancing = b
	return
}

func (quotes Stock) GetAdvance() bool {
	return quotes.Advancing
}

//------------------------------------------------------------------
func sanitize(body []byte) []byte {
	return bytes.Replace(bytes.TrimSpace(body), []byte{'"'}, []byte{}, -1)
}
