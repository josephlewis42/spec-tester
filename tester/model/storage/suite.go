package storage

import (
	"fmt"
	"io"

	"github.com/josephlewis42/scheme-compliance/tester/executor"
	v1 "github.com/josephlewis42/scheme-compliance/tester/model/v1"
	"github.com/josephlewis42/scheme-compliance/tester/validation"
)

// Suite represents the full set of files that make up a test suite on-disk.
type Suite struct {
	RootPath string

	Implementations []YamlFile[v1.Implementation]
	Specifications  []YamlFile[v1.Specification]
	Tests           []YamlFile[v1.Test]
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

	// Exclusions should not overlap with inclusions
}

// Tidy cleans up the structure to remove validation warnings.
func (s *Suite) Tidy() {
	for _, impl := range s.Implementations {
		impl.Value.Tidy()
	}

	for _, spec := range s.Specifications {
		spec.Value.Tidy()
	}

	for _, test := range s.Tests {
		test.Value.Tidy()
	}
}

// Save updates the files on disk.
func (s *Suite) Save() {
	for _, impl := range s.Implementations {
		impl.Save()
	}

	for _, spec := range s.Specifications {
		spec.Save()
	}

	for _, test := range s.Tests {
		test.Save()
	}
}

// Diff returns a human-readable diff for the whole suite compared to its stored version.
// The output isn't stable.
func (s *Suite) Diff(w io.Writer) {
	for _, impl := range s.Implementations {
		fmt.Fprint(w, impl.Diff())
	}

	for _, spec := range s.Specifications {
		fmt.Fprint(w, spec.Diff())
	}

	for _, test := range s.Tests {
		fmt.Fprint(w, test.Diff())
	}
}

func (s *Suite) ListTests() (hydrated []*executor.TestCase) {
	for _, t := range s.Tests {
		t.Value.WalkCases(func(test *executor.TestCase) {
			hydrated = append(hydrated, test)
		})
	}

	return
}

func (s *Suite) ListImplementations() (hydrated []*executor.Implementation) {
	for _, impl := range s.Implementations {
		hydrated = append(hydrated, impl.Value.ConvertToInternal())
	}

	return
}

func (s *Suite) ListSpecifications() (hydrated []*executor.Specification) {
	for _, spec := range s.Specifications {
		hydrated = append(hydrated, spec.Value.ConvertToInternal())
	}
	return
}
