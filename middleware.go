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

func (mwc MiddlewareChain) WrapHandlerFunc(final HandlerFunc) http.Handler {
	return mwc.Wrap(final)
}

func (mwc MiddlewareChain) Add(middlewareFuncs ...MiddlewareFunc) MiddlewareChain {
	newChain := mwc
	copy(newChain.chain, mwc.chain)
	newChain.chain = append(newChain.chain, middlewareFuncs...)
	return newChain
}
