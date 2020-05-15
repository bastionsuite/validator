package validator

import "testing"

func TestParseRequestForm(t *testing.T) {
}

func TestPaseRequestJSONBody(t *testing.T) {
}

func TestParseStruct(t *testing.T) {
	type target struct {
		TestString string `source:"Test" validate:"length=4-10,required"`
		TestInt    int    `source:"Number" validate:"range=0-10,required"`
	}

	source := struct {
		Test   string
		Number int
	}{
		"Test",
		2,
	}
	targ := target{}

	err := ParseStruct(&targ, source)
	if err != nil {
		t.Errorf("Error! %v", err)
	}
}
