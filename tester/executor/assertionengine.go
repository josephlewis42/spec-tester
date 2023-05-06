package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/xeipuuv/gojsonschema"
)

// Checks that a script parses as valid ES5.
func CheckValid(program string) error {
	_, err := NewRuntime(program)
	return err
}

type Runtime struct {
	prg *goja.Program

	funcLookup map[string]string

	schemaDefinitions map[string]*gojsonschema.Schema
}

type Evaluation struct {
}

func NewRuntime(program string) (*Runtime, error) {
	prg, err := goja.Compile("runtime", program, true)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		prg:               prg,
		funcLookup:        make(map[string]string),
		schemaDefinitions: make(map[string]*gojsonschema.Schema),
	}, nil
}

func (runtime *Runtime) AddAssertion(name, functionName string, schema json.RawMessage) error {
	vm := goja.New()
	_, err := vm.RunProgram(runtime.prg)
	if err != nil {
		return fmt.Errorf("error in assertion script: %e", err)
	}
	if _, ok := goja.AssertFunction(vm.Get(functionName)); !ok {
		return fmt.Errorf("%q isn't a function", functionName)
	}

	parsedSchemaJson := make(map[string]any)
	if err := json.Unmarshal(schema, &parsedSchemaJson); err != nil {
		return fmt.Errorf("schema is invalid JSON: %e", err)
	}
	parsedSchemaJson["$id"] = fmt.Sprintf("spectest://%s", name)
	parsedSchemaJson["$schema"] = "http://json-schema.org/draft-07/schema#"

	jsonSchema, err := gojsonschema.NewSchemaLoader().Compile(gojsonschema.NewGoLoader(parsedSchemaJson))
	if err != nil {
		return fmt.Errorf("schema is invalid: %e", err)
	}

	runtime.funcLookup[name] = functionName
	runtime.schemaDefinitions[name] = jsonSchema

	return nil
}

type testSettings struct {
	functionName  string
	parsedOptions any
}

func (runtime *Runtime) CheckTestInput(evalTest *EvalTest) (*testSettings, error) {
	functionName, foundFunctionName := runtime.funcLookup[evalTest.ExpectationType]
	if !foundFunctionName {
		return nil, fmt.Errorf("couldn't find function for expectation type %q", evalTest.ExpectationType)
	}

	schema, foundSchema := runtime.schemaDefinitions[evalTest.ExpectationType]
	if !foundSchema {
		return nil, fmt.Errorf("couldn't find schema for expectation type %q", evalTest.ExpectationType)
	}

	var parsedOptionsJson interface{}
	if err := json.Unmarshal([]byte(evalTest.ExpectationOptionsJson), &parsedOptionsJson); err != nil {
		return nil, fmt.Errorf("options are invalid JSON: %e", err)
	}

	result, err := schema.Validate(gojsonschema.NewGoLoader(parsedOptionsJson))
	if err != nil {
		return nil, fmt.Errorf("couldn't load expectation options: %e", err)
	}

	if result.Valid() {
		return &testSettings{
			functionName:  functionName,
			parsedOptions: parsedOptionsJson,
		}, nil
	}

	return nil, fmt.Errorf("validation errors: %q", result.Errors())
}

func (runtime *Runtime) EvaluateTestResult(ctx context.Context, testCase *TestCase, output *ProcessOutput) (*TestResult, error) {

	switch testType := testCase.TestType.(type) {
	default:
		return nil, fmt.Errorf("can't evaluate tests with type: %t", testType)

	case *TestCase_Eval:
		// Lookup type and validation

		vm := goja.New()
		_, err := vm.RunProgram(runtime.prg)
		if err != nil {
			return nil, fmt.Errorf("error in assertion script: %w", err)
		}

		testSettings, err := runtime.CheckTestInput(testType.Eval)
		if err != nil {
			return nil, fmt.Errorf("error in test case: %w", err)
		}

		assertionFunction, ok := goja.AssertFunction(vm.Get(testSettings.functionName))
		if !ok {
			return nil, fmt.Errorf(
				"function %q used by expectation type %q not defined in assertion script",
				testSettings.functionName,
				testType.Eval.ExpectationType,
			)
		}

		res, err := assertionFunction(goja.Undefined(), vm.ToValue(
			map[string]any{
				"metadata": map[string]any{
					"uid":    testCase.Metadata.Uid,
					"labels": testCase.Metadata.Labels,
				},
				"input":  testType.Eval.Input,
				"config": testSettings.parsedOptions,
				"output": map[string]any{
					"stdout":   output.Stdout,
					"stderr":   output.Stderr,
					"exitCode": output.ExitCode,
				},
			}))
		if err != nil {
			return nil, fmt.Errorf("error executing evaluation function: %w", err)
		}

		bytes, err := json.Marshal(res)
		if err != nil {
			return nil, fmt.Errorf("couldn't convert assertion response to JSON: %w", err)
		}

		out := new(TestResult)
		if err := json.Unmarshal(bytes, out); err != nil {
			return nil, fmt.Errorf("invalid response format: %w", err)
		}

		return out, nil
	}
}
