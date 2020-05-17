package validator

import (
	"math"
	"testing"
)

func TestInvalidParseTag(t *testing.T) {
	// Test Error cases
	_, err := parseTag(`source:"Test" source:"Tast"`)
	if err != ErrDuplicateQualifier {
		t.Errorf("Expected duplicate qualifier on source. Got nothing")
	}

	_, err = parseTag(`validate:"length=0-10" validate:"length=10-20"`)
	if err != ErrDuplicateQualifier {
		t.Errorf("Expected duplicate qualifier on validate. Got nothing")
	}

	_, err = parseTag(`nope:"doesn't exist"`)
	if err != ErrInvalidQualifier {
		t.Errorf("Expected invalid qualifier on `nope'. Got nothing")
	}

	_, err = parseTag(`validate:"honk=12"`)
	if err != ErrInvalidValidation {
		t.Errorf("Expected invalid validation. Got nothing")
	}
}

func TestParseTag(t *testing.T) {
	opts, err := parseTag(`validate:"length=0-10" required source:"test"`)
	if err != nil {
		t.Errorf("Problem with parsing tags %v", err)
		return
	}

	if opts.Source != "test" {
		t.Errorf("parseTag failed, source %v, expected 'source'", opts.Source)
	}
}

func TestInvalidLengthValidations(t *testing.T) {
	var err error
	_, err = parseValidations("length=test-10")
	if err == nil {
		t.Errorf("Expected invalid min length validation. Got nothing")
	}

	_, err = parseValidations("length=test-test")
	if err == nil {
		t.Errorf("Expected invalid min and max length validation. Got nothing")
	}

	_, err = parseValidations("length=0-test")
	if err == nil {
		t.Errorf("Expected invalid max length validation. Got nothing")
	}
}

func TestValidLengthValidations(t *testing.T) {
	v, _ := parseValidations("length=5-10")

	if len(v) != 1 {
		t.Errorf("Expected one validation, got %v", len(v))
		return
	}

	testFunc := v[0]

	tables := []struct {
		Input  string
		Expect bool
	}{
		{"1234", false},
		{"12345", true},
		{"1234567890", true},
		{"12345678901", false},
	}

	for _, row := range tables {
		res := testFunc(row.Input)
		if res != row.Expect {
			t.Errorf("validation(%v) == %v, expected: %v", row.Input, res, row.Expect)
		}
	}
}

func TestInvalidRangeValidations(t *testing.T) {
	var err error
	_, err = parseValidations("range=t-10")
	if err == nil {
		t.Errorf("Expected invalid min range validation. Got nothing")
	}

	_, err = parseValidations("range=t-t")
	if err == nil {
		t.Errorf("Expected invalid min and max range validation. Got nothing")
	}

	_, err = parseValidations("range=0-t")
	if err == nil {
		t.Errorf("Expected invalid max range validation. Got nothing")
	}
}

func TestValidIntRangeValidations(t *testing.T) {
	v, _ := parseValidations("range=5-10")
	if len(v) != 1 {
		t.Errorf("Expected one validation, got %v", len(v))
		return
	}

	testFunc := v[0]

	tables := []struct {
		Input  interface{}
		Expect bool
	}{
		{1, false},
		{5, true},
		{10, true},
		{11, false},
		{"abc", false},
	}

	for _, row := range tables {
		res := testFunc(row.Input)
		if res != row.Expect {
			t.Errorf("validation(%v) == %v, expected: %v", row.Input, res, row.Expect)
		}
	}
}

func TestValidStringValidation(t *testing.T) {
	v, _ := parseValidations("string")
	if len(v) != 1 {
		t.Errorf("Expected one validation, got %v", len(v))
		return
	}

	testFunc := v[0]

	res := testFunc("test")
	if !res {
		t.Errorf("string validation should really always work")
	}
}

func TestValidIntValidation(t *testing.T) {
	v, _ := parseValidations("int")
	if len(v) != 1 {
		t.Errorf("Expected one validation, got %v", len(v))
		return
	}

	testFunc := v[0]

	tables := []struct {
		Input  interface{}
		Expect bool
	}{
		{1, true},
		{12, true},
		{1092323, true},
		{0xdeadbeaf, true},
		{"abc", false},
	}

	for _, row := range tables {
		res := testFunc(row.Input)
		if res != row.Expect {
			t.Errorf("validation(%v) == %v, expected: %v", row.Input, res, row.Expect)
		}
	}
}

func TestFloatValidation(t *testing.T) {
	v, _ := parseValidations("float")
	if len(v) != 1 {
		t.Errorf("Expected one validation, got %v", len(v))
		return
	}

	testFunc := v[0]

	tables := []struct {
		Input  interface{}
		Expect bool
	}{
		{1.0, true},
		{1., true},
		{0.123, true},
		{-0.123, true},
		{"falsdkjas", false},
		{math.NaN(), true},
		{math.Inf(1), true},
		{math.Inf(-1), true},
	}

	for _, row := range tables {
		res := testFunc(row.Input)
		if res != row.Expect {
			t.Errorf("validation(%v) == %v, expected: %v", row.Input, res, row.Expect)
		}
	}
}
