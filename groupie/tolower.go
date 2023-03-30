package groupie

func ToHigher(s string) string { // Permet de mettre en MAJUSCULE une string
	modstr := []rune(s)
	for j := 0; j < len(s); j++ {
		if modstr[j] >= 97 && s[j] <= 122 {
			modstr[j] = modstr[j] - 32
		} else {
			continue
		}
	}
	s = string(modstr)
	return s
}

func ToLower(modstr []rune) []rune {
	var runetab []rune
	for j := 0; j < len(modstr); j++ {
		if modstr[j] <= 64 {
			j++
		} else {
			if modstr[j] >= 65 && modstr[j] <= 90 {
				runetab = append(runetab, modstr[j]+32)
			} else if modstr[j] >= 97 && modstr[j] <= 122 {
				runetab = append(runetab, modstr[j])
			} else {
				j++
			}
		}
	}
	return runetab
}
