gen:
	go-grpc protoc --config=protoc.yaml;
createMain:
	sh ./docker/builder/run.sh docker.dev-friends.com ms-accounts-grpc main
publishMain:
	sh ./docker/publisher/run.sh docker.dev-friends.com ms-accounts-grpc main 15001 ./listener/main.go
updateMain:
	sh ./docker/updater/run.sh docker.dev-friends.com ms-accounts-grpc main;
