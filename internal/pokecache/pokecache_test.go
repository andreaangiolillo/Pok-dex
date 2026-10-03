package pokecache

import (
	"testing"
	"time"
)



func TestCache(t *testing.T) {
	m, _ := time.ParseDuration("1h30m")
	testCase := []struct {
		input map[string][]byte
		duration time.Duration
		expected Cache
	}{
		{
			input: map[string][]byte{
				"test": []byte("test"),
				"test2": []byte("test2"),
			},
			expected: Cache{
				Entries: map[string]*CacheEntry{
					"test": &CacheEntry{
						val: []byte("test"),
					},
					"test2": &CacheEntry{
						val: []byte("test2"),
					},
				},
			},
			duration: m,
		},
	}

	for _, c := range testCase {
		cache := NewCache(c.duration)
		for k, v := range c.input {
			cache.Add(k, v)
		}

		for k, expected := range c.input {
			v, ok := cache.Get(k)

			if !ok || string(v) != string(expected) {
				t.Errorf("got %s, expected %s", v, expected)
			}
		}

	}

}


