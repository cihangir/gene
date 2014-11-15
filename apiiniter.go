package gene

import "github.com/rcrowley/go-tigertonic"

// Initer is for initializing the api endpoints for a module
type Initer interface {
	Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux
}

// Customizer is for adding custom handlers to the module
type Customizer interface {
	Customize(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux
}
