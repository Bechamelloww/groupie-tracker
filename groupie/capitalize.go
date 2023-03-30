package groupie

func Capitalize(s string) string {
	modstr := []rune(s)
	for j := 0; j < len(s); j++ {
		if modstr[j] >= 97 && s[j] <= 122 {
			modstr[j] = modstr[j] - 32
		}
		if j > 0 {
			if modstr[j] >= 65 && modstr[j] <= 90 {
				modstr[j] = modstr[j] + 32
				if modstr[j] >= 97 && s[j] <= 122 && modstr[j-1] < 48 {
					modstr[j] = modstr[j] - 32
				}
				if modstr[j] >= 97 && s[j] <= 122 && modstr[j-1] > 57 && modstr[j-1] < 65 {
					modstr[j] = modstr[j] - 32
				}
				if modstr[j] >= 97 && s[j] <= 122 && modstr[j-1] > 90 && modstr[j-1] < 97 {
					modstr[j] = modstr[j] - 32
				}
				if modstr[j] >= 97 && s[j] <= 122 && modstr[j-1] > 122 {
					modstr[j] = modstr[j] - 32
				}
			}
		}
	}
	s = string(modstr)
	return s
}
