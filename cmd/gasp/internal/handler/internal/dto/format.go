package dto

func FormatBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
