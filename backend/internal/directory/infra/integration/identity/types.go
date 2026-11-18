package identity

import (
	"fmt"
	"messenger/messenger/internal/directory/domain/model"
)

type identityType string

func (i identityType) String() string {
	return string(i)
}

const (
	identityTypeUser identityType = "user"
	identityTypeBot  identityType = "bot"
)

func parseIdentityType(t string) (identityType, error) {
	typ := identityType(t)
	switch typ {
	case identityTypeUser,
		identityTypeBot:
		return typ, nil
	default:
		return "", fmt.Errorf("unknown identity type %q", t)
	}
}

var typeMap = map[identityType]model.TargetType{
	identityTypeUser: model.TargetTypeUser,
	identityTypeBot:  model.TargetTypeBot,
}

func mapIdentityTypeToTarget(t string) (model.TargetType, error) {
	idT, err := parseIdentityType(t)
	if err != nil {
		return "", err
	}

	trgtT, ok := typeMap[idT]
	if !ok {
		return "", fmt.Errorf("no target type corresponding to identity type %q", idT)
	}

	return trgtT, nil
}
