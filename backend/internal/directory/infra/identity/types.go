package identity

import (
	"fmt"
	"messenger/messenger/internal/directory/domain/entry"
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

var typeMap = map[identityType]entry.TargetType{
	identityTypeUser: entry.TargetTypeUser,
	identityTypeBot:  entry.TargetTypeBot,
}

func mapIdentityTypeToTarget(t string) (entry.TargetType, error) {
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
