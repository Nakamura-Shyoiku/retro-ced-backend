package ternary

func String(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

func Int(a, b int) int {
	if a == 0 {
		return b
	}
	return a
}

func Int64(a, b int64) int64 {
	if a == 0 {
		return b
	}
	return a
}
