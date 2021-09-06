#!/bin/sh

set -e
test -n "$DEBUG" && set -x

server="${REDIS_HOST:-127.0.0.1:6379}"
bind="${MEMCACHE_SERVER:-0.0.0.0:11211}"

/bin/memcache-redis-adapter -server=$server -bind=$bind
