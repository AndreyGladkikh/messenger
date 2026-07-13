package uuid

import "github.com/google/uuid"

type Provider struct {

}

func (p *Provider) ID() string {
	return uuid.NewString()
}