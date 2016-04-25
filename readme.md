# Stock Market Clone

_now in golang!_

### TODO
- [x] ledger service  
  - [x] listen for http (from orderbook)  
  - [x] track quantity and asset for each user (cash is an asset)  

- [x] orderbook service  
  - [x] listen for http (for submitting of orders)  
  - [x] priotiy queues for buy and sell orders  
  - [x] dequeue priority queues into trades  
  - [x] message to ledger service  
  - [x] message to ticker service  


- [ ] ticker service (chart and quote stream)  
  - [x] on price data http, update minute high, low, open, close, vol information  
  - [x] every 60 seconds, persist tick data to DB  
  - [ ] on price data, publish price to "QUOTE" channel; rate limit 1 / second  
  - [ ] every 60 seconds, publish tick data to "CHART" channel  

- [ ] web service  
  - [ ] serve front end javascript  
  - [ ] handle accounts and authentication  

- [ ] web client (front end)  
  - [ ] graph "CHART" data into a stock chart  
  - [ ] listen to "QUOTE" and "CHART" channels  
  - [ ] append new data to chart  
  - [ ] display account info  
  - [ ] place orders  

- [ ] bot service (automate trades)  
  - [ ] trading bots   
  - [ ] front end UI for creating bots  
  - [ ] rule based trading using JSON {$when: ..., $buy: ... }  

- [ ] nginx reverse proxy all services to one domain  
- [ ] message queue between services  

Ideas:
 - Ticker names: Kryptonite, Adamantium, Puppies
 - rule based trading:

    {
      $when: {
        $comparison: {
          itemA: { $ticker: "STOCKA"},
          itemB: 20.55,
          type: "gte"
        }
      },
      $do: {
        $order: { kind: "LIMIT", intent: "BUY", bid: 20.60, shares: 100 }
      },
      $then: {
        $when: {
          $comparison: {
            itemA: { $ticker: "STOCKA"},
            itemB: 21.00,
            type: "gte"
          }
        },
        $do: {
          $order: { kind: "LIMIT", intent: "SELL", ask: 20.90, shares 100 }
        }
      }
    }
