

cd src
protoc -I ../../../server/common-plugin/proto/ ../../../server/common-plugin/proto/kaetzchen.proto --rust-grpc_out=.
protoc -I ../../../server/common-plugin/proto/ ../../../server/common-plugin/proto/kaetzchen.proto --rust_out=.

