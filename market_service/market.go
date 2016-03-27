package main

import (
	"fmt"
	"time"
	"github.com/nickstefan/market/market_service/heap"
	"bytes"
	"encoding/json"
	"net/http"
)

// how does a stock market organize the orders? Depth of Market or OrderBook

type OrderBook struct {
	handleTrade tradehandler
	buyQueue heap.Heap
	sellQueue heap.Heap
	buyHash map[string]*Order
	sellHash map[string]*Order
}

type tradehandler func(Trade)

func NewOrderBook() *OrderBook {
	return &OrderBook{
		handleTrade: func(t Trade) { },
		buyHash: make(map[string]*Order),
		sellHash: make(map[string]*Order),
		buyQueue: heap.Heap{Priority: "max"},
		sellQueue: heap.Heap{Priority: "min"},
	}
}

func (o *OrderBook) setTradeHandler(execTrade tradehandler) {
	o.handleTrade = execTrade
}

func (o *OrderBook) add(order Order) {

	if order.getOrder().intent == "BUY" {
		o.buyHash[order.lookup()] = &order
		o.buyQueue.Enqueue(&heap.Node{
			Value: order.price(),
			Lookup: order.lookup(),
		})

	} else if order.getOrder().intent == "SELL" {
		o.sellHash[order.lookup()] = &order
		o.sellQueue.Enqueue(&heap.Node{
			Value: order.price(),
			Lookup: order.lookup(),
		})
	}
}

func (o *OrderBook) run() {
	buyTop := o.buyQueue.Peek()
	sellTop := o.sellQueue.Peek()
	
	for (buyTop != nil && sellTop != nil && buyTop.Value >= sellTop.Value) {
		
		buy := *(o.buyHash[ buyTop.Lookup ])
		sell := *(o.sellHash[ sellTop.Lookup ])

		if buy.getOrder().shares == sell.getOrder().shares {
			o.buyQueue.Dequeue()
			o.sellQueue.Dequeue()

			price := buy.price()
			o.handleTrade(buy.fill(price))
			o.handleTrade(sell.fill(sell.price()))
			
			delete(o.buyHash, buyTop.Lookup)
			delete(o.sellHash, sellTop.Lookup)

		} else if buy.getOrder().shares < sell.getOrder().shares {
			o.buyQueue.Dequeue()
			remainderSell := sell.getOrder().shares - buy.getOrder().shares
			
			price := buy.price()
			o.handleTrade(buy.fill(price))
			o.handleTrade(sell.partialFill(sell.price(), remainderSell))
 
			delete(o.buyHash, buyTop.Lookup)
		
		} else if buy.getOrder().shares > sell.getOrder().shares {
			o.sellQueue.Dequeue()
			remainderBuy := buy.getOrder().shares - sell.getOrder().shares
			
			price := buy.price()
			o.handleTrade(sell.fill(sell.price()))
			o.handleTrade(buy.partialFill(price, remainderBuy))

			delete(o.sellHash, sellTop.Lookup)
		}
		
		buyTop = o.buyQueue.Peek()
		sellTop = o.sellQueue.Peek()
	}
}


type Trade struct {
	Actor string `json:"actor"`
	Shares int `json:"shares"`
	Ticker string `json:"ticker"`
	Price float64 `json:"price"`
	Intent string `json:"intent"`
	Kind string `json:"kind"`
	State  string `json:"state"`
}

func main() {

	orderBook := NewOrderBook()

	orderBook.setTradeHandler(func (t Trade) {
		fmt.Println("hello handler")
		url := "http://localhost:8000/fill"
		trade, err := json.Marshal(t)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(trade))
		if err != nil {
			panic(err)
		}
		fmt.Println("response Status:", resp.Status)
	})

	anOrder := SellLimit{
		ask: 10.05, 
		BaseOrder: &BaseOrder{
			actor: "Tim", timecreated: time.Now().Unix(),
			intent: "SELL", shares: 100, state: "OPEN", ticker: "STOCK",
			kind: "LIMIT",

		},
	}

	anotherOrder := BuyLimit{
		bid: 10.10, 
		BaseOrder: &BaseOrder{
			actor: "Bob", timecreated: time.Now().Unix(),
			intent: "BUY", shares: 100, state: "OPEN", ticker: "STOCK",
			kind: "LIMIT",
		},
	}
	orderBook.add(anOrder)
	orderBook.add(anotherOrder)
	orderBook.run()
}