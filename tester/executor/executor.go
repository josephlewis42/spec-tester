package executor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/josephlewis42/scheme-compliance/internal/specctx"
	"golang.org/x/exp/slog"
	"k8s.io/apimachinery/pkg/labels"
)

type ExecutionOptions struct {
	// Filter for which specifications to run.
	SpecificationFilter Filter[*Specification]

	// Filter for which implementations to run.
	ImplementationFilter Filter[*Implementation]

	// Filter for which tests to run.
	TestFilter Filter[*TestCase]

	// Only execute missing items.
	OnlyMissing bool
}

func Execute(ctx context.Context, suite TestSuite, opts ExecutionOptions) error {
	// For spec, implementation, suite
	log := specctx.GetLogger(ctx)

	runtime, err := suite.TestAssertionEngine()
	if err != nil {
		return fmt.Errorf("couldn't create execution engine: %e", err)
	}

	implementations := opts.ImplementationFilter.Apply(suite.ListImplementations())
	allSpecifications := opts.SpecificationFilter.Apply(suite.ListSpecifications())
	allTests := opts.TestFilter.Apply(suite.ListTests())

	for _, impl := range implementations {
		log := log.With(slog.String("implementation", impl.GetMetadata().GetUid()))
		for _, variant := range impl.Variants {
			log := log.With(slog.String("variant", variant.GetMetadata().GetUid()))

			specifications := opts.SpecificationFilter.
				WithUid(variant.GetSpecificationUids()...).
				Apply(allSpecifications)

			testsToRun := NewFilter[*TestCase]()
			for _, spec := range specifications {
				log := log.With(slog.String("specification", spec.GetMetadata().GetUid()))

				testFilters, err := getTestFilters(spec.Sections)
				if err != nil {
					return err
				}

				for _, filter := range testFilters {
					filter.ForEach(allTests, func(matching *TestCase) {
						testsToRun = testsToRun.WithUid(matching.GetMetadata().GetUid())
					})
				}

				for _, test := range testsToRun.Apply(allTests) {
					log := log.With(slog.String("test", test.GetMetadata().GetUid()))

					switch testType := test.TestType.(type) {
					case *TestCase_Skip:
						log.Debug("Test skipped", "reason", testType.Skip.Message)
						continue
					case *TestCase_Invalid:
						log.Error("Invalid test", "reason", testType.Invalid.Message)
						return fmt.Errorf("invalid test, message: %s", testType.Invalid.Message)

					case *TestCase_Eval:
						log.Debug("Running test")
						out, err := executeTest(specctx.WithLogger(ctx, log), test, variant)
						if err != nil {
							return err
						}

						result, err := runtime.EvaluateTestResult(ctx, test, out)
						if err != nil {
							return fmt.Errorf("couldn't evaluate result: %w", err)
						}

						log.Info("Completed evaluation", "result", result)

					default:
						return fmt.Errorf("invalid test type %T", test.TestType)
					}
				}
			}
		}
	}

	return nil
}

func executeTest(ctx context.Context, testCase *TestCase, variant *ImplementationVariant) (*ProcessOutput, error) {
	// For spec, implementation, suite
	log := specctx.GetLogger(ctx)

	if _, ok := testCase.TestType.(*TestCase_Eval); !ok {
		return nil, fmt.Errorf("can't execute test with type %T", testCase.TestType)
	}

	switch runtime := variant.Runtime.(type) {
	case *ImplementationVariant_Local:

		tmp, err := os.CreateTemp("", "spec-test")
		if err != nil {
			return nil, err
		}
		defer os.Remove(tmp.Name())

		program := testCase.GetEval().GetInput()

		if _, err := tmp.Write([]byte(program)); err != nil {
			return nil, fmt.Errorf("couldn't write test to file %q: %w", tmp.Name(), err)
		}

		cmd := formatEvalCommand(variant.GetTestCommand(), program, tmp.Name())
		if len(cmd) == 0 {
			return nil, errors.New("can't execute an empty command")
		}

		stdoutBuffer := &bytes.Buffer{}
		stderrBuffer := &bytes.Buffer{}

		log.Info("Running command", slog.Any("command", cmd))
		executingCommand := exec.Command(cmd[0], cmd[1:]...)
		executingCommand.Stdin = nil
		executingCommand.Stdout = stdoutBuffer
		executingCommand.Stderr = stderrBuffer
		err = executingCommand.Run()
		switch {
		case err == nil:
			return &ProcessOutput{
				Stdout:   stdoutBuffer.String(),
				Stderr:   stderrBuffer.String(),
				ExitCode: 0,
			}, nil
		default:
			if exitErr, ok := err.(*exec.ExitError); ok {
				return &ProcessOutput{
					Stdout:   stdoutBuffer.String(),
					Stderr:   stderrBuffer.String(),
					ExitCode: int64(exitErr.ExitCode()),
				}, nil

			}
			return nil, fmt.Errorf("couldn't run command: %q: %w", cmd, err)
		}
	default:
		return nil, fmt.Errorf("Unknown runtime type: %t", runtime)
	}
}

func formatEvalCommand(cmd []string, program, programPath string) []string {
	replacer := strings.NewReplacer("$(PROGRAM)", program, "$(PROGRAM_PATH)", programPath)

	var replaced []string
	for _, part := range cmd {
		replaced = append(replaced, replacer.Replace(part))
	}

	return replaced
}

func getTestFilters(sections []*SpecificationSection) ([]Filter[*TestCase], error) {
	var filters []Filter[*TestCase]

	for _, section := range sections {
		switch content := section.Content.(type) {
		case *SpecificationSection_SectionSummary:

			subsectionFilters, err := getTestFilters(content.SectionSummary.GetSubsections())
			if err != nil {
				return nil, err
			}

			filters = append(filters, subsectionFilters...)

		case *SpecificationSection_TestSummary:
			testSelector, err := labels.Parse(content.TestSummary.TestSelector)
			if err != nil {
				return nil, err
			}

			filters = append(filters, NewFilter[*TestCase]().WithSelector(testSelector))
		}
	}

	return filters, nil
}
