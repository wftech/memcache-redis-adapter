package proxy

import (
	"github.com/wftech/memcache-redis-adapter/protocol"
)

type ProtocolProxy interface {
	Process(*protocol.McRequest) protocol.McResponse
}
