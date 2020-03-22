package maps

func CopyMap(m map[string]string) map[string]string {
	c := make(map[string]string)
	if m == nil {
		return c
	}

	for v, value := range m {
		c[v] = value
	}

	return c
}

func MergeMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		if m == nil {
			continue
		}

		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
