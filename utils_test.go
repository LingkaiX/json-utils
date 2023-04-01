package jsonutils

import (
	"testing"
)

var json = []string{
	"",
	`{}`,
	`[123, "123, 12.3, {}]`,
	`{
		"name": "John Smith"
	}`,
	`{
		"name": ["Jane Doe"],
		"age": 27,
		"isMarried": false
	}`,
	`[  {    "id": 1,    "name": "Product 1",    "price": 19.99  },  
	{    "id": 2,    "name": "Product 2",    "price": 29.99  },  
	{    "id": 3,    "name": "Product 3",    "price": 39.99  }]`,
	`{
		"name": "John Smith",
		"age": 32,
		"address": {
			"street": "123 Main St",
			"city": "Anytown",
			"state": "CA",
			"zip": "12345"
		}
	}`,
}

// only test valid json
func TestDepth(t *testing.T) {
	var expectation = []int{0, 1, 2, 1, 2, 2, 2}

	for k, v := range json {
		if i, _ := Depth(v); i != expectation[k] {
			t.Errorf("faulty depth: %v for %s", i, v)
		}
	}
}
