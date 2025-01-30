package test

type TestContainer struct {
}

func (TestContainer) CreateTestContainer() (*TestContainer, error) {
	tc := TestContainer{}
	initTestContainer(&tc)
	return &tc, nil
}

func initTestContainer(tc *TestContainer) error {
	// initDependencies...
	return nil
}
