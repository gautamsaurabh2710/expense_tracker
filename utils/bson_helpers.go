package utils

import "go.mongodb.org/mongo-driver/bson"

func BsonFloat(m bson.M, key string) float64 {
	if m == nil {
		return 0
	}

	switch v := m[key].(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}
