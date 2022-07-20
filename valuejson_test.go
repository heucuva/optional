package optional_test

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/heucuva/optional"
)

type marshalTestJSON[T any] struct {
	test     string
	value    optional.Value[T]
	expected string
	run      func(*testing.T)
}

func (ti marshalTestJSON[T]) runSupported(t *testing.T) {
	t.Helper()
	blob, err := json.Marshal(ti.value)
	if err != nil {
		t.Fatal(err)
	}
	if observed := string(blob); strings.Compare(ti.expected, observed) != 0 {
		t.Fatalf("expected %q, got %q", ti.expected, observed)
	}
}

func (ti marshalTestJSON[T]) runUnsupportedValue(t *testing.T) {
	t.Helper()
	_, err := json.Marshal(ti.value)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedValue *json.UnsupportedValueError
	if !errors.As(err, &unsupportedValue) {
		t.Fatal(err)
	}
}

func (ti marshalTestJSON[T]) runUnsupportedType(t *testing.T) {
	t.Helper()
	_, err := json.Marshal(ti.value)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedType *json.UnsupportedTypeError
	if !errors.As(err, &unsupportedType) {
		t.Fatal(err)
	}
}

func marshalSupportedJSON[T any](name string, value T, expected string) marshalTestJSON[T] {
	ti := marshalTestJSON[T]{
		test:     name,
		value:    optional.NewValue(value),
		expected: expected,
	}
	ti.run = ti.runSupported
	return ti
}

func marshalUnsupportedJSONValue[T any](name string, value T) marshalTestJSON[T] {
	ti := marshalTestJSON[T]{
		test:  name,
		value: optional.NewValue(value),
	}
	ti.run = ti.runUnsupportedValue
	return ti
}

func marshalUnsupportedJSONType[T any](name string, value T) marshalTestJSON[T] {
	ti := marshalTestJSON[T]{
		test:  name,
		value: optional.NewValue(value),
	}
	ti.run = ti.runUnsupportedType
	return ti
}

func testMarshalJSON[T any](t *testing.T, tests ...marshalTestJSON[T]) {
	t.Helper()

	t.Run("Unset", marshalTestJSON[T]{expected: "null"}.runSupported)

	for _, ti := range tests {
		t.Run(ti.test, ti.run)
	}
}

func TestMarshalJSON(t *testing.T) {
	// Boolean
	t.Run("Bool", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON("True", true, `true`),
			marshalSupportedJSON("False", false, `false`),
		)
	})

	// Signed Integer
	t.Run("Int", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON("Zero", 0, `0`),
			marshalSupportedJSON("Positive", math.MaxInt, `9223372036854775807`),
			marshalSupportedJSON("Negative", math.MinInt, `-9223372036854775808`),
		)
	})
	t.Run("Int8", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[int8]("Zero", 0, `0`),
			marshalSupportedJSON[int8]("Positive", math.MaxInt8, `127`),
			marshalSupportedJSON[int8]("Negative", math.MinInt8, `-128`),
		)
	})
	t.Run("Int16", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[int16]("Zero", 0, `0`),
			marshalSupportedJSON[int16]("Positive", math.MaxInt16, `32767`),
			marshalSupportedJSON[int16]("Negative", math.MinInt16, `-32768`),
		)
	})
	t.Run("Int32", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[int32]("Zero", 0, `0`),
			marshalSupportedJSON[int32]("Positive", math.MaxInt32, `2147483647`),
			marshalSupportedJSON[int32]("Negative", math.MinInt32, `-2147483648`),
		)
	})

	// Unsigned integer
	t.Run("Uint", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[uint]("Zero", 0, `0`),
			marshalSupportedJSON[uint]("Max", math.MaxUint, `18446744073709551615`),
		)
	})
	t.Run("Uint8", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[uint8]("Zero", 0, `0`),
			marshalSupportedJSON[uint8]("Max", math.MaxUint8, `255`),
		)
	})
	t.Run("Uint16", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[uint16]("Zero", 0, `0`),
			marshalSupportedJSON[uint16]("Max", math.MaxUint16, `65535`),
		)
	})
	t.Run("Uint32", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[uint32]("Zero", 0, `0`),
			marshalSupportedJSON[uint32]("Max", math.MaxUint32, `4294967295`),
		)
	})

	// Floating point
	t.Run("Float32", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[float32]("ZeroPositive", 0.0, `0`),
			marshalSupportedJSON("ZeroNegative", math.Float32frombits(0x80000000), `-0`),
			marshalSupportedJSON[float32]("Positive", math.MaxFloat32, `3.4028235e+38`),
			marshalSupportedJSON[float32]("Negative", -math.MaxFloat32, `-3.4028235e+38`),
			marshalSupportedJSON[float32]("Smallest", math.SmallestNonzeroFloat32, `1e-45`),
			marshalUnsupportedJSONValue("QNaN", math.Float32frombits(0x7FFFFFFF)),
			marshalUnsupportedJSONValue("SNaN", math.Float32frombits(0x7FbFFFFF)),
			marshalUnsupportedJSONValue("PositiveInf", math.Float32frombits(0x7F800000)),
			marshalUnsupportedJSONValue("NegativeInf", math.Float32frombits(0xFF800000)),
		)
	})
	t.Run("Float64", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON("ZeroPositive", 0.0, `0`),
			marshalSupportedJSON("ZeroNegative", math.Float64frombits(0x8000000000000000), `-0`),
			marshalSupportedJSON("Positive", math.MaxFloat64, `1.7976931348623157e+308`),
			marshalSupportedJSON("Negative", -math.MaxFloat64, `-1.7976931348623157e+308`),
			marshalSupportedJSON("Smallest", math.SmallestNonzeroFloat64, `5e-324`),
			marshalUnsupportedJSONValue("QNaN", math.Float64frombits(0x7FFFFFFFFFFFFFFF)),
			marshalUnsupportedJSONValue("SNaN", math.Float64frombits(0x7FF7FFFFFFFFFFFF)),
			marshalUnsupportedJSONValue("PositiveInf", math.Float64frombits(0x7FF0000000000000)),
			marshalUnsupportedJSONValue("NegativeInf", math.Float64frombits(0xFFF0000000000000)),
		)
	})

	// Complex
	t.Run("Complex64", func(t *testing.T) {
		testMarshalJSON(t,
			marshalUnsupportedJSONType("BothZeroPositive", complex(float32(0.0), float32(0.0))),
		)
	})
	t.Run("Complex128", func(t *testing.T) {
		testMarshalJSON(t,
			marshalUnsupportedJSONType("BothZeroPositive", complex(float64(0.0), float64(0.0))),
		)
	})

	// Rune
	// NOTE: rune is effectively uint16
	t.Run("Rune", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON("Alpha", 'A', `65`),
			marshalSupportedJSON("Unicode", '\u2E9F', `11935`),
		)
	})

	// String
	t.Run("String", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON("Empty", "", `""`),
			marshalSupportedJSON("NonEmpty", "The quick brown fox", `"The quick brown fox"`),
		)
	})

	// Slice
	t.Run("Slice", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[[]string]("Null", nil, `null`),
			marshalSupportedJSON("Empty", []string{}, `[]`),
			marshalSupportedJSON("NonEmpty", []string{"The quick brown fox"}, `["The quick brown fox"]`),
		)
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		testMarshalJSON(t,
			marshalSupportedJSON[map[string]string]("Null", nil, `null`),
			marshalSupportedJSON("Empty", map[string]string{}, `{}`),
			marshalSupportedJSON("NonEmpty", map[string]string{"entry": "The quick brown fox"}, `{"entry":"The quick brown fox"}`),
		)
	})

	// Struct
	t.Run("Struct", func(t *testing.T) {
		t.Run("NotMarshalled", func(t *testing.T) {
			type testStructNotMarshalled struct {
				value int
			}
			var notMarshalled testStructNotMarshalled
			_ = notMarshalled.value
			testMarshalJSON(t,
				marshalSupportedJSON("Set", notMarshalled, `{}`),
			)
		})
		t.Run("Hidden", func(t *testing.T) {
			type testStructHidden struct {
				Hidden int `json:"-"`
			}
			var hidden testStructHidden
			testMarshalJSON(t,
				marshalSupportedJSON("Set", hidden, `{}`),
			)
		})
		t.Run("OneField", func(t *testing.T) {
			type testStructOneField struct {
				Value int `json:"value"`
			}
			var oneField testStructOneField
			testMarshalJSON(t,
				marshalSupportedJSON("Set", oneField, `{"value":0}`),
			)
		})
		t.Run("TwoFields", func(t *testing.T) {
			type testStructTwoFields struct {
				A int  `json:"a"`
				B bool `json:"b"`
			}
			var twoFields testStructTwoFields
			testMarshalJSON(t,
				marshalSupportedJSON("Set", twoFields, `{"a":0,"b":false}`),
			)
		})
		t.Run("EmbeddedOptional", func(t *testing.T) {
			type testStructEmbeddedOptional struct {
				Value optional.Value[int] `json:"value"`
			}
			var embeddedUnset testStructEmbeddedOptional
			embeddedSet := testStructEmbeddedOptional{
				Value: optional.NewValue(5),
			}
			testMarshalJSON(t,
				marshalSupportedJSON("SetValueUnset", embeddedUnset, `{"value":null}`),
				marshalSupportedJSON("SetValueSet", embeddedSet, `{"value":5}`),
			)
		})
	})
}

type unmarshalTestJSON[T any] struct {
	test     string
	data     string
	comparer func(observed optional.Value[T]) (optional.Value[T], bool)
	run      func(*testing.T)
}

func (ti unmarshalTestJSON[T]) runSupported(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := json.Unmarshal([]byte(ti.data), &observed)
	if err != nil {
		t.Fatal(err)
	}
	if expected, success := ti.comparer(observed); !success {
		t.Fatalf("expected %+v, got %+v", expected, observed)
	}
}

func (ti unmarshalTestJSON[T]) runUnsupportedValue(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := json.Unmarshal([]byte(ti.data), &observed)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedValue *json.SyntaxError
	if !errors.As(err, &unsupportedValue) {
		t.Fatal(err)
	}
}

func (ti unmarshalTestJSON[T]) runUnsupportedType(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := json.Unmarshal([]byte(ti.data), &observed)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedType *json.SyntaxError
	if !errors.As(err, &unsupportedType) {
		t.Fatal(err)
	}
}

func unmarshalSupported[T any](name string, data string, value T) unmarshalTestJSON[T] {
	ti := unmarshalTestJSON[T]{
		test: name,
		data: data,
		comparer: func(observed optional.Value[T]) (optional.Value[T], bool) {
			expected := optional.NewValue(value)
			if observed.IsSet() != expected.IsSet() {
				return expected, false
			}
			observedValue, _ := observed.Get()
			return expected, reflect.DeepEqual(observedValue, value)
		},
	}
	ti.run = ti.runSupported
	return ti
}

func unmarshalUnsupportedValue[T any](name string, data string) unmarshalTestJSON[T] {
	ti := unmarshalTestJSON[T]{
		test: name,
		data: data,
	}
	ti.run = ti.runUnsupportedValue
	return ti
}

func unmarshalUnsupportedType[T any](name string, data string) unmarshalTestJSON[T] {
	ti := unmarshalTestJSON[T]{
		test: name,
		data: data,
	}
	ti.run = ti.runUnsupportedType
	return ti
}

func testUnmarshalJSON[T any](t *testing.T, tests ...unmarshalTestJSON[T]) {
	t.Helper()

	t.Run("Unset", marshalTestJSON[T]{expected: "null"}.runSupported)

	for _, ti := range tests {
		t.Run(ti.test, ti.run)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	// Boolean
	t.Run("Bool", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported("True", `true`, true),
			unmarshalSupported("False", `false`, false),
		)
	})

	// Signed Integer
	t.Run("Int", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported("Zero", `0`, 0),
			unmarshalSupported("Positive", `9223372036854775807`, math.MaxInt),
			unmarshalSupported("Negative", `-9223372036854775808`, math.MinInt),
		)
	})
	t.Run("Int8", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[int8]("Zero", `0`, 0),
			unmarshalSupported[int8]("Positive", `127`, math.MaxInt8),
			unmarshalSupported[int8]("Negative", `-128`, math.MinInt8),
		)
	})
	t.Run("Int16", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[int16]("Zero", `0`, 0),
			unmarshalSupported[int16]("Positive", `32767`, math.MaxInt16),
			unmarshalSupported[int16]("Negative", `-32768`, math.MinInt16),
		)
	})
	t.Run("Int32", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[int32]("Zero", `0`, 0),
			unmarshalSupported[int32]("Positive", `2147483647`, math.MaxInt32),
			unmarshalSupported[int32]("Negative", `-2147483648`, math.MinInt32),
		)
	})

	// Unsigned integer
	t.Run("Uint", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[uint]("Zero", `0`, 0),
			unmarshalSupported[uint]("Max", `18446744073709551615`, math.MaxUint),
		)
	})
	t.Run("Uint8", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[uint8]("Zero", `0`, 0),
			unmarshalSupported[uint8]("Max", `255`, math.MaxUint8),
		)
	})
	t.Run("Uint16", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[uint16]("Zero", `0`, 0),
			unmarshalSupported[uint16]("Max", `65535`, math.MaxUint16),
		)
	})
	t.Run("Uint32", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[uint32]("Zero", `0`, 0),
			unmarshalSupported[uint32]("Max", `4294967295`, math.MaxUint32),
		)
	})

	// Floating point
	t.Run("Float32", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[float32]("ZeroPositive", `0`, 0.0),
			unmarshalSupported("ZeroNegative", `-0`, math.Float32frombits(0x80000000)),
			unmarshalSupported[float32]("Positive", `3.4028235e+38`, math.MaxFloat32),
			unmarshalSupported[float32]("Negative", `-3.4028235e+38`, -math.MaxFloat32),
			unmarshalSupported[float32]("Smallest", `1e-45`, math.SmallestNonzeroFloat32),
			unmarshalUnsupportedValue[float32]("QNaN", `qnan`),
			unmarshalUnsupportedValue[float32]("SNaN", `snan`),
			unmarshalUnsupportedValue[float32]("PositiveInf", `inf`),
			unmarshalUnsupportedValue[float32]("NegativeInf", `-inf`),
		)
	})
	t.Run("Float64", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported("ZeroPositive", `0`, 0.0),
			unmarshalSupported("ZeroNegative", `-0`, math.Float64frombits(0x8000000000000000)),
			unmarshalSupported("Positive", `1.7976931348623157e+308`, math.MaxFloat64),
			unmarshalSupported("Negative", `-1.7976931348623157e+308`, -math.MaxFloat64),
			unmarshalSupported("Smallest", `5e-324`, math.SmallestNonzeroFloat64),
			unmarshalUnsupportedValue[float64]("QNaN", `qnan`),
			unmarshalUnsupportedValue[float64]("SNaN", `snan`),
			unmarshalUnsupportedValue[float64]("PositiveInf", `inf`),
			unmarshalUnsupportedValue[float64]("NegativeInf", `-inf`),
		)
	})

	// Complex
	t.Run("Complex64", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalUnsupportedType[complex64]("BothZeroPositive", `(0.0,0.0)`),
		)
	})
	t.Run("Complex128", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalUnsupportedType[complex128]("BothZeroPositive", `(0.0,0.0)`),
		)
	})

	// Rune
	// NOTE: rune is effectively uint16
	t.Run("Rune", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported("Alpha", `65`, 'A'),
			unmarshalSupported("Unicode", `11935`, '\u2E9F'),
		)
	})

	// String
	t.Run("String", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported("Empty", `""`, ""),
			unmarshalSupported("NonEmpty", `"The quick brown fox"`, "The quick brown fox"),
		)
	})

	// Slice
	t.Run("Slice", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[[]string]("Null", `null`, nil),
			unmarshalSupported("Empty", `[]`, []string{}),
			unmarshalSupported("NonEmpty", `["The quick brown fox"]`, []string{"The quick brown fox"}),
		)
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		testUnmarshalJSON(t,
			unmarshalSupported[map[string]string]("Null", `null`, nil),
			unmarshalSupported("Empty", `{}`, map[string]string{}),
			unmarshalSupported("NonEmpty", `{"entry":"The quick brown fox"}`, map[string]string{"entry": "The quick brown fox"}),
		)
	})

	// Struct
	t.Run("Struct", func(t *testing.T) {
		t.Run("NotMarshalled", func(t *testing.T) {
			type testStructNotMarshalled struct {
				value int
			}
			var notMarshalled testStructNotMarshalled
			_ = notMarshalled.value
			testUnmarshalJSON(t,
				unmarshalSupported("Set", `{}`, notMarshalled),
			)
		})
		t.Run("Hidden", func(t *testing.T) {
			type testStructHidden struct {
				Hidden int `json:"-"`
			}
			var hidden testStructHidden
			testUnmarshalJSON(t,
				unmarshalSupported("Set", `{}`, hidden),
			)
		})
		t.Run("OneField", func(t *testing.T) {
			type testStructOneField struct {
				Value int `json:"value"`
			}
			var oneField testStructOneField
			testUnmarshalJSON(t,
				unmarshalSupported("Set", `{"value":0}`, oneField),
			)
		})
		t.Run("TwoFields", func(t *testing.T) {
			type testStructTwoFields struct {
				A int  `json:"a"`
				B bool `json:"b"`
			}
			var twoFields testStructTwoFields
			testUnmarshalJSON(t,
				unmarshalSupported("Set", `{"a":0,"b":false}`, twoFields),
			)
		})
		t.Run("EmbeddedOptional", func(t *testing.T) {
			type testStructEmbeddedOptional struct {
				Value optional.Value[int] `json:"value"`
			}
			var embeddedUnset testStructEmbeddedOptional
			embeddedSet := testStructEmbeddedOptional{
				Value: optional.NewValue(5),
			}
			testUnmarshalJSON(t,
				unmarshalSupported("SetValueUnset", `{}`, embeddedUnset),
				unmarshalSupported("SetValueSet", `{"value":5}`, embeddedSet),
			)
		})
	})
}
