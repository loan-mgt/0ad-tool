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
	return SPrintMapRaw(m, "")
}

func SPrintMapRaw(m map[string]any, indent string) string {
	var result string = "\n" + indent
	for key, value := range m {
		switch v := value.(type) {
		case map[string]any:
			result += "\n" + indent + key + ":" + SPrintMapRaw(v, indent+"  ") + "\n"
		default:
			result += key + ": " + v.(string) + ", "
		}
	}
	return result
}

func PrintMap(m map[string]any) {
	println(SPrintMap(m))
}
