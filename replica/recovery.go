package replica

import (
	"fmt"
)

var _ = fmt.Printf

func (r *Replica) sendPrepare(L int, instanceId InstanceId, messageChan chan Message) {
	if r.InstanceMatrix[L][instanceId] == nil {
		// TODO: we need to prepare an instance that doesn't exist.
		r.InstanceMatrix[L][instanceId] = &Instance{
			// TODO:
			// Assumed no-op to be nil here.
			// we need to do more since state machine needs to know how to interpret it.
			cmds:   nil,
			deps:   make([]InstanceId, r.Size), // TODO: makeInitialDeps
			status: -1,                         // 'none' might be a conflicting name. We currenctly pick '-1' for it
			ballot: r.makeInitialBallot(),
			info:   NewInstanceInfo(),
		}
	}

	inst := r.InstanceMatrix[L][instanceId]

	inst.ballot.incNumber()
	inst.ballot.setReplicaId(r.Id)

	prepare := &Prepare{
		ballot:     inst.ballot,
		replicaId:  L,
		instanceId: instanceId,
	}

	go func() {
		for i := 0; i < r.Size-1; i++ {
			messageChan <- prepare
		}
	}()
}

func (r *Replica) recvPrepare(pp *Prepare, messageChan chan Message) {
	inst := r.InstanceMatrix[pp.replicaId][pp.instanceId]
	if inst == nil {
		// reply PrepareReply with no-op and invalid status
		pr := &PrepareReply{
			ok:         true,
			ballot:     pp.ballot,
			status:     -1,                         // TODO: hardcode, not a best approach
			deps:       make([]InstanceId, r.Size), // TODO: makeInitialDeps
			replicaId:  pp.replicaId,
			instanceId: pp.instanceId,
		}
		r.sendPrepareReply(pr, messageChan)
		return
	}

	// we have some info about the instance
	pr := &PrepareReply{
		status:     inst.status,
		replicaId:  pp.replicaId,
		instanceId: pp.instanceId,
	}

	if inst.isAfterStatus(preAccepted) {
		pr.cmds = inst.cmds
		pr.deps = inst.deps
	}

	if inst.status == preAccepted {
		if pp.ballot.Compare(inst.ballot) > 0 {
			pr.cmds = inst.cmds
			pr.deps = inst.deps
			pr.ok = true
			inst.ballot = pp.ballot
		} else {
			pr.ok = false
		}
	}

	pr.ballot = inst.ballot
	r.sendPrepareReply(pr, messageChan)
}

func (r *Replica) sendPrepareReply(pr *PrepareReply, messageChan chan Message) {
	go func() {
		messageChan <- pr
	}()
}

func (r *Replica) recvPrepareReply(p *PrepareReply, m chan Message) {
	inst := r.InstanceMatrix[p.replicaId][p.instanceId]

	if inst == nil {
		// it shouldn't happen
		return
	}

	if inst.isAfterStatus(preAccepted) {
		return
	}

	if p.status > preAccepted {
		inst.cmds = p.cmds
		inst.deps = p.deps
		inst.status = p.status
		inst.ballot = p.ballot

		switch p.status {
		case committed:
			// send commit
		case accepted:
			// send accept
		}
		return
	}

	if p.status == preAccepted {

	} else {
		// no op
	}
}
