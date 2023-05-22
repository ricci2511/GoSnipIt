package main

type contextKey string

const (
	isAuthenticatedContextKey = contextKey("isAuthenticated")
	snippetContextKey         = contextKey("snippet")
)
