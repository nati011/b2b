package container

type TestContainer interface {
	NewTestContainer() (*TestContainer, error)
}

type MasterTestContainer struct {
	TestContainer []*any
}

func (mt MasterTestContainer) NewMasterTestContainer() (*MasterTestContainer, error) {
	tc, err := mt.init()
	if err != nil {
		panic("failed to create MasterTestContainer")
	}

	return tc, nil
}

func (mt *MasterTestContainer) RegisterTestContainer(tc *any) {
	mt.TestContainer = append(mt.TestContainer, tc)
}

func (mt MasterTestContainer) init() (*MasterTestContainer, error) {
	tc := MasterTestContainer{}
	return &tc, nil
}
