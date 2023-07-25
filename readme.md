# grpc 패키지 설치하기

~~~bash
go get google.golang.org/protobuf/cmd/protoc-gen-go
go run google.golang.org/protobuf/cmd/protoc-gen-go

~~~

# cli 설치하기

~~~bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install github.com/d3v-friends/go-grpc/go-grpc

~~~

# proto 생성

~~~bash
go-grpc protoc --config=protoc.yaml

~~~

# 환경변수

* 기본 배포 포트 15001

TZ=
PORT=
MG_HOST=
MG_USERNAME=
MG_PASSWORD=
MG_DATABASE=
