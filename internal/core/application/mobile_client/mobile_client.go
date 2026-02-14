package mobileclient

import "context"

type GetResponse struct {
}

type Provider interface {
	GetMinCompatibleVersion(context.Context) string
}

type MobileClientProvider struct {
	MinMobileClientVersion string
}

func NewMobileClientProvider(MinMobileClientVersion string) Provider {
	return &MobileClientProvider{
		MinMobileClientVersion: MinMobileClientVersion}
}

func (m *MobileClientProvider) GetMinCompatibleVersion(context.Context) string {
	return m.MinMobileClientVersion
}
