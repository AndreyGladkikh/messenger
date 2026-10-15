package entry

import (
	"fmt"
	"slices"
)

type TargetType string

func (t TargetType) String() string {
	return string(t)
}

const (
	TargetTypeUser TargetType = "user"
	TargetTypeBot TargetType = "bot"
	TargetTypeGroup TargetType = "group"
	TargetTypeChannel TargetType = "channel"
)

var targetTypes = []TargetType{
	TargetTypeUser,
	TargetTypeBot,
	TargetTypeGroup,
	TargetTypeChannel,
}

type UnsupportableTargetTypeError struct {
	Type string
}

func (e UnsupportableTargetTypeError) Error() string {
	return fmt.Sprintf("entry target type %q is not supported", e.Type)
}

func (e UnsupportableTargetTypeError) AllowedValues() []TargetType {
	return targetTypes
}

func ParseTargetType(t string) (TargetType, error) {
	trgt := TargetType(t)
	if slices.Contains(targetTypes, trgt) {
		return trgt, nil
	}
	return "", UnsupportableTargetTypeError{t}
}
