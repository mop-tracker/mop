package stooq

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mop-tracker/mop/provider"
)

type Quotes struct {
	ctx          context.Context
	market       *Market
	profile      provider.Profile
	stocks       []provider.Stock
	errors       string
	lastFetch    time.Time
	historyCache *HistoryCache

	onUpdate func()
	mu       sync.Mutex
}

func NewQuotes(ctx context.Context, market *Market, profile provider.Profile) *Quotes {
	return &Quotes{
		ctx:          ctx,
		market:       market,
		profile:      profile,
		stocks:       make([]provider.Stock, 0),
		historyCache: NewHistoryCache(),
	}
}

func (q *Quotes) Fetch() provider.Quotes {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Simple caching: if less than 5m since last fetch, return current state
	if time.Since(q.lastFetch) < 5*time.Minute && len(q.stocks) > 0 {
		return q
	}

	q.market.Throttle()
	q.errors = ""
	tickers := q.profile.GetTickers()
	if len(tickers) == 0 {
		return q
	}

	// Prepare list of tickers including aux fundamentals
	var allTickers []string
	// we keep track of which tickers are "stocks" (not indices)
	// so we know which ones to look up basics for.
	stockSet := make(map[string]bool)

	for _, t := range tickers {
		// Heuristic: Indices usually start with ^.
		if !strings.HasPrefix(t, "^") {
			stockSet[t] = true

			// Decide on suffix
			base := t
			suffix := ".US"
			hasSuffix := false

			if strings.HasSuffix(t, ".US") {
				base = strings.TrimSuffix(t, ".US")
				hasSuffix = true
			} else if strings.Contains(t, ".") {
				// E.g. "VOD.UK". Suffix is ".UK"
				parts := strings.Split(t, ".")
				base = parts[0]
				suffix = "." + parts[1]
				hasSuffix = true
			}

			// If no suffix provided by user, append .US for the REQUEST
			reqTicker := t
			if !hasSuffix {
				reqTicker = t + ".US"
			}
			allTickers = append(allTickers, reqTicker)

			// Append auxiliary tickers
			// Note: Stooq usually only supports _PE/_MV for US stocks.
			// We construct them from base + suffix.
			allTickers = append(allTickers, base+"_PE"+suffix)
			allTickers = append(allTickers, base+"_MV"+suffix)
		} else {
			// Indices or others, use as is
			allTickers = append(allTickers, t)
		}
	}

	joinedTickers := strings.Join(allTickers, "+")
	url := fmt.Sprintf(marketURL, joinedTickers)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		q.errors = fmt.Sprintf("Request creation failed: %s", err)
		return q
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		q.errors = fmt.Sprintf("Fetch failed: %s", err)
		return q
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		q.errors = fmt.Sprintf("Stooq returned status: %d", resp.StatusCode)
		return q
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		q.errors = fmt.Sprintf("CSV parse failed: %s", err)
		return q
	}

	// Index data by symbol for easy lookup
	dataMap := make(map[string][]string)
	for i, record := range records {
		if i == 0 || len(record) < 7 {
			continue
		}
		dataMap[record[0]] = record
	}

	q.stocks = make([]provider.Stock, 0, len(tickers))

	for _, t := range tickers {
		// Try to find the main record
		// Stooq sometimes modifies symbol name in response? reliable enough usually.
		rec, found := dataMap[t]
		if !found {
			// Try implicit .US
			if !strings.Contains(t, ".") {
				rec, found = dataMap[t+".US"]
			}
		}

		if !found {
			continue
		}

		last, _ := strconv.ParseFloat(rec[6], 64)
		open, _ := strconv.ParseFloat(rec[3], 64)
		high, _ := strconv.ParseFloat(rec[4], 64)
		low, _ := strconv.ParseFloat(rec[5], 64)
		vol, _ := strconv.ParseFloat(rec[7], 64)

		change := last - open
		changePct := 0.0
		if open != 0 {
			changePct = (change / open) * 100
		}

		// Look for fundamentals
		peVal := "-"
		mvVal := "-"

		// Reconstruct aux keys
		base := t
		suffix := ".US"
		if strings.HasSuffix(t, ".US") {
			base = strings.TrimSuffix(t, ".US")
		} else if strings.Contains(t, ".") {
			parts := strings.Split(t, ".")
			base = parts[0]
			suffix = "." + parts[1]
		} else {
			// default to .US for raw tickers like "IBM"
			suffix = ".US"
		}

		if peRec, ok := dataMap[base+"_PE"+suffix]; ok {
			// PE stored in Close usually
			if v, err := strconv.ParseFloat(peRec[6], 64); err == nil && v != 0 {
				peVal = fmt.Sprintf("%.2f", v)
			}
		}

		if mvRec, ok := dataMap[base+"_MV"+suffix]; ok {
			if v, err := strconv.ParseFloat(mvRec[6], 64); err == nil && v != 0 {
				// value is in millions
				if v >= 1000 {
					mvVal = fmt.Sprintf("%.2fB", v/1000)
				} else {
					mvVal = fmt.Sprintf("%.0fM", v)
				}
			}
		}

		s := provider.Stock{
			Ticker:     t,
			LastTrade:  fmt.Sprintf("%.2f", last),
			Change:     fmt.Sprintf("%+.2f", change),
			ChangePct:  fmt.Sprintf("%+.2f%%", changePct),
			Open:       fmt.Sprintf("%.2f", open),
			High:       fmt.Sprintf("%.2f", high),
			Low:        fmt.Sprintf("%.2f", low),
			Volume:     formatVolume(vol),
			Low52:      "-",
			High52:     "-",
			AvgVolume:  "-",
			PeRatio:    peVal,
			PeRatioX:   "-",
			Dividend:   "-",
			Yield:      "-",
			MarketCap:  mvVal,
			MarketCapX: "-",
			Currency:   guessCurrency(t),
		}

		if change > 0 {
			s.Direction = 1
		} else if change < 0 {
			s.Direction = -1
		}

		// Enrich with historical data if available
		if hData, ok := q.historyCache.Get(t); ok {
			s.Low52 = hData.Low52
			s.High52 = hData.High52
			s.AvgVolume = hData.AvgVolume

			// If today's high/low exceeds the cached 52-week range, adjust the display
			// to show the breakout in real-time.
			if h52, err := strconv.ParseFloat(hData.High52, 64); err == nil && high > h52 {
				s.High52 = fmt.Sprintf("%.2f", high)
			}
			if l52, err := strconv.ParseFloat(hData.Low52, 64); err == nil && low < l52 && low > 0 {
				s.Low52 = fmt.Sprintf("%.2f", low)
			}
		}

		q.stocks = append(q.stocks, s)
	}

	// Trigger background update for history
	q.historyCache.UpdateInBackground(q.ctx, tickers, func(t string) {
		q.mu.Lock()
		defer q.mu.Unlock()

		if hData, ok := q.historyCache.Get(t); ok {
			for i := range q.stocks {
				if q.stocks[i].Ticker == t {
					q.stocks[i].Low52 = hData.Low52
					q.stocks[i].High52 = hData.High52
					q.stocks[i].AvgVolume = hData.AvgVolume

					// Re-evaluate immediate breakout high/low
					curHigh, _ := strconv.ParseFloat(q.stocks[i].High, 64)
					curLow, _ := strconv.ParseFloat(q.stocks[i].Low, 64)
					h52, _ := strconv.ParseFloat(hData.High52, 64)
					l52, _ := strconv.ParseFloat(hData.High52, 64)

					if curHigh > h52 {
						q.stocks[i].High52 = q.stocks[i].High
					}
					if curLow < l52 && curLow > 0 {
						q.stocks[i].Low52 = q.stocks[i].Low
					}

					break
				}
			}
		}

		if q.onUpdate != nil {
			q.onUpdate()
		}
	})

	q.lastFetch = time.Now()
	q.market.fetchedQuotes = true
	return q
}

func (q *Quotes) Ok() (bool, string) {
	return q.errors == "", q.errors
}

func (q *Quotes) AddTickers(tickers []string) (int, error) {
	added, err := q.profile.AddTickers(tickers)
	if added > 0 {
		q.stocks = nil // Invalidate cache
		q.lastFetch = time.Time{}
	}
	return added, err
}

func (q *Quotes) RemoveTickers(tickers []string) (int, error) {
	removed, err := q.profile.RemoveTickers(tickers)
	if removed > 0 {
		q.stocks = nil // Invalidate cache
		q.lastFetch = time.Time{}
	}
	return removed, err
}

func (q *Quotes) GetStocks() []provider.Stock {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Returns a copy to be safe?
	// actually []Stock is a slice. If we append, we make a new one.
	// But modifying elements in place (which we do in background) is dangerous if caller is reading.
	// So we should return a copy.
	cpy := make([]provider.Stock, len(q.stocks))
	copy(cpy, q.stocks)
	return cpy
}

func (q *Quotes) RefreshAdvice() int {
	return 5
}

func (q *Quotes) BindOnUpdate(f func()) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.onUpdate = f
}

func formatVolume(v float64) string {
	switch {
	case v > 1e9:
		return fmt.Sprintf("%.2fB", v/1e9)
	case v > 1e6:
		return fmt.Sprintf("%.2fM", v/1e6)
	case v > 1e3:
		return fmt.Sprintf("%.2fk", v/1e3)
	default:
		return fmt.Sprintf("%.0f", v)
	}
}

func guessCurrency(ticker string) string {
	if strings.HasSuffix(ticker, ".US") {
		return "USD"
	}
	if strings.HasSuffix(ticker, ".UK") {
		return "GBP"
	}
	if strings.HasSuffix(ticker, ".DE") {
		return "EUR"
	}
	if strings.HasSuffix(ticker, ".JP") {
		return "JPY"
	}
	if strings.HasSuffix(ticker, ".HK") {
		return "HKD"
	}
	if strings.HasSuffix(ticker, ".PL") {
		return "PLN"
	}
	// Default to USD if no dot (implicit US)
	if !strings.Contains(ticker, ".") {
		return "USD"
	}
	return ""
}
