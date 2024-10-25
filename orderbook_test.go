package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrder := NewOrder(true, 5)
	buyOrderB := NewOrder(true, 8)
	buyOrderC := NewOrder(true, 10)

	l.AddOrder(buyOrder)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrder)
	fmt.Println(l)
}
func TestPlaceLimitOrder(t *testing.T) {

	ob := NewOrderBook()

	sellOrder := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)

	ob.PlaceLimitOrder(10_000, sellOrder)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketPlaceOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrder := NewOrder(false, 20)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := NewOrder(true, 10)
	matches := ob.PlaceMartketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1) // the ask order is not completly filled since only a buy order of 10 was filled
	assert(t, len(ob.bids), 0) // the bid order should be removed when filled completly
	assert(t, ob.AskTotalVolume(), 10.0)

	// checking out match algorithm works properly
	// hint: always check an array length before accessing it directly like this matches[0]
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.Size, 0.0) // the buy order should be completly filled
}

func TestPlaceOrderMultiFill(t *testing.T) {
	ob := NewOrderBook()

	buyOrderA := NewOrder(true, 25)
	buyOrderB := NewOrder(true, 50)
	buyOrderC := NewOrder(true, 1)
	buyOrderD := NewOrder(true, 10)
	// Volume -> 25 + 50 + 1 + 10 = 86

	/***
	 Explanation:
	 	There are 4 buy orders
		There are 3 distict price levels (also called price limits)
		Price levels are 10_000, 9_000, 8_000
		When a buy order is at level 10_000 it means someone wants to buy crypto at $10_000

	***/
	ob.PlaceLimitOrder(10_000, buyOrderA)
	ob.PlaceLimitOrder(10_000, buyOrderB)
	ob.PlaceLimitOrder(20_000, buyOrderC)
	ob.PlaceLimitOrder(20_000, buyOrderD)

	assert(t, ob.BidTotalVolume(), 86.0)

	sellOrder := NewOrder(false, 20)
	matches := ob.PlaceMartketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 66.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.bids), 2)
	assert(t, len(ob.asks), 0)

}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrder := NewOrder(true, 10)
	ob.PlaceLimitOrder(10_000, buyOrder)

	assert(t, len(ob.bids), 1)
	assert(t, ob.BidTotalVolume(), 10.0)

	ob.CancelOrder(buyOrder)

	assert(t, ob.BidTotalVolume(), 0.0)
}
