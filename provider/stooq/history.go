package stooq

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type DailyRecord struct {
	Date   string // YYYYMMDD
	High   float64
	Low    float64
	Volume float64
}

type HistoryData struct {
	Records     []DailyRecord
	High52      string
	Low52       string
	AvgVolume   string
	LastUpdated time.Time
}

type HistoryCache struct {
	mu       sync.RWMutex
	Data     map[string]HistoryData
	filePath string
}

func NewHistoryCache() *HistoryCache {
	usr, _ := user.Current()
	path := filepath.Join(usr.HomeDir, ".mop_history.json")

	c := &HistoryCache{
		Data:     make(map[string]HistoryData),
		filePath: path,
	}
	c.Load()
	return c
}

func (c *HistoryCache) Load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.filePath)
	if err == nil {
		json.Unmarshal(data, &c.Data)
	}
}

func (c *HistoryCache) Save() {
	c.mu.RLock()
	data, err := json.MarshalIndent(c.Data, "", "  ")
	c.mu.RUnlock()

	if err == nil {
		os.WriteFile(c.filePath, data, 0644)
	}
}

func (c *HistoryCache) Get(ticker string) (HistoryData, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.Data[ticker]
	return d, ok
}

func (c *HistoryCache) Set(ticker string, d HistoryData) {
	c.mu.Lock()
	c.Data[ticker] = d
	c.mu.Unlock()
	c.Save()
}

// Background worker
// Background worker
func (c *HistoryCache) UpdateInBackground(ctx context.Context, tickers []string, onUpdate func(ticker string)) {
	go func() {
		for _, t := range tickers {
			select {
			case <-ctx.Done():
				return
			default:
			}

			c.mu.RLock()
			cached, exists := c.Data[t]
			c.mu.RUnlock()

			// Check if update is needed (older than 24h)
			if exists && time.Since(cached.LastUpdated) < 24*time.Hour {
				continue
			}

			// Throttling
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}

			// Determine date range
			var d1, d2 string
			now := time.Now()
			oneYearAgo := now.AddDate(-1, 0, 0)
			d2 = now.Format("20060102")

			if exists && len(cached.Records) > 0 {
				// Incremental fetch: start from the day after the last record
				lastDate, _ := time.Parse("20060102", cached.Records[len(cached.Records)-1].Date)
				d1 = lastDate.AddDate(0, 0, 1).Format("20060102")
			} else {
				// Full fetch
				d1 = oneYearAgo.Format("20060102")
			}

			// Ensure d1 <= d2 (if up to date, skip fetch but update timestamp)
			if d1 > d2 {
				// Just update timestamp
				c.Set(t, cached)
				continue
			}

			newRecords, err := fetchHistory(t, d1, d2)
			if err == nil {
				// Merge records
				if exists {
					// Append new records
					cached.Records = append(cached.Records, newRecords...)
					// Prune old records (older than 1 year)
					cutoff := oneYearAgo.Format("20060102")
					validStart := 0
					for i, r := range cached.Records {
						if r.Date >= cutoff {
							validStart = i
							break
						}
					}
					cached.Records = cached.Records[validStart:]
				} else {
					cached.Records = newRecords
				}

				// Recalculate stats
				var maxHigh, minLow float64
				var totalVol float64
				minLow = math.MaxFloat64

				for _, r := range cached.Records {
					if r.High > maxHigh {
						maxHigh = r.High
					}
					if r.Low < minLow {
						minLow = r.Low
					}
					totalVol += r.Volume
				}

				if len(cached.Records) > 0 {
					cached.High52 = fmt.Sprintf("%.2f", maxHigh)
					cached.Low52 = fmt.Sprintf("%.2f", minLow)
					cached.AvgVolume = formatVolume(totalVol / float64(len(cached.Records)))
				} else {
					cached.High52 = "-"
					cached.Low52 = "-"
					cached.AvgVolume = "-"
				}

				cached.LastUpdated = time.Now()
				c.Set(t, cached)

				if onUpdate != nil {
					onUpdate(t)
				}
			}
		}
	}()
}

func fetchHistory(ticker, d1, d2 string) ([]DailyRecord, error) {
	// Handle implicit suffix for URL if missing
	reqTicker := ticker
	if !utils_hasSuffix(ticker) {
		reqTicker = ticker + ".US"
	}

	url := fmt.Sprintf("https://stooq.com/q/d/l/?s=%s&d1=%s&d2=%s&i=d", reqTicker, d1, d2)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	reader := csv.NewReader(resp.Body)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		// Possibly just no data for this range (e.g. weekend incremental fetch)
		return []DailyRecord{}, nil
	}

	var parsed []DailyRecord

	// Skip header
	for i, row := range records {
		if i == 0 || len(row) < 6 {
			continue
		}

		high, err1 := strconv.ParseFloat(row[2], 64)
		low, err2 := strconv.ParseFloat(row[3], 64)
		vol, err3 := strconv.ParseFloat(row[5], 64)
		date := row[0] // YYYYMMDD

		if err1 == nil && err2 == nil && err3 == nil {
			parsed = append(parsed, DailyRecord{
				Date:   date,
				High:   high,
				Low:    low,
				Volume: vol,
			})
		}
	}

	return parsed, nil
}

// simple helper to avoid circular dependency or duplication
func utils_hasSuffix(t string) bool {
	return len(t) > 3 && t[len(t)-3] == '.' || len(t) > 4 && t[len(t)-4] == '.'
}
