package testbed

import (
	"testing"
)

type TestCase struct {
	t *testing.T

	// Flex process
	flex FlexRunner

	// Execution validator
	validator ExecutionValidator
}

func NewTestCase(t *testing.T, f FlexRunner, v ExecutionValidator) *TestCase {
	return &TestCase{
		t:         t,
		flex:      f,
		validator: v,
	}
}

func (tc *TestCase) RunTest() {
	err := tc.flex.Run()
	if err != nil {
		tc.t.Error(err)
		return
	}

	stdout, stderr, err := tc.flex.Results()
	if err != nil {
		tc.t.Error(err)
		return
	}

	err = tc.validator.Validate(tc.t, stdout, stderr)
	if err != nil {
		tc.t.Error(err)
		return
	}
}
