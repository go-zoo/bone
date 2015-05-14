bone [![GoDoc](https://godoc.org/github.com/squiidz/bone?status.png)](http://godoc.org/github.com/go-zoo/bone) [![Build Status](https://travis-ci.org/go-zoo/bone.svg)](https://travis-ci.org/go-zoo/bone) [![Codeship Status for squiidz/bone](https://codeship.com/projects/172dae70-802a-0132-9038-321707412590/status?branch=master)](https://codeship.com/projects/57454)
=======

## What is bone ?

Bone is a lightweight and lightning fast HTTP Multiplexer for Golang. It support URL variables with regex parameters, Http method declaration
and custom NotFound handler.

![alt tag](https://c2.staticflickr.com/2/1070/540747396_5542b42cca_z.jpg)

## Speed

```
- BenchmarkBoneMux        10000000               118 ns/op
- BenchmarkZeusMux          100000               144 ns/op
- BenchmarkHttpRouterMux  10000000               134 ns/op
- BenchmarkNetHttpMux      3000000               580 ns/op
- BenchmarkGorillaMux       300000              3333 ns/op
- BenchmarkGorillaPatMux   1000000              1889 ns/op
```

 These test are just for fun, all these router are great and really efficient. 
 Bone do not pretend to be the fastest router for every job. 

## Example

``` go

package main

import(
  "net/http"

  "github.com/go-zoo/bone"
)

func main () {
  mux := bone.New()
  
  // mux.Get, Post, etc ... takes http.Handler
  mux.Get("/home/:id", HomeHandler)
  mux.Get("/profil/:id/:var", ProfilHandler)
  mux.Post("/data", DataHandler)

  // Support Regex Route params
  mux.Get("/index/#id^[0-9]$", IndexHandler)

  // Handle take http.Handler
  mux.Handle("/", http.HandlerFunc(RootHandler))

  http.ListenAndServe(":8080", mux)
}

func Handler(rw http.ResponseWriter, req *http.Request) {
  // Get the value of the "id" parameters.
  val := bone.GetValue(req, "id")

  rw.Write([]byte(val))
}

```
## Changelog

#### Update 26 April 2015

- Add Support for regex parameters, using ` # ` instead of ` : `.
- Add Mux method ` mux.GetFunc(), mux.PostFunc(), etc ... `, takes ` http.HandlerFunc ` instead of ` http.Handler `.

Example : 
``` go
func main() {
    mux.GetFunc("/route/#var^[a-z]$", handler) 
}

func handler(rw http.ResponseWriter, req *http.Request) {
    bone.GetValue(req, "var")
}
``` 

#### Update 29 january 2015

- Speed improvement for url parameters, from ```~ 1500 ns/op ``` to ```~ 1000 ns/op ```. 

#### Update 25 december 2014

After trying to find a way of using the default url.Query() for route parameters, i decide to change the way bone is dealing with this. url.Query() is too slow for good router performance.
So now to get the parameters value in your handler, you need to use 
` bone.GetValue(req, key) ` instead of ` req.Url.Query().Get(key) `.
This change give a big speed improvement for every kind of application using route parameters, like ~80x faster ...
Really sorry for breaking things, but i think it's worth it.  

## TODO

- DOC
- More Testing
- Debugging
- Optimisation

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Write Tests!
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

## License
MIT

## Links

Middleware Chaining module : [Claw](https://github.com/go-zoo/claw)


[![Bitdeli Badge](https://d2weczhvl823v0.cloudfront.net/go-zoo/bone/trend.png)](https://bitdeli.com/free "Bitdeli Badge")

