// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	"io/ioutil"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
	"github.com/vargi/geopublisher/geohash"
	"fmt"
	"strconv"
	"net/url"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options


func getGeoHash(vars url.Values)(string){
	lat := vars.Get("lat")
	if !(len(lat) >  0) {
		return ""
	}
	lng := vars.Get("lng")
	if !(len(lng) >  0) {
		return ""
	}
	precision := vars.Get("precision")
	if !(len(precision) >  0) {
		return ""
	}
	channel := vars.Get("channel_id")
	if !(len(channel) >  0) {
		return ""
	}
	xlat, err := strconv.Atoi(lat)
	xlng, err := strconv.Atoi(lng)
	xprecision, err := strconv.Atoi(precision)
	if err != nil {
		return ""
	}

	geohash := geohash.EncodeWithPrecision(float64(xlat), float64(xlng), xprecision)
	geohash = channel + "|" + geohash
	return geohash
}




func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	redis_conn, err := redis.Dial("tcp", ":6379")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	vars := r.URL.Query()
	geohash := getGeoHash(vars)
	if !(len(geohash) > 0) {
		return
	}

	if err != nil {
		panic(err)
	}
	fmt.Println(geohash)
	redis_conn.Do("PUBLISH", geohash, string(body))
}


func socketHandler(w http.ResponseWriter, r *http.Request) {
	redis_conn , err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	psc := redis.PubSubConn{redis_conn}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	vars := r.URL.Query()

	geohash := getGeoHash(vars)
	if !(len(geohash) > 0) {
		return
	}
	geohash = geohash + "*"
	err = psc.PSubscribe(geohash)
	if err != nil {
		panic(err)
	}
	fmt.Println("hello")
	for {
	    switch v := psc.Receive().(type) {
	    case redis.PMessage:
	        fmt.Println("hello")
		fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		err = c.WriteMessage(1, v.Data)
	        if err != nil {
			panic(err)
		}

	    case redis.Subscription:
		fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
	    }
	}
}


func main() {
	flag.Parse()
	log.SetFlags(0)
	r := mux.NewRouter()
	r.HandleFunc("/socket", socketHandler)
	r.HandleFunc("/messages", postMessageHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", r))

}