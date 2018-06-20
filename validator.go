package bone

// Validator can be passed to a route to validate the params
type Validator func(string) bool

type validator struct {
	start    int
	end      int
	name     string
	validate Validator
}

func containsValidators(path string) []validator {
	var index []int

	for i, c := range path {
		if c == '|' {
			index = append(index, i)
		}
	}

	if len(index) > 0 {
		var validators []validator
		for i, pos := range index {
			if i+1 == len(index) {
				validators = append(validators, validator{
					start: pos,
					end:   len(path),
					name:  path[pos:len(path)],
				})
			} else {
				validators = append(validators, validator{
					start: pos,
					end:   index[i+1],
					name:  path[pos:index[i+1]],
				})
			}
		}
		return validators
	}
	return nil
}
