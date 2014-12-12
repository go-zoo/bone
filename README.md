bone
=======

## What is bone ?

Bone is a lightweight and lightning fast HTTP Multiplexer for Golang. It support URL variables, http method declaration
and custom NotFound handler.
Also bone is always gonna be supported and updated.

![alt tag](https://c2.staticflickr.com/2/1070/540747396_5542b42cca_z.jpg)

## Speed

- bone : 	 				           555  ns/op
- daryl/zeus :				       590  ns/op
- julienschmidt/httprouter : 611  ns/op
- net/http : 				         924  ns/op
- gorilla/mux : 			       1158 ns/op
- gorilla/pat : 			       1313 ns/op

[ These test are just for fun, all these router are great and really efficient. ]

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
