package command

type Command []byte

type StateMachine interface {
	// Execute a group of commands, return the result
	// in an array, return any error that occurs
	Execute(c []Command) ([]interface{}, error)
	// Test if there exists any conflicts in two group of commands
	HaveConflicts(c1 []Command, c2 []Command) bool
}
