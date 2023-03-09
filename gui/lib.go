package gui

import (
	"log"
	"sort"
	"sync"
)

func SortedKeys(m sync.Map) []string {
	keys := []string{}
	m.Range(func(key, value any) bool {
		strKey, ok := key.(string)
		if ok {
			keys = append(keys, strKey)
		} else {
			log.Printf("SortedKeys got non-string key: %v", key)
		}
		return true
	})

	sort.Strings(keys)
	return keys
}
