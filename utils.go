package jsonutils

import (
	"fmt"
)

type report struct {
	str      string
	startAt  int
	endAt    int
	isObject bool
	depth    int
	isValid  bool
	errInfo  string
}

// validate a json string
// ref: https://restfulapi.net/json-data-types/
// ! UNCOMPLETED
func Validate(json string) bool {
	// pattern := regexp.MustCompile(`[a-z]`)
	var r report
	r.isValid = true
	var s stack
	s.push(Start)
	for k, v := range json {
		if s.peek() == Start {
			if v != ' ' && v != '\n' {
				r.isValid = false
				r.errInfo = "Hey, feed me JSON!"
				break
			}
		}
		if s.peek() == End {
			if v != ' ' && v != '\n' {
				r.isValid = false
				r.errInfo = "Is there something outside of the root?"
			}
			break
		}

		switch v {
		case '{':
			if s.peek() == Start {
				r.startAt = k
				r.isObject = true
			}
			s.push((InObject))
		case '}':
		case '[':
		case ']':
		case '"':
		case ':':
		case '\\':
		default:

		}
	}
	return false
}

// return the max depth of a json string
// with O(n) space complexity and time complexity
// todo: read from io stream
func Depth(json string) (int, error) {
	var depth, current int
	var s stack
	for _, v := range json {
		switch v {
		case '{':
			s.push(InObject)
			current++
			if current > depth {
				depth = current
			}
		case '[':
			s.push(InArray)
			current++
			if current > depth {
				depth = current
			}
		case '}':
			if s.peek() != InObject {
				return 0, fmt.Errorf("invalid json")
			} else {
				s.pop()
				current--
			}
		case ']':
			if s.peek() != InArray {
				return 0, fmt.Errorf("invalid json")
			} else {
				s.pop()
				current--
			}
		}
	}
	if !s.isEmpty() {
		return 0, fmt.Errorf("invalid json")
	}
	return depth, nil
}
