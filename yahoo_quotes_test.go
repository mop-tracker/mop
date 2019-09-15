// Copyright (c) 2013-2016 by Michael Dvorkin. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuotes(t *testing.T) {
	market := NewMarket()
	profile := NewProfile()

	profile.Tickers = []string{"GOOG", "BA"}

	quotes := NewQuotes(market, profile)
	require.NotNil(t, quotes)

	data, err := ioutil.ReadFile("./yahoo_quotes_sample.json")
	require.Nil(t, err)
	require.NotNil(t, data)

	require.True(t, quotes.isReady())
	//quotes.Fetch(data)
	_, err = quotes.parse2(data)
	assert.NoError(t, err)

	require.Equal(t, 2, len(quotes.stocks))
	assert.Equal(t, "BA", quotes.stocks[0].Ticker)
	assert.Equal(t, "331.76", quotes.stocks[0].LastTrade)
	assert.Equal(t, "GOOG", quotes.stocks[1].Ticker)
	assert.Equal(t, "1214.38", quotes.stocks[1].LastTrade)
}
