package command

type Command interface {
	IsCommand()
	Name() string
}