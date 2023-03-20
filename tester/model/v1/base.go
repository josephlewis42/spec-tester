package v1

import (
	"regexp"

	"github.com/josephlewis42/scheme-compliance/tester/validation"

	"k8s.io/apimachinery/pkg/labels"
)

const (
	APIVersion  = "compliancetest/v1"
	NameMatcher = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"
)

var qualifiedNameRegexp = regexp.MustCompile("^" + NameMatcher + "$")

type Base struct {
	Version  `json:",inline"`
	Metadata Metadata `json:"metadata"`
}

func (b *Base) Validate(validator *validation.Validator) {
	b.Version.Validate(validator)
	validator.WithField("metadata", b.Metadata.Validate)
}

var _ validation.Validatable = (*Base)(nil)

type Version struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind"`
}

func (v *Version) Validate(validator *validation.Validator) {
	validator.WithField("apiVersion", func(validator *validation.Validator) {
		validation.AssertEqual(validator, APIVersion, v.APIVersion)
	})

	validator.WithField("kind", func(validator *validation.Validator) {
		validation.AssertNotBlank(validator, v.Kind)
	})
}

var _ validation.Validatable = (*Version)(nil)

type Metadata struct {
	// Unique name for the object.
	Name string `json:"name,omitempty"`

	DisplayMetadata `json:",inline"`
}

func (m *Metadata) Validate(validator *validation.Validator) {
	validator.WithField("name", func(validator *validation.Validator) {
		validation.AssertNotBlank(validator, m.Name)

		if nameLen := len(m.Name); nameLen > 63 {
			validator.Error("is %d characters long, must be <= 63", nameLen)
		}

		if !qualifiedNameRegexp.MatchString(m.Name) {
			validator.Error("must match %s", NameMatcher)
		}
	})

	m.DisplayMetadata.Validate(validator)
}

type DisplayMetadata struct {
	// Labels for the object. Will be cascaded from the parent.
	Labels Labels `json:"labels,omitempty"`
	// Human readable name for the object.
	DisplayName string `json:"displayName,omitempty"`
	// Description for the object. Will be shown to users.
	Description string `json:"description,omitempty"`
	// Note for the object. Won't be displayed to users.
	Note string `json:"note,omitempty"`
}

func (m *DisplayMetadata) Validate(validator *validation.Validator) {
	m.Labels.Validate(validator.Field("labels"))
}

var _ validation.Validatable = (*Metadata)(nil)

type Labels map[string]string

func (m Labels) Validate(validator *validation.Validator) {
	// TODO: add validation
}

func (m Labels) MergeOver(over Labels) Labels {
	out := make(Labels)
	for k, v := range over {
		out[k] = v
	}
	for k, v := range m {
		out[k] = v
	}
	return out
}

type LabelSelector string

func (ls LabelSelector) Validate(validator *validation.Validator) {
	if _, err := labels.Parse(string(ls)); err != nil {
		validator.Error("invalid label selector: %s", err.Error())
	}
}

var _ validation.Validatable = (LabelSelector)("")
