package executor

type TestSuite interface {
	HydrateTests() []TestCase
}
