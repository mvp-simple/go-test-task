module github.com/ruslan-onishchenko/go-test-task

go 1.22.2

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.34.2-20240717164558-a6c49f84cc0f.2
	github.com/golang/protobuf v1.5.4
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	google.golang.org/genproto v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/stretchr/testify v1.9.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240528184218-531527333157
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.5.4
	google.golang.org/protobuf => google.golang.org/protobuf v1.34.1
)
