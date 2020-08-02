package tea

type testError string

func (e testError) Error() string { return string(e) }

const (
	PlanError = testError("test plan error")
	RunError  = testError("test run error")
)
