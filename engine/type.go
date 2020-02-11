package engine

type ParseFunc func(contents []byte, url string) ParseResult

type Request struct {
	Url           string
	ParseFunction ParseFunc
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type Item struct {
	Url     string
	Payload interface{}
}
