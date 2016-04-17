# geopublisher
Location based Go PubSub server 

## what is geopublisher?

Provides you a websocket server, that you may connect with your geo information (latitude, longitude, precision)
and receive the messages from people connected to same channel.

## How it works ? 

- Open a websocket connection to `ws://localhost:8080/socket?lat=47.12&lng=120.43&precision=5&channel_id=xyz`
- This will connect to channel with the **geohash** token generated with latitude longitude information like `<channel_id>|<geohash>`
- You will start to listen to changes in this channel and subgeohashes(for higher precision) and receive via websocket.
- You may send message by basic http request 
` POST {message:hello, user:tuna} http://localhost:8080/messages?lat=47.12&lng=120.43&precision=5&channel_id=xyz` 
- And the clients listening to those channels via websocket will receive your message

## Geohash

[You may check here to have a clue about geohashes](https://www.elastic.co/guide/en/elasticsearch/guide/current/geohashes.html)


## Requirements

- Redis
- Go

## How to run from source code?

- `mkdir /your/path`
- `export GOPATH=/your/path`
- `go get github.com/garyburd/redigo/redis`
- `go get github.com/gorilla/mux`
- `go get github.com/gorilla/websocket`
- `go get github.com/vargi/geopublisher`
- `cd /your/path/src/github.com/vargi/geopublisher/`
- `go run main.go`

## Reference

- Used [TomiHiltunen's geohash implementation](github.com/TomiHiltunen/geohash-golang) for generating geohashes

