package v1

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/josephlewis42/scheme-compliance/tester/executor"
	"github.com/josephlewis42/scheme-compliance/tester/validation"
	"github.com/xeipuuv/gojsonschema"
)

const (
	KindTestSuite = "TestSuite"
)

type TestSuite struct {
	Base `json:",inline"`
	Spec TestSuiteSpec `json:"spec"`
}

var _ validation.Validatable = (*TestSuite)(nil)

func (suite *TestSuite) Tidy() {
	// no-op
}

func (suite *TestSuite) Validate(validator *validation.Validator) {
	suite.Base.Validate(validator)

	validator.WithField("spec", suite.Spec.Validate)
}

type TestSuiteSpec struct {
	Assertions TestSuiteSpecAssertionConfig `json:"assertionConfig"`
}

var _ validation.Validatable = (*TestSuiteSpec)(nil)

func (suiteSpec *TestSuiteSpec) Validate(validator *validation.Validator) {

	validator.WithField("assertions", suiteSpec.Assertions.Validate)
}

type TestSuiteSpecAssertionConfig struct {
	// Script to evaluate functions within.
	Script string `json:"script"`

	Definitions []TestSuiteSpecAssertionDefintion `json:"definitions"`
}

func (cfg *TestSuiteSpecAssertionConfig) createEmptyRuntime() (*executor.Runtime, error) {
	return executor.NewRuntime(cfg.Script)
}

func (cfg *TestSuiteSpecAssertionConfig) CreateRuntime() (*executor.Runtime, error) {
	runtime, err := cfg.createEmptyRuntime()
	if err != nil {
		return nil, err
	}

	var errs []error
	for _, defn := range cfg.Definitions {
		if err := defn.Register(runtime); err != nil {
			errs = append(errs, fmt.Errorf("couldn't register %s: %w", defn.Name, err))
		}
	}

	return runtime, errors.Join(errs...)
}

var _ validation.Validatable = (*TestSuiteSpecAssertionConfig)(nil)

func (cfg *TestSuiteSpecAssertionConfig) Validate(validator *validation.Validator) {

	runtime, err := cfg.createEmptyRuntime()
	validator.WithField("script", func(validator *validation.Validator) {
		if err != nil {
			validator.Error("invalid script: %w", err)
		}
	})

	if err != nil {
		return
	}

	validator.WithField("defintions", func(validator *validation.Validator) {
		for idx, assertion := range cfg.Definitions {
			assertion.Validate(validator.AtIndex(idx), runtime)
		}

		validation.AssertDistinctMapping(validator, cfg.Definitions, func(d TestSuiteSpecAssertionDefintion) string {
			return d.Name
		}, []string{"name"})
	})

}

type TestSuiteSpecAssertionDefintion struct {
	Name         string          `json:"name"`
	InputSchema  json.RawMessage `json:"inputSchema"`
	FunctionName string          `json:"functionName"`
}

func (defn *TestSuiteSpecAssertionDefintion) Register(runtime *executor.Runtime) error {
	return runtime.AddAssertion(defn.Name, defn.FunctionName, defn.InputSchema)
}

func (defn *TestSuiteSpecAssertionDefintion) Validate(validator *validation.Validator, runtime *executor.Runtime) {
	validator.WithField("name", func(validator *validation.Validator) {
		validation.AssertNotBlank(validator, defn.Name)
	})

	validator.WithField("inputSchema", func(validator *validation.Validator) {
		sl := gojsonschema.NewSchemaLoader()
		err := sl.AddSchema(defn.FunctionName, gojsonschema.NewStringLoader(string(defn.InputSchema)))
		if err != nil {
			validator.Error("invalid schema: %w", err)
		}
	})

	validator.WithField("functionName", func(validator *validation.Validator) {
		validation.AssertNotBlank(validator, defn.FunctionName)
	})

	if err := defn.Register(runtime); err != nil {
		validator.Error("bad assertion: %w", err)
	}

}
