package engine

import "simple-golang-crawler/fetcher"

type ParseFunc func(contents []byte, url string) ParseResult

type Request struct {
	Url           string
	ParseFunction ParseFunc
	FetchFun      fetcher.FetchFun
}

func NewRequest(url string, parseFunction ParseFunc, fetchFun fetcher.FetchFun) *Request {
	return &Request{Url: url, ParseFunction: parseFunction, FetchFun: fetchFun}
}

type ParseResult struct {
	Requests []*Request
	Items    []*Item
}

type Item struct {
	Payload interface{}
}

func NewItem(payload interface{}) *Item {
	return &Item{Payload: payload}
}
