// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the MIT/X11 license.

// Package restapi implements basic functions to be used with RESTful API's returning JSON responses.
package restapi

import (
	"appengine"
	"appengine/urlfetch"
	"os"
	"io/ioutil"
	"json"
	"http"
)

// GetJson fetches a URL using a GET request, parses the returned JSON response, and stores the response in struct |a|.
func GetJson(c appengine.Context, url string, a interface{}) (err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	b := getUrl(c, url)
	err = json.Unmarshal(b, &a)
	return
}

// PostJson submits data to a URL using a POST request, parses the returned JSON response, and stores the response in struct |a|.
func PostJson(c appengine.Context, url string, data map[string][]string, a interface{}) (err os.Error) {
	defer func() {
		if e, ok := recover().(os.Error); ok {
			err = e
		}
	}()
	b := postForm(c, url, data)
	err = json.Unmarshal(b, &a)
	return
}

func check(err os.Error) {
	if err != nil {
		panic(err)
	}
}

func getUrl(c appengine.Context, url string) []byte {
	/*client := urlfetch.Client(c)
	res, err := client.Get(url)
	check(err)*/
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	t := &urlfetch.Transport{c, 10, true} // Set timeout to 10 seconds (instead of the default 5)
    res, err := t.RoundTrip(req)
    check(err)
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	check(err)
	return b
}

func postForm(c appengine.Context, url string, data map[string][]string) []byte {
	client := urlfetch.Client(c)
	res, err := client.PostForm(url, data)
	check(err)
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	check(err)
	return b
}
