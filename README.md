bone
=======

## What is bone ?

bone is a lightweight HTTP Multiplexer. It support URL variables and http method declaration.

## Example

``` go

package main

import(
  "net/http"

  "github.com/squiidz/bone"
)

func main () {
  mux := bone.NewMux()
  
  mux.Handle("/home/:id", HomeHandler")
  
  http.ListenAndServe(":8080", mux)
}

```
## TODO

- Url parameters
- Custom Regex parameters

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Write Tests!
4. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

## License
MIT
