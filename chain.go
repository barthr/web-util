package web

import "net/http"

// MiddlewareFunc is the signature for a middleware function
type MiddlewareFunc func(http.Handler) http.Handler

// MiddlewareChain is the chain of handlers which are getting called LIFO
type MiddlewareChain struct {
	chain []MiddlewareFunc
}

// New returns a new MiddlewareChain with the given middleware functions
func New(middlewareFuncs ...MiddlewareFunc) MiddlewareChain {
	return MiddlewareChain{
		chain: middlewareFuncs,
	}
}

// Wrap takes a final handler and returns the handler through all
// the middleware functions
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

// WrapHandlerFunc does the same as Wrap but takes an HandlerFunc instead of a http.Handler
func (mwc MiddlewareChain) WrapHandlerFunc(final HandlerFunc) http.Handler {
	return mwc.Wrap(final)
}

// Add create's a copy of the current chain and add's the given middlewareFuncs to the chain
func (mwc MiddlewareChain) Add(middlewareFuncs ...MiddlewareFunc) MiddlewareChain {
	newChain := mwc
	copy(newChain.chain, mwc.chain)
	newChain.chain = append(newChain.chain, middlewareFuncs...)
	return newChain
}
