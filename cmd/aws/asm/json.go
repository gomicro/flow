package asm

import (
	gofmt "fmt"
)

func getJsonValue(heap map[string]interface{}, keys []string) (string, error) {
	if len(keys) > 1 {
		var key string
		key, keys = keys[0], keys[1:]

		next, found := heap[key]
		if !found {
			return "", gofmt.Errorf("next key not found: %v", key)
		}

		newHeap, ok := next.(map[string]interface{})
		if !ok {
			return "", gofmt.Errorf("next level of json is not an object: [%v] = %v", key, next)
		}

		return getJsonValue(newHeap, keys)
	}

	val, found := heap[keys[0]]
	if !found {
		return "", gofmt.Errorf("key not found: %v", keys[0])
	}

	return gofmt.Sprintf("%v", val), nil
}
