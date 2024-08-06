# download proto google
proto-google:
	curl https://raw.githubusercontent.com/googleapis/googleapis/974ad5bdfc9ba768db16b3eda2850aadd8c10a2c/google/api/annotations.proto --create-dirs -o proto/google/api/annotations.proto
	curl https://raw.githubusercontent.com/googleapis/googleapis/974ad5bdfc9ba768db16b3eda2850aadd8c10a2c/google/api/http.proto --create-dirs -o proto/google/api/http.proto

# download proto validate
proto-validate:
	curl https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/validate.proto --create-dirs -o proto/buf/validate/validate.proto
	curl https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/expression.proto --create-dirs -o proto/buf/validate/expression.proto
	curl https://raw.githubusercontent.com/bufbuild/protovalidate/main/proto/protovalidate/buf/validate/priv/private.proto --create-dirs -o proto/buf/validate/priv/private.proto

# download proto openapiv2
proto-openapiv2:
	curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/annotations.proto --create-dirs -o proto/protoc-gen-openapiv2/options/annotations.proto
	curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/main/protoc-gen-openapiv2/options/openapiv2.proto --create-dirs -o proto/protoc-gen-openapiv2/options/openapiv2.proto

PACKAGE=github.com/ruslan-onishchenko/go-test-task

# generate grpc rest and openapi
generate:
	protoc -I proto \
		--go_opt=module=$(PACKAGE) --go_out=. \
		--go-grpc_opt=module=$(PACKAGE) --go-grpc_out=. \
		--grpc-gateway_opt=module=$(PACKAGE) --grpc-gateway_out=. \
		./proto/queue_service.proto
#		--openapiv2_out=allow_merge=true:./pkg/microservice/v1 \

install:
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc