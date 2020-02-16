package main

type PriceCache interface {
	Contains(ticker string) bool
}

type NullCache struct {
}

func (cache NullCache) Contains(ticker string) bool {
	return false
}
