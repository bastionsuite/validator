package validator

import (
	"testing"
)

func TestTagParsing(t *testing.T) {
	opts, err := parseTag(`validate:"length=0-10" source:"test"`)
	if err != nil {
		t.Errorf("Problem with parsing tags %v", err)
		return
	}

	if opts.Source != "test" {
		t.Errorf("parseTag failed, source %v, expected 'source'", opts.Source)
	}
}
