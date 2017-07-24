package web

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

type MiddlewareChain struct {
	chain []MiddlewareFunc
}

func New(middlewareFuncs ...MiddlewareFunc) MiddlewareChain {
	return MiddlewareChain{
		chain: middlewareFuncs,
	}
}

func (mwc MiddlewareChain) Wrap(final http.Handler) http.Handler {
	if final == nil {
		panic("cannot wrap nil handler")
	}
	topIndex := len(mwc.chain) - 1
	for i := range mwc.chain {
		final = mwc.chain[topIndex-i](final)
	}
	return final
}
