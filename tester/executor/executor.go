package executor

import (
	"context"
	"fmt"

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

	implementations := opts.ImplementationFilter.Apply(suite.ListImplementations())
	allSpecifications := opts.SpecificationFilter.Apply(suite.ListSpecifications())
	allTests := opts.TestFilter.Apply(suite.ListTests())

	for _, impl := range implementations {
		for _, variant := range impl.Variants {
			specifications := opts.SpecificationFilter.
				WithUid(variant.GetSpecificationUids()...).
				Apply(allSpecifications)

			testsToRun := NewFilter[*TestCase]()
			for _, spec := range specifications {
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
					fmt.Printf(
						"%s/%s/%s: Running test: %q\n",
						impl.GetMetadata().GetUid(),
						variant.GetMetadata().GetUid(),
						spec.GetMetadata().GetUid(),
						test.GetMetadata().GetUid(),
					)
				}
			}
		}
	}

	return nil
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
