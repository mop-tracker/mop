// Copyright (c) 2013 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mop

import (
	`sort`
	`encoding/json`
	`io/ioutil`
	`os/user`
	`strings`
)

const moprc = `/.moprc`

type Profile struct {
	MarketRefresh	 int
	QuotesRefresh	 int
	Grouped          bool
	Tickers          []string
	SortColumn       int
	Ascending        bool
	selected_column	 int
}

//-----------------------------------------------------------------------------
func (self *Profile) Initialize() *Profile {
	data, err := ioutil.ReadFile(self.default_file_name())
	if err != nil {
		// Set default values.
		self.MarketRefresh = 12
		self.QuotesRefresh = 5
		self.Grouped = false
		self.Tickers = []string{`AAPL`, `C`, `GOOG`, `IBM`, `KO`, `ORCL`, `V`}
		self.SortColumn = 0
		self.Ascending = true
		self.Save()
	} else {
		json.Unmarshal(data, self)
	}
	self.selected_column = -1

	return self
}

//-----------------------------------------------------------------------------
func (self *Profile) Save() error {
	if data, err := json.Marshal(self); err != nil {
		return err
	} else {
		return ioutil.WriteFile(self.default_file_name(), data, 0644)
	}
}

//-----------------------------------------------------------------------------
func (self *Profile) ListOfTickers() string {
	return strings.Join(self.Tickers, `+`)
}

//-----------------------------------------------------------------------------
func (self *Profile) AddTickers(tickers []string) (added int, err error) {
	added = 0
	existing := make(map[string]bool)

	for _, ticker := range self.Tickers {
		existing[ticker] = true
	}

	for _, ticker := range tickers {
		if _, found := existing[ticker]; !found {
			self.Tickers = append(self.Tickers, ticker)
			added++
		}
	}
	sort.Strings(self.Tickers)
	err = self.Save()
	return
}

//-----------------------------------------------------------------------------
func (self *Profile) RemoveTickers(tickers []string) (removed int, err error) {
	removed = 0
	for _, ticker := range tickers {
		for i, existing := range self.Tickers {
			if ticker == existing { // Requested ticker is there: remove i-th slice item.
				self.Tickers = append(self.Tickers[:i], self.Tickers[i+1:]...)
				removed++
			}
		}
	}
	err = self.Save()
	return
}

//-----------------------------------------------------------------------------
func (self *Profile) Reorder() error {
	if self.selected_column == self.SortColumn {
		self.Ascending = !self.Ascending
	} else {
		self.SortColumn = self.selected_column
	}
	return self.Save()
}

//-----------------------------------------------------------------------------
func (self *Profile) Regroup() error {
	self.Grouped = !self.Grouped
	return self.Save()
}

// private
//-----------------------------------------------------------------------------
func (self *Profile) default_file_name() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir + moprc
}
