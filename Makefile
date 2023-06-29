proto_include := internal/testprotos
proto_src :=	internal/testprotos/hotstuff/hotstuff.proto \

proto_go := $(proto_src:%.proto=%.pb.go)

.PHONY: protos

protos: $(proto_go)

%.pb.go %_gorums.pb.go : %.proto
	@protoc -I=.:$(proto_include) \
		--go_out=paths=source_relative:. \
		$<
