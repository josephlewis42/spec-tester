package v1

import (
	"fmt"
	"reflect"

	"github.com/josephlewis42/scheme-compliance/tester/validation"
)

const (
	KindTest = "Test"
)

type Test struct {
	Base        `json:",inline"`
	TestContext `json:",inline"`
}

var _ validation.Validatable = (*Test)(nil)

func (t *Test) effectiveTemplate() HydratedTestCaseTemplate {
	return HydratedTestCaseTemplate{
		Path: "/" + t.Metadata.Name,
		TestCaseTemplate: TestCaseTemplate{
			Labels: t.Metadata.Labels,
		},
	}
}

func (t *Test) WalkCases(callback func(HydratedTestCase)) {
	t.TestContext.WalkCases(t.effectiveTemplate(), callback)
}

func (t *Test) Validate(validator *validation.Validator) {
	t.Base.Validate(validator)
	validation.AssertEqual(validator.Field("kind"), KindTest, t.Base.Kind)

	validator.WithField("tests", func(validator *validation.Validator) {
		t.TestContext.ValidateEffective(validator, t.effectiveTemplate())
	})
}

// Tidy cleans up the structure to remove validation warnings.
func (t *Test) Tidy() {
	t.TestContext.Tidy()
}

// TestContext holds a set of related tests and a template that can be applied to them.
type TestContext struct {
	Template TestCaseTemplate    `json:"template,omitempty"`
	Tests    []TestContextOrCase `json:"tests"`
}

// Tidy cleans up the structure to remove validation warnings.
func (t *TestContext) Tidy() {
	for _, test := range t.Tests {
		test.Tidy()
	}
}

// WalkCases executes callback for each child test case that's had the template applied.
func (t *TestContext) WalkCases(parent HydratedTestCaseTemplate, callback func(HydratedTestCase)) {
	parent = t.Template.Hydrate(parent.WithPathSuffix("/tests"))

	for idx, entry := range t.Tests {
		entry.WalkCases(parent.WithPathSuffix(fmt.Sprintf("/%d", idx)), callback)
	}
}

func (t *TestContext) ValidateEffective(validator *validation.Validator, parent HydratedTestCaseTemplate) {
	validator.WithField("template", func(validator *validation.Validator) {
		// TODO: Decide which fields shouldn't be set in the temlpate (if any)
		t.Template.Validate(validator)
	})

	parent = t.Template.Hydrate(parent)
	validator.WithField("tests", func(validator *validation.Validator) {
		for idx, value := range t.Tests {
			value.ValidateEffective(validator.AtIndex(idx), parent)
		}
	})
}

// TestSectionOrCase is a union of section and case only one field may be set.
type TestContextOrCase struct {
	Context *TestContext `json:"context,omitempty"`
	Case    *TestCase    `json:"case,omitempty"`
}

// Tidy cleans up the structure to remove validation warnings.
func (t *TestContextOrCase) Tidy() {
	if t.Case != nil {
		t.Case.Tidy()
	}
	if t.Context != nil {
		t.Context.Tidy()
	}
}

// WalkCases executes callback for each child test case that's had the template applied.
func (t *TestContextOrCase) WalkCases(parent HydratedTestCaseTemplate, callback func(HydratedTestCase)) {
	if t.Case != nil {
		callback(t.Case.Hydrate(parent))
	}
	if t.Context != nil {
		t.Context.WalkCases(parent, callback)
	}
}

func (t *TestContextOrCase) ValidateEffective(validator *validation.Validator, parent HydratedTestCaseTemplate) {
	validation.OneOf().
		ValidatedField("context", t.Context != nil, func(validatior *validation.Validator) {
			t.Context.ValidateEffective(validator, parent)
		}).
		ValidatedField("case", t.Case != nil, func(validatior *validation.Validator) {
			t.Case.ValidateEffective(validator, parent)
		}).
		Validate(validator)
}

type TestCaseTemplate struct {
	DisplayName StringMutator `json:"displayName,omitempty"`
	Description StringMutator `json:"description,omitempty"`

	Labels Labels           `json:"labels,omitempty"`
	Expect *TestExpectation `json:"expect,omitempty"`
}

func (template *TestCaseTemplate) Validate(validator *validation.Validator) {
	validator.WithField("labels", template.Labels.Validate)

	if expect := template.Expect; expect != nil {
		validator.WithField("expect", expect.Validate)
	}
}

func (template *TestCaseTemplate) Hydrate(parent HydratedTestCaseTemplate) (hydrated HydratedTestCaseTemplate) {
	hydrated.Path = parent.Path
	hydrated.Description = template.Description.MergeOver(parent.Description)
	hydrated.DisplayName = template.DisplayName.MergeOver(parent.DisplayName)

	hydrated.Labels = template.Labels.MergeOver(parent.Labels)
	hydrated.Expect = coalesce(template.Expect, parent.Expect)

	return
}

type HydratedTestCaseTemplate struct {
	Path string
	TestCaseTemplate
}

func (template *HydratedTestCaseTemplate) WithPathSuffix(suffix string) HydratedTestCaseTemplate {
	copy := *template
	copy.Path += suffix
	return copy
}

type StringMutator struct {
	Prefix string `json:"prefix,omitempty"`
	Suffix string `json:"suffix,omitempty"`
	Value  string `json:"value,omitempty"`
}

func (sm *StringMutator) Apply(input string) string {
	if input == "" {
		input = sm.Value
	}

	return sm.Prefix + input + sm.Suffix
}

func (sm *StringMutator) MergeOver(base StringMutator) (merged StringMutator) {
	merged.Prefix = base.Prefix + sm.Prefix
	merged.Value = coalesce(sm.Value, base.Value)
	merged.Suffix = sm.Suffix + base.Suffix
	return
}

type TestCase struct {
	IdentifiableMetadata `json:",inline"`
	Input                *string          `json:"input"`
	Expect               *TestExpectation `json:"expect,omitempty"`
	Skip                 *string          `json:"skip,omitempty"`
}

// Tidy cleans up the structure to remove validation warnings.
func (t *TestCase) Tidy() {
	t.IdentifiableMetadata.Tidy()
}

func (t *TestCase) ValidateEffective(validator *validation.Validator, parent HydratedTestCaseTemplate) {
	effectiveTestCase := t.Hydrate(parent)
	effectiveTestCase.DisplayMetadata.Validate(validator)

	validator.WithField("input", func(validator *validation.Validator) {
		if effectiveTestCase.Input == nil {
			validator.Error("must be defined")
		} else {
			validation.AssertNotBlank(validator, *effectiveTestCase.Input)
		}
	})

	validator.WithField("expect", func(validator *validation.Validator) {
		if effectiveTestCase.Expect == nil {
			validator.Error("must be defined")
		} else {
			effectiveTestCase.Expect.Validate(validator)
		}
	})

	validator.WithField("skip", func(validator *validation.Validator) {
		if effectiveTestCase.Skip != nil {
			validator.Warning("test is skipped, reason: %q", *t.Skip)
		}
	})
}

// IsSkipped checks whether this test case should be skipped.
func (tc *TestCase) IsSkipped() bool {
	return tc.Skip != nil
}

// Hydrate merges properties from a parent test case into this one
// to produce a single test output.
func (tc *TestCase) Hydrate(parent HydratedTestCaseTemplate) HydratedTestCase {
	var out TestCase
	// Metadata fields
	out.DisplayMetadata.DisplayName = parent.DisplayName.Apply(tc.DisplayMetadata.DisplayName)
	out.DisplayMetadata.Description = parent.Description.Apply(tc.DisplayMetadata.Description)

	out.DisplayMetadata.Labels = out.DisplayMetadata.Labels.MergeOver(parent.Labels)

	// Test Fields
	out.Input = tc.Input
	out.Expect = coalesce(tc.Expect, parent.Expect)

	out.Skip = tc.Skip

	return HydratedTestCase{
		Path:     parent.Path,
		TestCase: out,
	}
}

type HydratedTestCase struct {
	Path     string `json:"path"`
	TestCase `json:",inline"`
}

func coalesce[T any](args ...T) (zero T) {
	for _, arg := range args {
		if reflect.ValueOf(arg).IsZero() {
			continue
		}
		return arg
	}

	return
}

type TestExpectation struct {
	// Indicates an exact value is required.
	Exact *string `json:"exact,omitempty"`
	// Indicates undefined behavior.
	Undefined *bool `json:"undefined,omitempty"`
}

func (t *TestExpectation) Validate(validator *validation.Validator) {
	validation.OneOf().
		Field("exact", t.Exact != nil).
		ValidatedField("undefined", t.Undefined != nil, func(validator *validation.Validator) {
			switch {
			case t.Undefined == nil:
				return
			case *t.Undefined == false:
				validator.Error("may only be true")
			}
		}).
		Validate(validator)
}
