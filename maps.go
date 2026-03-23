package goutils

// MapKeys returns a slice containing all keys from the provided map.
// The key order is not guaranteed.
func MapKeys(m map[string]any) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// MapMerge merges maps into initMap.
// When the same key appears more than once, conflictFunc is used to resolve
// the value, where v0 is the current value in initMap and v1 is the incoming value.
// It mutates and returns initMap.
func MapMerge[K comparable, V any](conflictFunc func(key K, v0, v1 V) V, initMap map[K]V, maps ...map[K]V) map[K]V {
	for _, _map := range maps {
		if _map == nil {
			continue
		}
		for mk, mv := range _map {
			if initMapVal, ok := initMap[mk]; !ok {
				initMap[mk] = mv
			} else {
				initMap[mk] = conflictFunc(mk, initMapVal, mv)
			}
		}
	}
	return initMap
}
