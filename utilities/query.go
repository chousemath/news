package utilities

// AddToQuery adds to a query string if the string value is present
func AddToQuery(key string, value string) string {
	if len(value) > 0 {
		return "&" + key + "=" + value
	}
	return ""
}
