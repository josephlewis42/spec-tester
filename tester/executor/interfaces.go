package executor

type TestSuite interface {
	ListTests() []*TestCase

	ListImplementations() []*Implementation

	ListSpecifications() []*Specification

	TestAssertionEngine() (*Runtime, error)
}
