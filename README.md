# Protobuf to Go struct literal conversion

The `proto2` package provides a single function `GoString()` that returns a Go string literal representation of a protobuf message.
This can be handy if you want to generate Go code from a protobuf message.

## Example

```go
// msg is an arbitrary protobuf message
msg := &hotstuff.Proposal{
 Block: &hotstuff.Block{
  Parent:   []byte{1, 2, 3},
  QC:       &hotstuff.QuorumCert{},
  View:     1,
  Command:  []byte{4, 5, 6},
  Proposer: 2,
 },
 AggQC: &hotstuff.AggQC{},
}
fmt.Println(proto2.GoString(msg))
```

This produces the following output:

```go
&hotstuff.Proposal{
 Block: &hotstuff.Block{
  Parent:   []byte{1, 2, 3},
  QC:       &hotstuff.QuorumCert{},
  View:     1,
  Command:  []byte{4, 5, 6},
  Proposer: 2,
 },
 AggQC: &hotstuff.AggQC{},
}
```

## Install

```bash
go get -u github.com/meling/proto2
```

## License

MIT -- see [LICENSE](LICENSE) file
