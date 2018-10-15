package muxie

import (
	"net/http"
)

// GetParam returns the path parameter value based on its key, i.e
// "/hello/:name", the parameter key is the "name".
// For example if a route with pattern of "/hello/:name" is inserted to the `Trie` or handlded by the `Mux`
// and the path "/hello/kataras" is requested through the `Mux#ServeHTTP -> Trie#Search`
// then the `GetParam("name")` will return the value of "kataras".
// If not associated value with that key is found then it will return an empty string.
//
// The function will do its job only if the given "w" http.ResponseWriter interface is an `paramsWriter`.
func GetParam(w http.ResponseWriter, key string) string {
	if store, ok := w.(*paramsWriter); ok {
		return store.Get(key)
	}

	return ""
}

// GetParams returns all the available parameters based on the "w" http.ResponseWriter which should be a *paramsWriter.
//
// The function will do its job only if the given "w" http.ResponseWriter interface is an `paramsWriter`.
func GetParams(w http.ResponseWriter) []ParamEntry {
	if store, ok := w.(*paramsWriter); ok {
		return store.params
	}

	return nil
}

// SetParam sets manually a parameter to the "w" http.ResponseWriter which should be a *paramsWriter.
// This is not commonly used by the end-developers.
func SetParam(w http.ResponseWriter, key, value string) bool {
	if store, ok := w.(*paramsWriter); ok {
		store.Set(key, value)
		return true
	}

	return false
}

type paramsWriter struct {
	http.ResponseWriter
	params []ParamEntry
}

// ParamEntry holds the Key and the Value of a named path parameter.
type ParamEntry struct {
	Key   string
	Value string
}

// Set implements the `ParamsSetter` which `Trie#Search` needs to store the parameters, if any.
// These are decoupled because end-developers may want to use the trie to design a new Mux of their own
// or to store different kind of data inside it.
func (pw *paramsWriter) Set(key, value string) {
	if ln := len(pw.params); cap(pw.params) > ln {
		pw.params = pw.params[:ln+1]
		p := &pw.params[ln]
		p.Key = key
		p.Value = value
		return
	}

	pw.params = append(pw.params, ParamEntry{
		Key:   key,
		Value: value,
	})
}

// Get returns the value of the associated parameter based on its key/name.
func (pw *paramsWriter) Get(key string) string {
	n := len(pw.params)
	for i := 0; i < n; i++ {
		if kv := pw.params[i]; kv.Key == key {
			return kv.Value
		}
	}

	return ""
}

func (pw *paramsWriter) reset(w http.ResponseWriter) {
	pw.ResponseWriter = w
	pw.params = pw.params[0:0]
}