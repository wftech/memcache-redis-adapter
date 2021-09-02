package proxy

import (
	"github.com/gomodule/redigo/redis"
	"github.com/wftech/memcache-redis-adapter/protocol"
)

type RedisProxy struct {
	conn redis.Conn
}

func NewRedisProxy(conn redis.Conn) *RedisProxy {
	r := new(RedisProxy)
	r.conn = conn
	return r
}

func serverError(err error) protocol.McResponse {
	return protocol.McResponse{Response: "SERVER_ERROR " + err.Error()}
}

func serverErrorText(err error, text string) protocol.McResponse {
	return protocol.McResponse{Response: "SERVER_ERROR " + err.Error() + " (" + text + ")"}
}

// process a request and generate a response
func (p *RedisProxy) Process(req *protocol.McRequest) protocol.McResponse {

	switch req.Command {
	case "get":
		res := protocol.McResponse{}
		for i := range req.Keys {

			r, err := redis.Values(p.conn.Do("MGET", req.Keys[i], req.Keys[i]+"_mcflags"))
			if err != nil {
				// hmm, barf errors, or just ignore?
				return serverError(err)
			}
			if r[0] != nil {
				data, err := redis.Bytes(r[0], nil)
				if err != nil {
					return serverErrorText(err, "data")
				}
				flags := "0"
				if r[1] != nil {
					flags, err = redis.String(r[1], nil)
					if err != nil {
						return serverErrorText(err, "flags")
					}
				}
				// todo, both can return error
				res.Values = append(res.Values, protocol.McValue{req.Keys[i], flags, data})
			}
		}
		res.Response = "END"
		return res

	// TODO - check `add`
	case "set", "add":
		r, err := redis.String(p.conn.Do("MSET", req.Key, req.Data, req.Key+"_mcflags", req.Flags))
		if err != nil || r != "OK" {
			return serverError(err)
		}

		if req.Exptime != 0 {
			_, err = p.conn.Do("EXPIREAT", req.Key, req.Exptime)
			if err != nil {
				return serverError(err)
			}
		}

		return protocol.McResponse{Response: "STORED"}

	case "delete":
		r, err := redis.Int(p.conn.Do("DEL", toInterface(req.Keys)...))
		if err != nil {
			return serverError(err)
		}
		if r > 0 {
			return protocol.McResponse{Response: "DELETED"}
		}
		return protocol.McResponse{Response: "NOT_FOUND"}

		// todo "touch"...
	}

	return protocol.McResponse{Response: "ERROR"}

}

func toInterface(s []string) []interface{} {

	ret := make([]interface{}, len(s))
	for i, v := range s {
		ret[i] = interface{}(v)
	}
	return ret
}
