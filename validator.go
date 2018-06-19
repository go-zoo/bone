package bone

// Validator can be passed to a route to validate the params
type Validator func(string) bool
