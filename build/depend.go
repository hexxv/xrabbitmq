package build

import (
	"xrabbitmq/pkg/external"
)

type Required func(*depend)

type depend struct {
	conn *external.XConnection
}

func DependConn(conn *external.XConnection) Required {
	return func(b *depend) {
		b.conn = conn
	}
}
