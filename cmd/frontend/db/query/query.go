// Package query provides an expression tree structure which can be converted
// into where queries. It is used by DB APIs to expose a more powerful query
// interface.
package query

// Q is a query item. It is converted into a *sqlf.Query by Eval.
type Q interface{}
