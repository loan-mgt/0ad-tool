package utils

func SPrintMapAny(m any) string {
	switch v := m.(type) {
	case map[string]any:
		return SPrintMap(v)
	case string:
		return v
	default:
		return ""
	}
}

func SPrintMap(m map[string]any) string {
	var result string = "\n[ "
	for key, value := range m {
		switch v := value.(type) {
		case map[string]any:
			result += key + ":\n" + SPrintMap(v)
		default:
			result += key + ": " + v.(string) + ", "
		}
	}
	result += "]"
	return result
}

func PrintMap(m map[string]any) {
	println(SPrintMap(m))
}
