package stooq

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mop-tracker/mop/provider"
)

const (
	// Stooq uses standard symbols for major indices
	// ^DJI - Dow Jones
	// ^NDQ - Nasdaq 100 (or ^IC for Composite? Stooq usages vary, let's try ^NDQ or ^IC)
	// ^SPX - S&P 500
	// ^NKX - Nikkei 225
	// ^HSI - Hang Seng seemed to be supported?
	// ^DAX - DAX (German)
	// ^FTM - FTSE 100 (UK) - tentative
	// 10USY.B - 10 Year Treasury Yield
	marketSymbols = "^DJI ^NDQ ^SPX ^NKX ^HSI ^FTM ^DAX 10USY.B CL.F USDJPY EURUSD GC.F"
	marketURL     = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

type MarketIndex = provider.MarketIndex

type Market struct {
	provider.MarketData
	errors    string
	lastFetch time.Time
	cache     *provider.MarketData // Simple cache

	// Rate limiting key
	reqMu       sync.Mutex
	lastRequest time.Time

	// Synchronization flags
	fetchedMarket bool
	fetchedQuotes bool
}

func NewMarket() *Market {
	return &Market{
		MarketData: provider.MarketData{},
	}
}

func (market *Market) Fetch() provider.Market {
	// Simple throttling/caching
	if time.Since(market.lastFetch) < 5*time.Minute && market.cache != nil {
		market.MarketData = *market.cache
		return market
	}

	market.Throttle()
	market.errors = ""

	// Prepare symbols - space separated for Stooq usually means + in URL
	symbols := strings.ReplaceAll(marketSymbols, " ", "+")
	url := fmt.Sprintf(marketURL, symbols)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		market.errors = fmt.Sprintf("Request creation failed: %s", err)
		return market
	}

	// User-Agent spoofing
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		market.errors = fmt.Sprintf("Fetch failed: %s", err)
		return market
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		market.errors = fmt.Sprintf("Stooq returned status: %d", resp.StatusCode)
		return market
	}

	reader := csv.NewReader(resp.Body)
	// Stooq CSV format with header (Symbol,Date,Time,Open,High,Low,Close,Volume)
	records, err := reader.ReadAll()
	if err != nil {
		market.errors = fmt.Sprintf("CSV parse failed: %s", err)
		return market
	}

	if len(records) < 2 {
		market.errors = "No data returned from Stooq"
		return market
	}

	// Map symbols to MarketData fields
	// We iterate through records (skip header)
	for i, record := range records {
		if i == 0 {
			continue // Header
		}
		// Columns: Symbol(0), Date(1), Time(2), Open(3), High(4), Low(5), Close(6), Volume(7)
		if len(record) < 7 {
			continue
		}

		symbol := record[0]
		currentPrice, _ := strconv.ParseFloat(record[6], 64)
		openPrice, _ := strconv.ParseFloat(record[3], 64)

		// Approximate change since we don't have prev close
		// Ideally we would store prev close or Stooq would provide it.
		// For now, let's use Open vs Close as a rough intraday prox or 0 if N/A.
		// NOTE: Stooq often returns "N/A" for missing data in CSV which ParseFloat treats as error (0).

		change := currentPrice - openPrice
		percent := 0.0
		if openPrice != 0 {
			percent = (change / openPrice) * 100
		}

		idx := MarketIndex{
			Latest:  fmt.Sprintf("%.2f", currentPrice),
			Change:  fmt.Sprintf("%+.2f", change),
			Percent: fmt.Sprintf("%+.2f%%", percent),
		}

		switch symbol {
		case "^DJI":
			idx.Name = "Dow Jones"
			market.Dow = idx
		case "^NDQ":
			idx.Name = "Nasdaq"
			market.Nasdaq = idx
		case "^SPX":
			idx.Name = "S&P 500"
			market.Sp500 = idx
		case "^NKX":
			idx.Name = "Nikkei 225"
			market.Tokyo = idx
		case "^HSI":
			idx.Name = "Hang Seng"
			market.HongKong = idx
		case "^FTM":
			idx.Name = "FTSE 100"
			market.London = idx
		case "^DAX":
			idx.Name = "DAX"
			market.Frankfurt = idx
		case "10USY.B":
			idx.Name = "10y Yield"
			market.Yield = idx
		case "CL.F":
			idx.Name = "Crude Oil"
			market.Oil = idx
		case "USDJPY":
			idx.Name = "Yen"
			market.Yen = idx
		case "EURUSD":
			idx.Name = "Euro"
			market.Euro = idx
		case "GC.F":
			idx.Name = "Gold"
			market.Gold = idx
		}
	}

	market.lastFetch = time.Now()
	market.fetchedMarket = true
	// Cache the result
	cached := market.MarketData
	market.cache = &cached

	return market
}

func (market *Market) Ok() (bool, string) {
	if market.errors != "" {
		return false, market.errors
	}
	// Pause display until both are fetched
	if !market.fetchedMarket || !market.fetchedQuotes {
		return false, "<tag>Stooq:</> Waiting for quotes and market data to synchronize..."
	}
	return true, ""
}

func (market *Market) IsClosed() bool {
	return market.Closed
}

func (market *Market) GetData() *provider.MarketData {
	return &market.MarketData
}

func (market *Market) RefreshAdvice() int {
	return 5
}

func (market *Market) Throttle() {
	market.reqMu.Lock()
	defer market.reqMu.Unlock()

	since := time.Since(market.lastRequest)
	if since < 5*time.Second {
		time.Sleep(5*time.Second - since)
	}
	market.lastRequest = time.Now()
}
