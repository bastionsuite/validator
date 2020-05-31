package validator

import (
	"net/url"
	"testing"
	"time"
)

func TestParseRequestForm(t *testing.T) {
	type target struct {
		TestString   string    `source:"test" validate:"length=4-10"`
		TestInt      int       `source:"number" validate:"range=0-10"`
		TestFloat    float64   `source:"floaty"`
		TestDateTime time.Time `source:"date"`
	}

	source := url.Values{}
	source.Set("test", "test")
	source.Set("number", "2")
	source.Set("floaty", "1.2")
	source.Set("date", "2020-06-30")

	targ := target{}
	err := ParseRequestForm(&targ, source)
	if err != nil {
		t.Errorf("Error: %+v", err)
	}

	if targ.TestString != "test" {
		t.Errorf("target.TestString = '%+v', expected 'test'", targ.TestString)
	}
	if targ.TestInt != 2 {
		t.Errorf("target.TestInt = '%+v', expected '2'", targ.TestInt)
	}
	if targ.TestFloat != 1.2 {
		t.Errorf("target.TestFloat = '%+v', expected '1.2'", targ.TestFloat)
	}
	refDT := time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC)
	if targ.TestDateTime != refDT {
		t.Errorf("target.TestFloat = '%+v', expected '2020-06-30'", targ.TestDateTime)
	}
}

func TestPaseRequestJSONBody(t *testing.T) {
	type target struct {
		TestString string    `source:"test" validate:"length=4-10"`
		TestInt    int       `source:"number" validate:"range=0-10"`
		TestFloat  float64   `source:"floaty"`
		TestDT     time.Time `source:"date"`
	}

	var jsonBlob = []byte(`{
	"test": "test",
	"number": 2,
	"floaty": 1.2,
	"date": "2020-06-30"}`)

	targ := target{}
	err := ParseJSONBody(&targ, jsonBlob)
	if err != nil {
		t.Errorf("Error: %+v", err)
	}

	if targ.TestString != "test" {
		t.Errorf("target.TestString = '%+v', expected 'test'", targ.TestString)
	}
	if targ.TestInt != 2 {
		t.Errorf("target.TestInt = '%+v', expected '2'", targ.TestInt)
	}
}
