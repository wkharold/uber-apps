// Package uber provides a set of types and interfaces for constructing UBER hypermedia documents.
package uber

// udata represents the individual data elements of an Uber hypermedia document.
type Data struct {
	ID         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Rel        []string `json:"rel,omitempty"`
	Label      string   `json:"label,omitempty"`
	URL        string   `json:"url,omitempty"`
	Templated  bool     `json:"templated,omitempty"`
	Action     string   `json:"action,omitempty"`
	Transclude bool     `json:"transclude,omitempty"`
	Model      string   `json:"model,omitempty"`
	Sending    string   `json:"sending,omitempty"`
	Accepting  []string `json:"accepting,omitempty"`
	Value      string   `json:"value,omitempty"`
	Data       []Data   `json:"data,omitempty"`
}

// ubody is the body of an Uber hypermedia document.
type Body struct {
	Version string `json:"version"`
	Data    []Data `json:"data,omitempty"`
	Error   []Data `json:"error,omitempty"`
}

// udoc represents an Uber hypermedia document.
type Doc struct {
	Uber Body `json:"uber"`
}
