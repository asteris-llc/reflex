package reflex

import (
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
)

func mesosExecutorMeta() (exec *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo) {
	exec = &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("reflex"),
		Name:       proto.String("reflex scheduler"),
	}

	fwinfo = &mesos.FrameworkInfo{
		User:       proto.String(""),
		Name:       proto.String("reflex"),
		Checkpoint: proto.Bool(true),
		// TODO: FailoverTimeout: proto.Int64(???),
	}

	return
}
