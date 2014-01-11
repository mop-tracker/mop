// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	`encoding/json`
	`io/ioutil`
	`os/user`
	`sort`
)

// File name in user's home directory where we store the settings.
const moprc = `/.moprc`

// Profile manages Mop program settings as defined by user (ex. list of
// stock tickers). The settings are serialized using JSON and saved in
// the ~/.moprc file.
type Profile struct {
	Tickers        []string // List of stock tickers to display.
	MarketRefresh  int      // Time interval to refresh market data.
	QuotesRefresh  int      // Time interval to refresh stock quotes.
	SortColumn     int      // Column number by which we sort stock quotes.
	Ascending      bool     // True when sort order is ascending.
	Grouped        bool     // True when stocks are grouped by advancing/declining.
	selectedColumn int      // Stores selected column number when the column editor is active.
}

// Initialize attempts to load the settings from ~/.moprc file. If the
// file is not there it gets created with the default values.
func (profile *Profile) Initialize() *Profile {
	data, err := ioutil.ReadFile(profile.defaultFileName())
	if err != nil { // Set default values:
		profile.MarketRefresh = 12 // Market data gets fetched every 12s (5 times per minute).
		profile.QuotesRefresh = 5  // Stock quotes get updated every 5s (12 times per minute).
		profile.Grouped = false    // Stock quotes are *not* grouped by advancing/declining.
		profile.Tickers = []string{`AAPL`, `C`, `GOOG`, `IBM`, `KO`, `ORCL`, `V`}
		profile.SortColumn = 0   // Stock quotes are sorted by ticker name.
		profile.Ascending = true // A to Z.
		profile.Save()
	} else {
		json.Unmarshal(data, profile)
	}
	profile.selectedColumn = -1

	return profile
}

// Save serializes settings using JSON and saves them in ~/.moprc file.
func (profile *Profile) Save() error {
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(profile.defaultFileName(), data, 0644)
}

// AddTickers updates the list of existing tikers to add the new ones making
// sure there are no duplicates.
func (profile *Profile) AddTickers(tickers []string) (added int, err error) {
	added, err = 0, nil
	existing := make(map[string]bool)

	// Build a hash of existing tickers so we could look it up quickly.
	for _, ticker := range profile.Tickers {
		existing[ticker] = true
	}

	// Iterate over the list of new tickers excluding the ones that
	// already exist.
	for _, ticker := range tickers {
		if _, found := existing[ticker]; !found {
			profile.Tickers = append(profile.Tickers, ticker)
			added++
		}
	}

	if added > 0 {
		sort.Strings(profile.Tickers)
		err = profile.Save()
	}

	return
}

// RemoveTickers removes requested stock tickers from the list we track.
func (profile *Profile) RemoveTickers(tickers []string) (removed int, err error) {
	removed, err = 0, nil
	for _, ticker := range tickers {
		for i, existing := range profile.Tickers {
			if ticker == existing {
				// Requested ticker is there: remove i-th slice item.
				profile.Tickers = append(profile.Tickers[:i], profile.Tickers[i+1:]...)
				removed++
			}
		}
	}

	if removed > 0 {
		err = profile.Save()
	}

	return
}

// Reorder gets called by the column editor to either reverse sorting order
// for the current column, or to pick another sort column.
func (profile *Profile) Reorder() error {
	if profile.selectedColumn == profile.SortColumn {
		profile.Ascending = !profile.Ascending // Reverse sort order.
	} else {
		profile.SortColumn = profile.selectedColumn // Pick new sort column.
	}
	return profile.Save()
}

// Regroup flips the flag that controls whether the stock quotes are grouped
// by advancing/declining issues.
func (profile *Profile) Regroup() error {
	profile.Grouped = !profile.Grouped
	return profile.Save()
}

//-----------------------------------------------------------------------------
func (profile *Profile) defaultFileName() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir + moprc
}
