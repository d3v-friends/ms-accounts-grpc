protogen:
	go-grpc protoc --config ./protoc.yaml;
create:
	sh ./docker/create/run.sh docker.dev-friends.com ms-accounts-grpc;
update:
	sh ./docker/update/run.sh docker.dev-friends.com ms-accounts-grpc;
build:
	sh ./docker/publish/run.sh docker.dev-friends.com ms-accounts-grpc ./listener/main.go;