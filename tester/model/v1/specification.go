package v1

import (
	"github.com/josephlewis42/scheme-compliance/tester/executor"
	"github.com/josephlewis42/scheme-compliance/tester/validation"
)

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
func (spec *Specification) Tidy() {
	for _, section := range spec.Sections {
		section.Tidy()
	}
}

func (spec *Specification) ConvertToInternal() *executor.Specification {
	internal := &executor.Specification{
		Metadata: spec.Metadata.ConvertToInternal(),
	}

	for _, childSection := range spec.Sections {
		internal.Sections = append(internal.Sections, childSection.ConvertToInternal())
	}

	return internal
}

type SpecificationSection struct {
	Metadata Metadata `json:"metadata"`
	Optional bool     `json:"optional,omitempty"`

	TestSelector LabelSelector          `json:"testSelector"`
	Sections     []SpecificationSection `json:"sections,omitempty"`
}

var _ validation.Validatable = (*Specification)(nil)

// Validate implements validation.Validatable
func (section *SpecificationSection) Validate(validator *validation.Validator) {
	section.Metadata.Validate(validator.Field("metadata"))

	section.TestSelector.Validate(validator.Field("testSelctor"))

	validation.OneOf().
		ValidatedField("testSelector", section.TestSelector != "", section.TestSelector.Validate).
		ValidatedField("sections", len(section.Sections) > 0, func(validator *validation.Validator) {
			for idx, section := range section.Sections {
				section.Validate(validator.AtIndex(idx))
			}
		}).Validate(validator)
}

// Tidy cleans up the structure to remove validation warnings.
func (section *SpecificationSection) Tidy() {
	for _, section := range section.Sections {
		section.Tidy()
	}
}

func (section *SpecificationSection) ConvertToInternal() *executor.SpecificationSection {
	internal := &executor.SpecificationSection{
		Metadata: section.Metadata.ConvertToInternal(),
	}

	switch {
	case section.TestSelector != "":
		internal.Content = &executor.SpecificationSection_TestSummary{
			TestSummary: &executor.SpecificationTestSummary{
				TestSelector: string(section.TestSelector),
			},
		}

	case len(section.Sections) > 0:
		var subsections []*executor.SpecificationSection
		for _, childSection := range section.Sections {
			subsections = append(subsections, childSection.ConvertToInternal())
		}

		internal.Content = &executor.SpecificationSection_SectionSummary{
			SectionSummary: &executor.SpecificationSectionSummary{
				Subsections: subsections,
			},
		}
	}

	return internal
}
