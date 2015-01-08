/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// Route content the required information for a valid route
// Path: is the Route URL
// Size: is the length of the path
// Token: is the value of each part of the path, split by /
// Pattern: is content information about the route, if it's have a route variable
// handler: is the handler who handle this route
// Method: define HTTP method on the route
type Route struct {
	Path    string
	Size    int
	Token   Token
	Pattern Pattern
	Handler http.Handler
	Method  string
}

// Token content all value of a spliting route path
// Tokens: string value of each token
// size: number of token
type Token struct {
	Tokens []string
	Size   int
}

type byLength []*Route

func (b byLength) Len() int {
	return len(b)
}

func (b byLength) Swap(i int, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byLength) Less(i int, j int) bool {
	return b[i].Token.Size < b[j].Token.Size
}

func (b byLength) Sort() {
	sort.Sort(b)
}

// Pattern content the required information for the route Pattern
// Exist: check if a variable was declare on the route
// ID: the name of the variable
// Pos: postition of var in the route path
// Value: is the value of the request parameters
type Pattern struct {
	Exist bool
	ID    string
	Pos   int
	Value map[string]string
}

// NewRoute return a pointer to a Route instance and call save() on it
func NewRoute(url string, h http.Handler) *Route {
	r := &Route{Path: url, Handler: h}
	r.save()
	return r
}

// Save, set automaticly the the Route.Size and Route.Pattern value
func (r *Route) save() {
	r.Size = len(r.Path)
	r.Token.Tokens = strings.Split(r.Path, "/")

	for i, s := range r.Token.Tokens {
		if len(s) >= 1 {
			if s[:1] == ":" {
				r.Pattern.Exist = true
				r.Pattern.ID = s[1:]
				r.Pattern.Pos = i
				r.Pattern.Value = make(map[string]string)
			}
		}
		r.Token.Size++
	}
}

// Match check if the request match the route Pattern
func (r *Route) Match(path string) (url.Values, bool) {
	ss := strings.Split(path, "/")
	if len(ss) == r.Token.Size && ss[r.Token.Size-1] != "" {
		if r.Path[:r.Pattern.Pos] == path[:r.Pattern.Pos] {
			uV := url.Values{}
			uV.Add(r.Pattern.ID, ss[r.Pattern.Pos])
			return uV, true
		}
	}
	return nil, false
}

// Get set the route method to Get
func (r *Route) Get() *Route {
	r.Method = "GET"
	return r
}

// Post set the route method to Post
func (r *Route) Post() *Route {
	r.Method = "POST"
	return r
}

// Put set the route method to Put
func (r *Route) Put() *Route {
	r.Method = "PUT"
	return r
}

// Delete set the route method to Delete
func (r *Route) Delete() *Route {
	r.Method = "DELETE"
	return r
}

// Head set the route method to Head
func (r *Route) Head() *Route {
	r.Method = "HEAD"
	return r
}

// Patch set the route method to Patch
func (r *Route) Patch() *Route {
	r.Method = "PATCH"
	return r
}

// Options set the route method to Options
func (r *Route) Options() *Route {
	r.Method = "OPTIONS"
	return r
}
