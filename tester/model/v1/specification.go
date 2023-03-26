package v1

import "github.com/josephlewis42/scheme-compliance/tester/validation"

const (
	TypeSpecification = "Specification"
)

// Specification represents a real world specifcation for example an RFC.
//
// Specifications are made up of several sections which themselves can be sub-divided.
// Each section can be marked as optional, and may have a selector associated that
// indicates which tests are used to assert compliance.
type Specification struct {
	Base `json:",inline"`

	Sections              []SpecificationSection `json:"sections"`
	ExclusionTestSelector LabelSelector          `json:"exclusionTestSelector,omitempty"`
}

var _ validation.Validatable = (*Specification)(nil)

// Validate implements validation.Validatable
func (spec *Specification) Validate(validator *validation.Validator) {
	spec.Base.Validate(validator)

	validator.WithField("sections", func(validator *validation.Validator) {
		if len(spec.Sections) == 0 {
			validator.Error("must supply at least one section")
		}

		for idx, section := range spec.Sections {
			section.Validate(validator.AtIndex(idx))
		}
	})

	validator.WithField("exclusionTestSelector", func(validator *validation.Validator) {
	})
}

// Tidy cleans up the structure to remove validation warnings.
func (impl *Specification) Tidy() {
	// no-op for specifications.
}

type SpecificationSection struct {
	Metadata     Metadata               `json:"metadata"`
	TestSelector LabelSelector          `json:"testSelector"`
	Optional     bool                   `json:"optional,omitempty"`
	Sections     []SpecificationSection `json:"sections,omitempty"`
}

var _ validation.Validatable = (*Specification)(nil)

// Validate implements validation.Validatable
func (section *SpecificationSection) Validate(validator *validation.Validator) {
	section.Metadata.Validate(validator.Field("metadata"))

	section.TestSelector.Validate(validator.Field("testSelctor"))

	validator.WithField("sections", func(validator *validation.Validator) {
		for idx, section := range section.Sections {
			section.Validate(validator.AtIndex(idx))
		}
	})
}

// HydratedSection is a runtime representation of a report section.
type HydratedSection struct {
	Metadata
	TestSelector LabelSelector
	Optional     bool
	Depth        int
}
