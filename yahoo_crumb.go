// Copyright (c) 2013-2024 by Michael Dvorkin and contributors. All Rights Reserved.
// Use of this source code is governed by a MIT-style license that can
// be found in the LICENSE file.

package mop

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const crumbURL = "https://query1.finance.yahoo.com/v1/test/getcrumb"
const cookieURL = "https://login.yahoo.com"
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/113.0"

func fetchCrumb(cookies string) string {
	client := http.Client{}
	request, err := http.NewRequest("GET", crumbURL, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip, deflate, br"},
		"Accept-Language": {"en-US,en;q=0.5"},
		"Connection":      {"keep-alive"},
		"Content-Type":    {"text/plain"},
		"Cookie":          {cookies},
		"Host":            {"query1.finance.yahoo.com"},
		"Sec-Fetch-Dest":  {"empty"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Site":  {"same-site"},
		"TE":              {"trailers"},
		"User-Agent":      {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return string(body[:])
}

func fetchCookies() string {
	client := http.Client{}
	request, err := http.NewRequest("GET", cookieURL, nil)
	if err != nil {
		panic(err)
	}

	request.Header = http.Header{
		"Accept":                   {"*/*"},
		"Accept-Encoding":          {"gzip, deflate, br"},
		"Accept-Language":          {"en-US,en;q=0.5"},
		"Connection":               {"keep-alive"},
		"Host":                     {"login.yahoo.com"},
		"Sec-Fetch-Dest":           {"document"},
		"Sec-Fetch-Mode":           {"navigate"},
		"Sec-Fetch-Site":           {"none"},
		"Sec-Fetch-User":           {"?1"},
		"TE":                       {"trailers"},
		"Update-Insecure-Requests": {"1"},
		"User-Agent":               {userAgent},
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var result string
	for _, cookie := range response.Cookies() {
		if cookie.Name != "AS" {
			result += cookie.Name + "=" + cookie.Value + "; "
		}
	}
	result = strings.TrimSuffix(result, "; ")
	return result
}
