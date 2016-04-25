package main

import (
	"time"
 	"strconv"
 	// "fmt"
)

// what makes up a stock market? competing buy and sell orders

// market buy order "i will buy 100 shares at the lowest available sell price"
// limit buy order "i will buy 100 shares at the lowest available price below ____"

// market sell order "i will sell 100 shares at the highest available price"
// limit sell order "i will sell 100 shares highest available price above _____"

type Order struct {
	Bid float64 `json:"bid"`
	Ask float64 `json:"ask"`
	Shares int `json:"shares"`
	Ticker string `json:"ticker"` 						// STOCK
	Actor string `json:"actor"`   						// BOB
	Intent string `json:"intent"` 						// BUY || SELL
	Kind string `json:"kind"`									// MARKET || LIMIT
	State string `json:"state"`   						// OPEN || FILLED || CANCELED
	Timecreated int64 `json:"timecreated"` 		// unix time
}

func (o *Order) lookup() string {
	return o.Actor + strconv.FormatInt(o.Timecreated, 10)
}

func (o *Order) partialFill(price float64, newShares int) Trade {
	o.Shares = newShares
	return Trade{
		Actor: o.Actor, Shares: o.Shares - newShares,
		Price: price, Intent: o.Intent, Kind: o.Kind, Ticker: o.Ticker,
		Time: time.Now().Unix(),
	}
}

func (o *Order) fill(price float64) Trade {
	return Trade{
		Actor: o.Actor, Shares: o.Shares, Price: price,
		Intent: o.Intent, Ticker: o.Ticker, Kind: o.Kind,
		Time: time.Now().Unix(),
	}
}

func (o *Order) price() float64 {
	K := o.Kind
	I := o.Intent

	if I == "BUY" && K == "LIMIT" {
		return o.Bid

	} else if I == "SELL" && K == "LIMIT" {
		return o.Ask

	} else if I == "BUY" && K == "MARKET" {
		return 1000000.00

	} else if I == "SELL" && K == "MARKET" {
		return 0.00
	} 
	// should never get here
	return 1000000.00
}

type OrderHash struct {
	data map[string]*Order
}

func (o *OrderHash) get(key string) *Order {
	return o.data[key]
}

func (o *OrderHash) set(key string, order *Order){
	o.data[key] = order
}

func (o *OrderHash) remove(key string){
	delete(o.data, key)
}

func NewOrderHash() *OrderHash {
	return &OrderHash{
		data: make(map[string]*Order),
	}
}