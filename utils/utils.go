package utils

// MergeMaps - merges two map of string(key and value)
// for any common fields, if keepA is true fields in map `a`
// is given priority otherwise `b`
func MergeMaps(a, b map[string]string, keepA bool) map[string]string {
	c := make(map[string]string)

	for k, v := range b {
		c[k] = v
	}

	for k, v := range a {
		c[k] = v

		if vb, prs := b[k]; prs {
			// common key
			if !keepA {
				c[k] = vb
			}
		}
	}

	return c
}
