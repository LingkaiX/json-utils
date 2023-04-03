package jsonutils

import "regexp"

// a FSM with transition table for validating json
const (
	Start         = iota + 1
	InObjectStart // after '{ , accepts key or '}'
	InArrayStart  // after '[', accepts value (string, number, bool, array or object) or ']'
	InObjectInKey
	InObjectAfterKey    // accepts ':'
	InObjectBeforeValue // after ':', accepts value
	InString
	AfterValue         // after a valid value in object or array, accepts ',' or ']' / '}'
	InObjectBeforeKey  // after ',', only accepts key
	InArrayBeforeValue // after ',', only accepts value
	End                // only accepts empty strings

	InObject
	InArray

	// actions need to be implemented
	Continue    // go to next loop without any change
	SkipNext    // index +1
	Reject      // error, break intermediately
	CheckTrue   // then goto AfterValue
	CheckFalse  // then goto AfterValue
	CheckNumber // check it is a valid number then goto AfterValue
	BeforeValueOrKey
	EnterObject
	EnterArray
	ExitObject
	ExitArray

	CheckSpecialChar // '"', '\', '/', b, f, n, r, t, u following '\', \u should be like \uffff
)

// space, newline or tab
var spaceRegexp = regexp.MustCompile(`[\s\n\r\t]`)

// match of the rune is a digit or '-'
var numRegexp = regexp.MustCompile(`[\d-]`)

type transition struct {
	nextStates    map[rune]int
	actions       map[rune]int
	regexpActions map[*regexp.Regexp]int
	// regexpTransitions map[*regexp.Regexp]int
	defaultAction int
}

var stateTable = map[int]transition{
	Start: {
		actions: map[rune]int{
			'{': EnterObject,
			'[': EnterArray,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
	InObjectStart: {
		nextStates: map[rune]int{
			'"': InObjectInKey,
		},
		actions: map[rune]int{
			'}': ExitObject,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
	InArrayStart: {
		nextStates: map[rune]int{
			'"': InString,
		},
		actions: map[rune]int{
			't': CheckTrue,
			'f': CheckFalse,
			']': ExitArray,
			'[': EnterArray,
			'{': EnterObject,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue, numRegexp: CheckNumber},
		defaultAction: Reject,
	},
	InObjectInKey: {
		nextStates: map[rune]int{
			'"': InObjectAfterKey,
		},
		actions: map[rune]int{
			'/': CheckSpecialChar,
		},
		defaultAction: Continue,
	},
	InObjectAfterKey: {
		nextStates: map[rune]int{
			':': InObjectBeforeValue,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
	InObjectBeforeValue: {
		nextStates: map[rune]int{
			'"': InString,
		},
		actions: map[rune]int{
			't': CheckTrue,
			'f': CheckFalse,
			'[': EnterArray,
			'{': EnterObject,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue, numRegexp: CheckNumber},
		defaultAction: Reject,
	},
	InString: {
		nextStates: map[rune]int{
			'"': AfterValue,
		},
		actions: map[rune]int{
			'/': CheckSpecialChar,
		},
		defaultAction: Continue,
	},
	AfterValue: {
		actions: map[rune]int{
			',': BeforeValueOrKey,
			'}': ExitObject,
			']': ExitArray,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
	InObjectBeforeKey: {
		nextStates: map[rune]int{
			'"': InObjectInKey,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
	InArrayBeforeValue: {
		nextStates: map[rune]int{
			'"': InString,
		},
		actions: map[rune]int{
			't': CheckTrue,
			'f': CheckFalse,
			'[': EnterArray,
			'{': EnterObject,
		},
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue, numRegexp: CheckNumber},
		defaultAction: Reject,
	},
	End: {
		regexpActions: map[*regexp.Regexp]int{spaceRegexp: Continue},
		defaultAction: Reject,
	},
}
