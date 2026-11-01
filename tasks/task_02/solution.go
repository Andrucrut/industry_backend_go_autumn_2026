package main

func rotateRunes(s string, shift int) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	if shift > 0 {
		shift = shift % len(runes)

		part := runes[shift:]
		first := runes[:shift]

		return string(part) + string(first)
	}
	if shift < 0 {
		shift = shift % len(runes)

		part := runes[len(runes)+shift:]
		first := runes[:len(runes)+shift]

		return string(part) + string(first)
	}
	return string(runes)
}
