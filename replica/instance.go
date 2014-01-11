package replica

import (
	cmd "github.com/go-epaxos/epaxos/command"
)

const (
	none int8 = iota
	preaccepted
	accepted
	committed
	executed
)

type Instance struct {
	cmds   []cmd.Command
	seq    int
	deps   []InstanceIdType
	status int8
}
