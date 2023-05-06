package validation

import (
	"fmt"
	"sort"
	"strings"
)

type Level int

const (
	LevelError Level = iota
	LevelWarning
	LevelInfo
)

type Validatable interface {
	Validate(v *Validator)
}

type Validator struct {
	Parent  *Validator
	Prefix  string
	Results []Result
}

type ValidationSummary struct {
	ErrorCount   int
	WarningCount int
	InfoCount    int
}

func (vs *ValidationSummary) Update(v *Validator) {
	for _, result := range v.Results {
		switch result.Level {
		case LevelError:
			vs.ErrorCount++
		case LevelWarning:
			vs.WarningCount++
		case LevelInfo:
			vs.InfoCount++
		}
	}
}

type Result struct {
	Level   Level
	Field   string
	Message string
}

func (r *Result) String() string {
	level := "UNKNOWN"
	switch r.Level {
	case LevelError:
		level = "ERROR"
	case LevelInfo:
		level = "INFO"
	case LevelWarning:
		level = "WARNING"
	}
	return fmt.Sprintf("%s: %s: %s", level, r.Field, r.Message)
}

func (v *Validator) addResult(r Result) {
	if v.Parent != nil {
		v.Parent.addResult(r)
	} else {
		v.Results = append(v.Results, r)
	}
}

// Error adds an error to the current validation context.
func (v *Validator) Error(format string, a ...any) {
	v.addResult(Result{
		Level:   LevelError,
		Field:   v.Prefix,
		Message: fmt.Sprintf(format, a...),
	})
}

// Warning adds a warning to the current validation context.
func (v *Validator) Warning(format string, a ...any) {
	v.addResult(Result{
		Level:   LevelWarning,
		Field:   v.Prefix,
		Message: fmt.Sprintf(format, a...),
	})

}

// Info adds a message to the current validation context.
func (v *Validator) Info(format string, a ...any) {
	v.addResult(Result{
		Level:   LevelInfo,
		Field:   v.Prefix,
		Message: fmt.Sprintf(format, a...),
	})
}

func (v *Validator) Field(name string) *Validator {
	return &Validator{
		Parent: v,
		Prefix: fmt.Sprintf("%s.%s", v.Prefix, name),
	}
}

func (v *Validator) AtKey(key string) *Validator {
	return &Validator{
		Parent: v,
		Prefix: fmt.Sprintf("%s[%s]", v.Prefix, key),
	}
}

func (v *Validator) AtIndex(index int) *Validator {
	return &Validator{
		Parent: v,
		Prefix: fmt.Sprintf("%s[%d]", v.Prefix, index),
	}
}

func (v *Validator) WithField(name string, callback func(*Validator)) {
	callback(v.Field(name))
}

func AssertDistinctMapping[T any, V comparable](v *Validator, slice []T, mapper func(T) V, subFieldPaths []string) {
	conflicts := make(map[V][]int)

	for idx, entry := range slice {
		key := mapper(entry)
		if existing, ok := conflicts[key]; ok {
			v.AtIndex(idx).
				Error(
					"conflicts with entries: %v, entries must have distinct sub-fields: %q",
					existing,
					subFieldPaths,
				)
		}

		conflicts[key] = append(conflicts[key], idx)
	}
}

func AssertEqual[T comparable](v *Validator, want, got T) {
	if want != got {
		v.Error("expected field to be %q got %q", want, got)
	}
}

func AssertNotBlank(v *Validator, got string) {
	if "" == strings.TrimSpace(got) {
		v.Error("must not to be blank")
	}
}

func OneOf() *OneOfBuilder {
	return &OneOfBuilder{}
}

type OneOfBuilder struct {
	entries []oneOfEntry
}

type oneOfEntry struct {
	key         string
	defined     bool
	validatable func(v *Validator)
}

func (oob *OneOfBuilder) ValidatedField(key string, defined bool, validator func(v *Validator)) *OneOfBuilder {
	oob.entries = append(oob.entries, oneOfEntry{key, defined, validator})

	return oob
}

func (oob *OneOfBuilder) Field(key string, defined bool) *OneOfBuilder {

	oob.entries = append(oob.entries, oneOfEntry{key, defined, func(v *Validator) {}})
	return oob
}

func (oob *OneOfBuilder) Validate(validator *Validator) {

	var definedFields []string
	var allFields []string

	for _, entry := range oob.entries {
		allFields = append(allFields, entry.key)
		if entry.defined {
			entry.validatable(validator.Field(entry.key))
			definedFields = append(definedFields, entry.key)
		}
	}

	sort.Strings(definedFields)
	sort.Strings(allFields)

	switch len(definedFields) {
	case 0:
		validator.Error("requires one of the following child fields %q", allFields)
	case 1:
		// Good.
	default:
		validator.Error("only one of %q may be set", definedFields)
	}

}
