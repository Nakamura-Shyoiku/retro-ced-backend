package bounds

func CheckString(a []string, i int) string {
	if len(a) < (i + 1) {
		return ""
	}
	if len(a) != 0 && (i+1) == 0 {
		return ""
	}
	return a[i]
}
