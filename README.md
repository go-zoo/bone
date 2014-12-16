bone [![GoDoc](https://godoc.org/github.com/squiidz/bone?status.png)](http://godoc.org/github.com/squiidz/bone) [![Build Status](https://travis-ci.org/squiidz/bone.svg)](https://travis-ci.org/squiidz/bone)
=======

## What is bone ?

Bone is a lightweight and lightning fast HTTP Multiplexer for Golang. It support URL variables, http method declaration
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

[ These test are just for fun, all these router are great and really efficient. ]
[ Bone do not pretend to be the fastest router for every job. ]
[ Big project should use a other solution for the moment. ]

## Example

``` go

package main

import(
  "net/http"

  "github.com/squiidz/bone"
)

func main () {
  mux := bone.New()
  
  // Method takes http.HandlerFunc
  mux.Get("/home/:id", HomeHandler)
  mux.Post("/data", DataHandler)

  // Handle take http.Handler
  mux.Handle("/", http.HandlerFunc(RootHandler))

  http.ListenAndServe(":8080", mux)
}

```
## TODO

- DOC
- More Testing
- Debugging
- Optimisation
- Refactoring

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

Middleware Chaining module : [Claw](https://github.com/squiidz/claw)
