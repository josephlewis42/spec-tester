package v1

import "github.com/josephlewis42/scheme-compliance/tester/validation"

const (
	TypeImplementation = "Implementation"
)

type Implementation struct {
	Base `json:",inline"`

	Variants []ImplementationVariant `json:"variants"`
}

var _ validation.Validatable = (*Implementation)(nil)

func (impl *Implementation) Validate(validator *validation.Validator) {
	impl.Base.Validate(validator)

	validator.WithField("variants", func(validator *validation.Validator) {
		if len(impl.Variants) == 0 {
			validator.Error("must supply at least one variant")
		}

		for idx, variant := range impl.Variants {
			variant.Validate(validator.AtIndex(idx))
		}
	})
}

type ImplementationVariant struct {
	Metadata       Metadata             `json:"metadata"`
	Runtime        ImplementationSource `json:"runtime"`
	Specifications []string             `json:"specifications"`

	// Command to run, $(PROGRAM) and $(PROGRAM_PATH) will be replaced
	TestCommand []string `json:"testCommand"`
}

var _ validation.Validatable = (*ImplementationVariant)(nil)

func (impl *ImplementationVariant) Validate(validator *validation.Validator) {
	validator.WithField("metadata", impl.Metadata.Validate)

	validator.WithField("runtime", impl.Runtime.Validate)

	validator.WithField("specifications", func(validator *validation.Validator) {
		if len(impl.Specifications) == 0 {
			validator.Error("must reference at least one specification")
		}
	})

	validator.WithField("testCommand", func(validator *validation.Validator) {
		if len(impl.Specifications) == 0 {
			validator.Error("must reference at least one specification")
		}
	})
}

type ImplementationSource struct {
	// Image      ImplementationSourceImage `json:"image,omitempty"`
	// Dockerfile ImplementationSourceBuild `json:"dockerfile,omitempty"`
	Local *ImplementationSourceLocal `json:"local,omitempty"`
}

var _ validation.Validatable = (*ImplementationSource)(nil)

func (impl *ImplementationSource) Validate(validator *validation.Validator) {
	validation.OneOf().
		ValidatedField("local", impl.Local != nil, impl.Local.Validate).
		Validate(validator)

}

type ImplementationSourceLocal struct {
}

var _ validation.Validatable = (*ImplementationSourceLocal)(nil)

func (impl *ImplementationSourceLocal) Validate(validator *validation.Validator) {
	// no-op
}
