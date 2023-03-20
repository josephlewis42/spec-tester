package storage

import (
	v1 "github.com/josephlewis42/scheme-compliance/tester/model/v1"
	"github.com/josephlewis42/scheme-compliance/tester/validation"
)

// Suite represents the full set of files that make up a test suite on-disk.
type Suite struct {
	RootPath string

	Implementations []File[v1.Implementation]
	Specifications  []File[v1.Specification]
	Tests           []File[v1.Test]
}

func (s *Suite) RunValidation(callback func(name string, v *validation.Validator)) {
	for _, impl := range s.Implementations {
		v := &validation.Validator{}
		impl.Value.Validate(v)
		callback(impl.Path, v)
	}

	for _, spec := range s.Specifications {
		v := &validation.Validator{}
		spec.Value.Validate(v)
		callback(spec.Path, v)
	}

	for _, test := range s.Tests {
		v := &validation.Validator{}
		test.Value.Validate(v)
		callback(test.Path, v)
	}

	// TODO: Add a global validation

	// Each item should have a unique name

	// Each test should be used

	// Each specification should explicitly include or exclude all tests

	// Exclusions should not overlap with inclusions

	// Implementations should have >= 1 version
}

func (s *Suite) HydrateTests() (hydrated []v1.HydratedTestCase) {

	for _, t := range s.Tests {
		t.Value.WalkCases(func(test v1.HydratedTestCase) {
			hydrated = append(hydrated, test)
		})
	}

	return
}

// File maps an on-disk file with a specific parsed value from that file.
type File[T any] struct {
	Path  string
	Value T
}
