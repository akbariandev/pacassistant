package helper

import "testing"

func TestCache(t *testing.T) {
	cache := NewCache[string, int]()

	// Test Get method
	cache.Add("one", 1)
	val, ok := cache.Get("one")
	if !ok || val != 1 {
		t.Errorf("Get method failed. Expected value: 1, Got: %v, Ok: %v", val, ok)
	}

	// Test Add method

	if ok := cache.Add("two", 2); !ok {
		t.Errorf("Add method failed")
	}
	val, ok = cache.Get("two")
	if !ok || val != 2 {
		t.Errorf("Add method failed. Expected value: 2, Got: %v, Ok: %v", val, ok)
	}

	// Test Keys method
	keys := cache.Keys()
	expectedKeys := []string{"one", "two"}
	if !equalStringSlices(keys, expectedKeys) {
		t.Errorf("Keys method failed. Expected keys: %v, Got: %v", expectedKeys, keys)
	}

	// Test Delete method
	ok = cache.Delete("one")
	if !ok {
		t.Errorf("Delete method failed. Expected: true, Got: %v", ok)
	}
	val, ok = cache.Get("one") //nolint
	if ok {
		t.Errorf("Delete method failed. Expected Ok: false, Got: true")
	}

	// Test NewCache method
	newCache := NewCache[string, int]()
	if newCache == nil {
		t.Error("NewCache method failed. Expected non-nil cache.")
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
