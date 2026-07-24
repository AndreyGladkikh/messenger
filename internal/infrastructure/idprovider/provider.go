package idprovider

import "github.com/google/uuid"

type Provider struct {
}

func NewProvider() *Provider {
	return new(Provider)
}

func (p *Provider) ID() string {
	return uuid.NewString()
}
