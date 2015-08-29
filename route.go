/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
	"regexp"
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
	Path     string
	Method   string
	Size     int
	Token    Token
	wildCard bool
	wildPos  int
	Params   bool
	Pattern  map[int]string
	Regex    bool
	Compile  map[int]*regexp.Regexp
	Tag      map[int]string
	Handler  http.Handler
}

// Token content all value of a spliting route path
// Tokens: string value of each token
// size: number of token
type Token struct {
	raw    []int
	Tokens []string
	Size   int
}

// NewRoute return a pointer to a Route instance and call save() on it
func NewRoute(url string, h http.Handler) *Route {
	r := &Route{Path: url, Handler: h}
	r.save()
	return r
}

// Save, set automaticly the the Route.Size and Route.Pattern values
func (r *Route) save() {
	r.Size = len(r.Path)
	r.Token.Tokens = strings.Split(r.Path, "/")
	if string(r.Path[0]) == "*" {
		r.wildCard = true
		r.wildPos = 0
	}
	for i, s := range r.Token.Tokens {
		if len(s) >= 1 {
			switch s[:1] {
			case ":":
				if !r.Params {
					r.Pattern = make(map[int]string)
				}
				r.Pattern[i] = s[1:]
				r.Params = true
			case "#":
				if !r.Regex {
					r.Compile = make(map[int]*regexp.Regexp)
					r.Tag = make(map[int]string)
				}
				tmp := strings.Split(s, "^")
				r.Tag[i] = tmp[0][1:]
				r.Compile[i] = regexp.MustCompile("^" + tmp[1][:len(tmp[1])-1])
				r.Regex = true
			case "*":
				r.wildCard = true
				r.wildPos = i
			default:
				r.Token.raw = append(r.Token.raw, i)
			}
		}
		r.Token.Size++
	}
}

// Match check if the request match the route Pattern
func (r *Route) Match(req *http.Request) bool {
	ss := strings.Split(req.URL.Path, "/")
	if len(ss) == r.Token.Size || r.wildCard {
		if r.wildCard && r.wildPos == 0 {
			rst := len(req.RequestURI) - len(r.Path[1:])
			if rst > 0 && r.Path[1:] == req.RequestURI[rst:] {
				return true
			}
			return false
		}
		for i, v := range r.Token.raw {
			if ss[v] != r.Token.Tokens[v] {
				if r.wildCard && i == r.wildPos {
					return true
				}
				return false
			}
		}
		vars.Lock()
		vars.m[req] = map[string]string{}
		vars.Unlock()
		if r.Regex {
			for k, v := range r.Compile {
				if v.MatchString(ss[k]) {
					vars.Lock()
					vars.m[req][r.Tag[k]] = ss[k]
					vars.Unlock()
				} else {
					return false
				}
			}
		}
		for k, v := range r.Pattern {
			vars.Lock()
			vars.m[req][v] = ss[k]
			vars.Unlock()
		}
		return true
	}
	return false
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

func (r *Route) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if r.Method != "" {
		if req.Method == r.Method {
			r.Handler.ServeHTTP(rw, req)
			return
		}
		http.NotFound(rw, req)
		return
	}
	r.Handler.ServeHTTP(rw, req)
}
