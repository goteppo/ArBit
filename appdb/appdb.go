// Copyright 2011 Teppo Salonen. All rights reserved.
// This file is part of ArBit and distributed under the terms of the GNU LGPLv3 license.

// Package appdb implements basic functions for storing and retrieving data from Google App Engine's datastore.
package appdb

// TODO: |filterStr| and |filterVal| not used in |Query|

import (
	"appengine"
	"appengine/datastore"
	"os"
)

type UniqueKeyer interface {
	UniqueKey() (string, int64)
}

// Put creates or updates a datastore entity of given kind - the entity is stored in |data| and implements interface |UniqueKeyer|.
func Put(c appengine.Context, kind string, data UniqueKeyer) (err os.Error) {
	stringId, intId := data.UniqueKey()
	key := datastore.NewKey(kind, stringId, intId, nil)
	_, err = datastore.Put(c, key, data)
	return
}

// KeyPut creates or updates a datastore entity of given kind - the entity is stored in |data| and has a unique key either in |sKey| or in |iKey|.
func KeyPut(c appengine.Context, kind string, data interface{}, sKey string, iKey int64) (err os.Error) {
	key := datastore.NewKey(kind, sKey, iKey, nil)
	_, err = datastore.Put(c, key, data)
	return
}

// Get retrieves a datastore entity of given kind based on the unique key found in |data|, and stores the entity back to |data|.
func Get(c appengine.Context, kind string, data UniqueKeyer) (err os.Error) {
	stringId, intId := data.UniqueKey()
	key := datastore.NewKey(kind, stringId, intId, nil)
	err = datastore.Get(c, key, data)
	return
}

// KeyGet retrieves a datastore entity of given kind based on the unique key in either |sKey| or |iKey|, and stores the entity to |data|.
func KeyGet(c appengine.Context, kind string, data interface{}, sKey string, iKey int64) (err os.Error) {
	key := datastore.NewKey(kind, sKey, iKey, nil)
	err = datastore.Get(c, key, data)
	return
}

// Query retrieves datastore entities of given kind.
func Query(c appengine.Context, kind string, filterStr string, filterVal interface{},
order string, offset int, limit int) (data []datastore.Map, err os.Error) {
	q := datastore.NewQuery(kind)
	if order != "" {
		q = q.Order(order)
	}
	q = q.Offset(offset).Limit(limit)
	data = make([]datastore.Map, 0)
	_, err = q.GetAll(c, &data)
	return
}

// Delete deletes a datastore entity of given kind based on the unique key found in |data|.
func Delete(c appengine.Context, kind string, data UniqueKeyer) (err os.Error) {
	stringId, intId := data.UniqueKey()
	key := datastore.NewKey(kind, stringId, intId, nil)
	err = datastore.Delete(c, key)
	return
}
