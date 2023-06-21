package main

func splitBySpace(s string) []string {
	return splitByChar(s, ' ')
}

func splitByChar(s string, char byte) []string {
	parts := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == char {
			if i > start {
				parts = append(parts, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}