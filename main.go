package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/wftech/memcache-redis-adapter/protocol"
	"github.com/wftech/memcache-redis-adapter/proxy"
	"github.com/wftech/memcache-redis-adapter/stats"
)

var redisServer = flag.String("server", "127.0.0.1:6379", "Redis server to connect to")
var listenAddress = flag.String("bind", "0.0.0.0:11211", "Bind address and port")

func initRedisPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			d, _ := time.ParseDuration("1s")
			c, err := redis.DialTimeout("tcp", *redisServer, d, d, d)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	
	defer c.Close()

	// take it per need
	pool := initRedisPool()
	conn := pool.Get()
	defer conn.Close()

	// process
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	redisProxy := proxy.NewRedisProxy(conn)
	proxy := stats.NewStatsProxy(redisProxy)

	for {
		req, err := protocol.ReadRequest(br)
		if perr, ok := err.(protocol.ProtocolError); ok {
			log.Printf("%v ReadRequest protocol err: %v", c, err)
			bw.WriteString("CLIENT_ERROR " + perr.Error() + "\r\n")
			bw.Flush()
			continue
		} else if err != nil {
			log.Printf("%v ReadRequest err: %v", c, err)
			return
		}
		log.Printf("%v Req: %+v\n", c, req)

		switch req.Command {
		case "quit":
			return
		case "version":
			res := protocol.McResponse{Response: "VERSION foobar"}
			bw.WriteString(res.Protocol())
			bw.Flush()
		default:
			res := proxy.Process(req)
			if !req.Noreply {
				log.Printf("%v Res: %+v\n", c, res)
				bw.WriteString(res.Protocol())
				bw.Flush()
			}
		}
	}
}

func main() {
	flag.Parse()

	l, err := net.Listen("tcp4", *listenAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}
