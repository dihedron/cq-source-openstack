package format

import "encoding/json"

// ToJSON returns the given object's representation as a JSON string.
func ToJSON(v any) string {
	data, _ := json.Marshal(v)
	return string(data)
}

// ToPrettyJSON returns the given object's representation as a pretty-printed JSON string.
func ToPrettyJSON(v any) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}
