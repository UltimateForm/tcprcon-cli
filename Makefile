RCON_PASSWORD=localpassword
RCON_PORT=7778


wire-podman-vlms:
	podman unshare chown 1000:1000 .serverdata/rust/common.cfg
	podman unshare chown 1000:1000 .serverdata/mh/common.cfg
	podman unshare chown 1000:1000 .serverdata/mh/Game.ini

lift-mh-server:
	docker compose run mh-server

lift-rust-server:
	docker compose run rust-server

build:
	go build -o .out/tcprcon

run: build
	.out/tcprcon -address=localhost -port=${RCON_PORT} -pw=${RCON_PASSWORD}

test:
	go test ./...

build-docker:
	docker build . -t tcprcon

run-docker: build-docker
	docker run -it tcprcon:latest -address=host.docker.internal -port=${RCON_PORT} -pw=${RCON_PASSWORD}
