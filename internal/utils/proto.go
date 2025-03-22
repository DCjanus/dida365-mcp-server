package utils

import "google.golang.org/protobuf/proto"

func ProtoClone[M proto.Message](m M) M {
	return proto.Clone(m).(M)
}
