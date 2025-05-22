package provider

import (
	port "b2b.nati011.github.com/internal/port/application/auth/provider"
)

type MapClaims map[string]interface{}

func (m MapClaims) GetEmail() (string, error) {
	return m.parseString("email")
}

func (m MapClaims) parseString(key string) (string, error) {
	var (
		ok  bool
		raw interface{}
		iss string
	)
	raw, ok = m[key]
	if !ok {
		return "", nil
	}

	iss, ok = raw.(string)
	if !ok {
		return "", port.ErrSysUnknown
	}

	return iss, nil
}
