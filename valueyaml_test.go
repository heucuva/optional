package optional_test

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/heucuva/optional"
	"gopkg.in/yaml.v2"
)

type marshalTestYAML[T any] struct {
	test     string
	value    optional.Value[T]
	expected string
	run      func(*testing.T)
}

func (ti marshalTestYAML[T]) runSupported(t *testing.T) {
	t.Helper()
	blob, err := yaml.Marshal(&ti.value)
	if err != nil {
		t.Fatal(err)
	}
	if observed := string(blob); strings.Compare(ti.expected, observed) != 0 {
		t.Fatalf("expected %q, got %q", ti.expected, observed)
	}
}

func (ti marshalTestYAML[T]) runUnsupportedValue(t *testing.T) {
	t.Helper()
	_, err := yaml.Marshal(&ti.value)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedValue *yaml.TypeError
	if !errors.As(err, &unsupportedValue) {
		t.Fatal(err)
	}
}

func (ti marshalTestYAML[T]) runUnsupportedType(t *testing.T) {
	t.Helper()
	_, err := yaml.Marshal(&ti.value)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedType *yaml.TypeError
	if !errors.As(err, &unsupportedType) {
		t.Fatal(err)
	}
}

func marshalSupportedYAML[T any](name string, value T, expected string) marshalTestYAML[T] {
	ti := marshalTestYAML[T]{
		test:     name,
		value:    optional.NewValue(value),
		expected: expected + "\n",
	}
	ti.run = ti.runSupported
	return ti
}

func marshalUnsupportedYAMLValue[T any](name string, value T) marshalTestYAML[T] {
	ti := marshalTestYAML[T]{
		test:  name,
		value: optional.NewValue(value),
	}
	ti.run = ti.runUnsupportedValue
	return ti
}

func marshalUnsupportedYAMLType[T any](name string, value T) marshalTestYAML[T] {
	ti := marshalTestYAML[T]{
		test:  name,
		value: optional.NewValue(value),
	}
	ti.run = ti.runUnsupportedType
	return ti
}

func testMarshalYAML[T any](t *testing.T, tests ...marshalTestYAML[T]) {
	t.Helper()

	t.Run("Unset", marshalTestYAML[T]{expected: "\n"}.runSupported)

	for _, ti := range tests {
		t.Run(ti.test, ti.run)
	}
}

func TestMarshalYAML(t *testing.T) {
	// TODO: fix up these tests
	// They're copy-pasted form JSON, so they probably are wrong.
	t.SkipNow()
	// Boolean
	t.Run("Bool", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML("True", true, `true`),
			marshalSupportedYAML("False", false, `false`),
		)
	})

	// Signed Integer
	t.Run("Int", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML("Zero", 0, `0`),
			marshalSupportedYAML("Positive", math.MaxInt, `9223372036854775807`),
			marshalSupportedYAML("Negative", math.MinInt, `-9223372036854775808`),
		)
	})
	t.Run("Int8", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[int8]("Zero", 0, `0`),
			marshalSupportedYAML[int8]("Positive", math.MaxInt8, `127`),
			marshalSupportedYAML[int8]("Negative", math.MinInt8, `-128`),
		)
	})
	t.Run("Int16", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[int16]("Zero", 0, `0`),
			marshalSupportedYAML[int16]("Positive", math.MaxInt16, `32767`),
			marshalSupportedYAML[int16]("Negative", math.MinInt16, `-32768`),
		)
	})
	t.Run("Int32", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[int32]("Zero", 0, `0`),
			marshalSupportedYAML[int32]("Positive", math.MaxInt32, `2147483647`),
			marshalSupportedYAML[int32]("Negative", math.MinInt32, `-2147483648`),
		)
	})

	// Unsigned integer
	t.Run("Uint", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[uint]("Zero", 0, `0`),
			marshalSupportedYAML[uint]("Max", math.MaxUint, `18446744073709551615`),
		)
	})
	t.Run("Uint8", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[uint8]("Zero", 0, `0`),
			marshalSupportedYAML[uint8]("Max", math.MaxUint8, `255`),
		)
	})
	t.Run("Uint16", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[uint16]("Zero", 0, `0`),
			marshalSupportedYAML[uint16]("Max", math.MaxUint16, `65535`),
		)
	})
	t.Run("Uint32", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[uint32]("Zero", 0, `0`),
			marshalSupportedYAML[uint32]("Max", math.MaxUint32, `4294967295`),
		)
	})

	// Floating point
	t.Run("Float32", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[float32]("ZeroPositive", 0.0, `0`),
			marshalSupportedYAML("ZeroNegative", math.Float32frombits(0x80000000), `-0`),
			marshalSupportedYAML[float32]("Positive", math.MaxFloat32, `3.4028235e+38`),
			marshalSupportedYAML[float32]("Negative", -math.MaxFloat32, `-3.4028235e+38`),
			marshalSupportedYAML[float32]("Smallest", math.SmallestNonzeroFloat32, `1e-45`),
			marshalUnsupportedYAMLValue("QNaN", math.Float32frombits(0x7FFFFFFF)),
			marshalUnsupportedYAMLValue("SNaN", math.Float32frombits(0x7FbFFFFF)),
			marshalUnsupportedYAMLValue("PositiveInf", math.Float32frombits(0x7F800000)),
			marshalUnsupportedYAMLValue("NegativeInf", math.Float32frombits(0xFF800000)),
		)
	})
	t.Run("Float64", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML("ZeroPositive", 0.0, `0`),
			marshalSupportedYAML("ZeroNegative", math.Float64frombits(0x8000000000000000), `-0`),
			marshalSupportedYAML("Positive", math.MaxFloat64, `1.7976931348623157e+308`),
			marshalSupportedYAML("Negative", -math.MaxFloat64, `-1.7976931348623157e+308`),
			marshalSupportedYAML("Smallest", math.SmallestNonzeroFloat64, `5e-324`),
			marshalUnsupportedYAMLValue("QNaN", math.Float64frombits(0x7FFFFFFFFFFFFFFF)),
			marshalUnsupportedYAMLValue("SNaN", math.Float64frombits(0x7FF7FFFFFFFFFFFF)),
			marshalUnsupportedYAMLValue("PositiveInf", math.Float64frombits(0x7FF0000000000000)),
			marshalUnsupportedYAMLValue("NegativeInf", math.Float64frombits(0xFFF0000000000000)),
		)
	})

	// Complex
	t.Run("Complex64", func(t *testing.T) {
		testMarshalYAML(t,
			marshalUnsupportedYAMLType("BothZeroPositive", complex(float32(0.0), float32(0.0))),
		)
	})
	t.Run("Complex128", func(t *testing.T) {
		testMarshalYAML(t,
			marshalUnsupportedYAMLType("BothZeroPositive", complex(float64(0.0), float64(0.0))),
		)
	})

	// Rune
	// NOTE: rune is effectively uint16
	t.Run("Rune", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML("Alpha", 'A', `65`),
			marshalSupportedYAML("Unicode", '\u2E9F', `11935`),
		)
	})

	// String
	t.Run("String", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML("Empty", "", `""`),
			marshalSupportedYAML("NonEmpty", "The quick brown fox", `"The quick brown fox"`),
		)
	})

	// Slice
	t.Run("Slice", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[[]string]("Null", nil, `null`),
			marshalSupportedYAML("Empty", []string{}, `[]`),
			marshalSupportedYAML("NonEmpty", []string{"The quick brown fox"}, `["The quick brown fox"]`),
		)
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		testMarshalYAML(t,
			marshalSupportedYAML[map[string]string]("Null", nil, `null`),
			marshalSupportedYAML("Empty", map[string]string{}, `{}`),
			marshalSupportedYAML("NonEmpty", map[string]string{"entry": "The quick brown fox"}, `{"entry":"The quick brown fox"}`),
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
			testMarshalYAML(t,
				marshalSupportedYAML("Set", notMarshalled, `{}`),
			)
		})
		t.Run("Hidden", func(t *testing.T) {
			type testStructHidden struct {
				Hidden int `yaml:"-"`
			}
			var hidden testStructHidden
			testMarshalYAML(t,
				marshalSupportedYAML("Set", hidden, `{}`),
			)
		})
		t.Run("OneField", func(t *testing.T) {
			type testStructOneField struct {
				Value int `yaml:"value"`
			}
			var oneField testStructOneField
			testMarshalYAML(t,
				marshalSupportedYAML("Set", oneField, `{"value":0}`),
			)
		})
		t.Run("TwoFields", func(t *testing.T) {
			type testStructTwoFields struct {
				A int  `yaml:"a"`
				B bool `yaml:"b"`
			}
			var twoFields testStructTwoFields
			testMarshalYAML(t,
				marshalSupportedYAML("Set", twoFields, `{"a":0,"b":false}`),
			)
		})
		t.Run("EmbeddedOptional", func(t *testing.T) {
			type testStructEmbeddedOptional struct {
				Value optional.Value[int] `yaml:"value"`
			}
			var embeddedUnset testStructEmbeddedOptional
			embeddedSet := testStructEmbeddedOptional{
				Value: optional.NewValue(5),
			}
			testMarshalYAML(t,
				marshalSupportedYAML("SetValueUnset", embeddedUnset, `{"value":null}`),
				marshalSupportedYAML("SetValueSet", embeddedSet, `{"value":5}`),
			)
		})
	})
}

type unmarshalTestYAML[T any] struct {
	test     string
	data     string
	comparer func(observed optional.Value[T]) (optional.Value[T], bool)
	run      func(*testing.T)
}

func (ti unmarshalTestYAML[T]) runSupported(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := yaml.Unmarshal([]byte(ti.data), &observed)
	if err != nil {
		t.Fatal(err)
	}
	if expected, success := ti.comparer(observed); !success {
		t.Fatalf("expected %+v, got %+v", expected, observed)
	}
}

func (ti unmarshalTestYAML[T]) runUnsupportedValue(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := yaml.Unmarshal([]byte(ti.data), &observed)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedValue *yaml.TypeError
	if !errors.As(err, &unsupportedValue) {
		t.Fatal(err)
	}
}

func (ti unmarshalTestYAML[T]) runUnsupportedType(t *testing.T) {
	t.Helper()
	var observed optional.Value[T]
	err := yaml.Unmarshal([]byte(ti.data), &observed)
	if err == nil {
		t.Fatal("expected serialization failure, but got success")
	}
	var unsupportedType *yaml.TypeError
	if !errors.As(err, &unsupportedType) {
		t.Fatal(err)
	}
}

func unmarshalSupportedYaml[T any](name string, data string, value T) unmarshalTestYAML[T] {
	ti := unmarshalTestYAML[T]{
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

func unmarshalUnsupportedYAMLValue[T any](name string, data string) unmarshalTestYAML[T] {
	ti := unmarshalTestYAML[T]{
		test: name,
		data: data,
	}
	ti.run = ti.runUnsupportedValue
	return ti
}

func unmarshalUnsupportedYAMLType[T any](name string, data string) unmarshalTestYAML[T] {
	ti := unmarshalTestYAML[T]{
		test: name,
		data: data,
	}
	ti.run = ti.runUnsupportedType
	return ti
}

func testUnmarshalYAML[T any](t *testing.T, tests ...unmarshalTestYAML[T]) {
	t.Helper()

	t.Run("Unset", marshalTestYAML[T]{expected: `{}`}.runSupported)

	for _, ti := range tests {
		t.Run(ti.test, ti.run)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	// TODO: fix up these tests
	// They're copy-pasted form JSON, so they probably are wrong.
	t.SkipNow()

	// Boolean
	t.Run("Bool", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml("True", `true`, true),
			unmarshalSupportedYaml("False", `false`, false),
		)
	})

	// Signed Integer
	t.Run("Int", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml("Zero", `0`, 0),
			unmarshalSupportedYaml("Positive", `9223372036854775807`, math.MaxInt),
			unmarshalSupportedYaml("Negative", `-9223372036854775808`, math.MinInt),
		)
	})
	t.Run("Int8", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[int8]("Zero", `0`, 0),
			unmarshalSupportedYaml[int8]("Positive", `127`, math.MaxInt8),
			unmarshalSupportedYaml[int8]("Negative", `-128`, math.MinInt8),
		)
	})
	t.Run("Int16", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[int16]("Zero", `0`, 0),
			unmarshalSupportedYaml[int16]("Positive", `32767`, math.MaxInt16),
			unmarshalSupportedYaml[int16]("Negative", `-32768`, math.MinInt16),
		)
	})
	t.Run("Int32", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[int32]("Zero", `0`, 0),
			unmarshalSupportedYaml[int32]("Positive", `2147483647`, math.MaxInt32),
			unmarshalSupportedYaml[int32]("Negative", `-2147483648`, math.MinInt32),
		)
	})

	// Unsigned integer
	t.Run("Uint", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[uint]("Zero", `0`, 0),
			unmarshalSupportedYaml[uint]("Max", `18446744073709551615`, math.MaxUint),
		)
	})
	t.Run("Uint8", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[uint8]("Zero", `0`, 0),
			unmarshalSupportedYaml[uint8]("Max", `255`, math.MaxUint8),
		)
	})
	t.Run("Uint16", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[uint16]("Zero", `0`, 0),
			unmarshalSupportedYaml[uint16]("Max", `65535`, math.MaxUint16),
		)
	})
	t.Run("Uint32", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[uint32]("Zero", `0`, 0),
			unmarshalSupportedYaml[uint32]("Max", `4294967295`, math.MaxUint32),
		)
	})

	// Floating point
	t.Run("Float32", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[float32]("ZeroPositive", `0`, 0.0),
			unmarshalSupportedYaml("ZeroNegative", `-0`, math.Float32frombits(0x80000000)),
			unmarshalSupportedYaml[float32]("Positive", `3.4028235e+38`, math.MaxFloat32),
			unmarshalSupportedYaml[float32]("Negative", `-3.4028235e+38`, -math.MaxFloat32),
			unmarshalSupportedYaml[float32]("Smallest", `1e-45`, math.SmallestNonzeroFloat32),
			unmarshalUnsupportedYAMLValue[float32]("QNaN", `qnan`),
			unmarshalUnsupportedYAMLValue[float32]("SNaN", `snan`),
			unmarshalUnsupportedYAMLValue[float32]("PositiveInf", `inf`),
			unmarshalUnsupportedYAMLValue[float32]("NegativeInf", `-inf`),
		)
	})
	t.Run("Float64", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml("ZeroPositive", `0`, 0.0),
			unmarshalSupportedYaml("ZeroNegative", `-0`, math.Float64frombits(0x8000000000000000)),
			unmarshalSupportedYaml("Positive", `1.7976931348623157e+308`, math.MaxFloat64),
			unmarshalSupportedYaml("Negative", `-1.7976931348623157e+308`, -math.MaxFloat64),
			unmarshalSupportedYaml("Smallest", `5e-324`, math.SmallestNonzeroFloat64),
			unmarshalUnsupportedYAMLValue[float64]("QNaN", `qnan`),
			unmarshalUnsupportedYAMLValue[float64]("SNaN", `snan`),
			unmarshalUnsupportedYAMLValue[float64]("PositiveInf", `inf`),
			unmarshalUnsupportedYAMLValue[float64]("NegativeInf", `-inf`),
		)
	})

	// Complex
	t.Run("Complex64", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalUnsupportedYAMLType[complex64]("BothZeroPositive", `(0.0,0.0)`),
		)
	})
	t.Run("Complex128", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalUnsupportedYAMLType[complex128]("BothZeroPositive", `(0.0,0.0)`),
		)
	})

	// Rune
	// NOTE: rune is effectively uint16
	t.Run("Rune", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml("Alpha", `65`, 'A'),
			unmarshalSupportedYaml("Unicode", `11935`, '\u2E9F'),
		)
	})

	// String
	t.Run("String", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml("Empty", `""`, ""),
			unmarshalSupportedYaml("NonEmpty", `"The quick brown fox"`, "The quick brown fox"),
		)
	})

	// Slice
	t.Run("Slice", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[[]string]("Null", `null`, nil),
			unmarshalSupportedYaml("Empty", `[]`, []string{}),
			unmarshalSupportedYaml("NonEmpty", `["The quick brown fox"]`, []string{"The quick brown fox"}),
		)
	})

	// Map
	t.Run("Map", func(t *testing.T) {
		testUnmarshalYAML(t,
			unmarshalSupportedYaml[map[string]string]("Null", `null`, nil),
			unmarshalSupportedYaml("Empty", `{}`, map[string]string{}),
			unmarshalSupportedYaml("NonEmpty", `{"entry":"The quick brown fox"}`, map[string]string{"entry": "The quick brown fox"}),
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
			testUnmarshalYAML(t,
				unmarshalSupportedYaml("Set", `{}`, notMarshalled),
			)
		})
		t.Run("Hidden", func(t *testing.T) {
			type testStructHidden struct {
				Hidden int `yaml:"-"`
			}
			var hidden testStructHidden
			testUnmarshalYAML(t,
				unmarshalSupportedYaml("Set", `{}`, hidden),
			)
		})
		t.Run("OneField", func(t *testing.T) {
			type testStructOneField struct {
				Value int `yaml:"value"`
			}
			var oneField testStructOneField
			testUnmarshalYAML(t,
				unmarshalSupportedYaml("Set", `{"value":0}`, oneField),
			)
		})
		t.Run("TwoFields", func(t *testing.T) {
			type testStructTwoFields struct {
				A int  `yaml:"a"`
				B bool `yaml:"b"`
			}
			var twoFields testStructTwoFields
			testUnmarshalYAML(t,
				unmarshalSupportedYaml("Set", `{"a":0,"b":false}`, twoFields),
			)
		})
		t.Run("EmbeddedOptional", func(t *testing.T) {
			type testStructEmbeddedOptional struct {
				Value optional.Value[int] `yaml:"value"`
			}
			var embeddedUnset testStructEmbeddedOptional
			embeddedSet := testStructEmbeddedOptional{
				Value: optional.NewValue(5),
			}
			testUnmarshalYAML(t,
				unmarshalSupportedYaml("SetValueUnset", `{}`, embeddedUnset),
				unmarshalSupportedYaml("SetValueSet", `{"value":5}`, embeddedSet),
			)
		})
	})
}
