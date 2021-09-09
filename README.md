# Memcache Redis Adapter

## Requirements

- Docker
- GNU Make

## How to compile

- `make` creates statically linked binary
- `make runshell` opens shell inside Docker container (`vim` setup for hacking included)
- `make image` creates Docker image

## Currently supported and tested commands

- `set`
- `add`
- `get`
- `delete`
- `touch`
- `incr`
- `decr`
- `replace`

## Credits

- [mrproxy](https://github.com/zobo/mrproxy)
