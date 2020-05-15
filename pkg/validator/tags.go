// Copyright 2020 BastionSuite. All rights resreved.
// Use of this source code is governed by a MIT license that can be
// found in the LICENSE file

// Most of this is an interpretation of the golang json tags
// implementation. Because that's exactly what we need, but their
// implementation is private

package validator

import (
	"errors"
	"regexp"
)

var ErrInvalidQualifier = errors.New("validator/tags: invalid qualifier")
var ErrDuplicateQualifier = errors.New("validator/tags: duplicate qualifier")

type ValidationOption struct {
	Name     string
	Required bool

	MinLength int
	MaxLength int
	MinRange  int
	MaxRange  int
}

type ValidationOptions struct {
	Source      string
	Validations []ValidationOption
}

func parseTag(tag string) (ValidationOptions, error) {
	r := ValidationOptions{}

	re, err := regexp.Compile(`([a-z]+):"([^"]+)"`)
	if err != nil {
		return r, err
	}

	qualifiers := re.FindAllStringSubmatch(tag, -1)
	hadSource := false
	hadValidations := false
	for _, q := range qualifiers {
		switch {
		case q[1] == "source":
			if hadSource {
				return r, ErrDuplicateQualifier
			}
			r.Source = q[2]
			hadSource = true
		case q[1] == "validate":
			if hadValidations {
				return r, ErrDuplicateQualifier
			}
			hadValidations = true
			validations, err := parseValidations(q[2])
			if err != nil {
				return r, err
			}
			r.Validations = validations
		default:
			return r, ErrInvalidQualifier
		}
	}

	return r, nil
}

func parseValidations(validations string) ([]ValidationOption, error) {
	v := make([]ValidationOption, 0)

	return v, nil
}
