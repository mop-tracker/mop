# Google Finance Data Provider Implementation Plan

This plan outlines the steps to add Google Finance as a data provider for the [mop](file:///c:/Users/speed/workspace/mop) application, enabling users to fetch market data and stock quotes from Google Finance as an alternative to Yahoo Finance.

## User Review Required

> [!IMPORTANT]
> Google Finance does not provide a public API. This implementation relies on scraping HTML and parsing embedded JSON logic (`AF_initDataCallback` / `window.WIZ_global_data`). This method is brittle and may break if Google changes their page structure.

> [!NOTE]
> Stock tickers on Google Finance often require an exchange prefix or suffix (e.g., `NASDAQ:AAPL` or `AAPL:NASDAQ`). The implementation will attempt to infer or handle standard suffixes, but users might need to update their `.moprc` if they have ambiguous tickers.

## Proposed Changes

### `provider/google` Package

Create a new package `provider/google` that implements `provider.Market` and `provider.Quotes`.

#### [NEW] [google_market.go](file:///c:/Users/speed/workspace/mop/provider/google/google_market.go)
- Implement [NewMarket() *Market](file:///c:/Users/speed/workspace/mop/provider/yahoo/yahoo_market.go#35-59)
- Implement [Fetch() provider.Market](file:///c:/Users/speed/workspace/mop/provider/yahoo/yahoo_market.go#60-107)
    - Scrape `https://www.google.com/finance/`
    - Extract index data (Dow, S&P 500, Nasdaq, etc.) using Regex or DOM parsing.
    - Map Google's data to `provider.MarketData`.

#### [NEW] [google_quotes.go](file:///c:/Users/speed/workspace/mop/provider/google/google_quotes.go)
- Implement [NewQuotes(market *Market, profile provider.Profile) *Quotes](file:///c:/Users/speed/workspace/mop/provider/yahoo/yahoo_quotes.go#43-51)
- Implement [Fetch() provider.Quotes](file:///c:/Users/speed/workspace/mop/provider/yahoo/yahoo_market.go#60-107)
    - Iterate through tickers in `profile`.
    - Fetch `https://www.google.com/finance/quote/{TICKER}:NASDAQ` (and fallback to NYSE or other if needed, or try clean search).
    - Parse specific fields (Price, Change, P/E, etc.) from the page.
    - Map to `provider.Stock`.

### `cmd/mop`

#### [MODIFY] [main.go](file:///c:/Users/speed/workspace/mop/cmd/mop/main.go)
- Add command-line flag `-provider` (default: "yahoo") or configuration option to switch between Yahoo and Google.
- Instantiate `google.NewMarket` and `google.NewQuotes` when Google is selected.
- Pass them to `mop.NewScreen` (which is already provider-agnostic).

## Verification Plan

### Automated Tests
- Create unit tests in `provider/google/parsing_test.go` (if feasible) to test the parsing logic against a saved HTML sample of Google Finance.
- Command: `go test ./provider/google/...`

### Manual Verification
1.  **Build**: `go build -o ./bin/mop ./cmd/mop`
2.  **Run with Yahoo (Default)**: `./bin/mop` -> Check if it still works.
3.  **Run with Google**: `./bin/mop -provider google`
    - Verify Market section (top 3 lines) shows Data.
    - Verify Quotes section shows data for default tickers (AAPL, C, GOOG, etc.).
    - Test adding a ticker: `+` then `AMZN`.
4.  **Edge Cases**:
    - Invalid ticker.
    - Network failure (disconnect internet).
