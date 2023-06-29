package proto2_test

import (
	"go/format"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/meling/proto2"
	"github.com/meling/proto2/internal/testprotos/hotstuff"
	"google.golang.org/protobuf/proto"
)

func TestGoString(t *testing.T) {
	testCases := []struct {
		input proto.Message
		want  string
	}{
		{
			input: &hotstuff.Proposal{
				Block: &hotstuff.Block{
					Parent:   []byte{1, 2, 3},
					QC:       &hotstuff.QuorumCert{},
					View:     1,
					Command:  []byte{4, 5, 6},
					Proposer: 2,
				},
				AggQC: &hotstuff.AggQC{},
			},
			want: `&hotstuff.Proposal{
				Block: &hotstuff.Block{
					Parent:   []byte{1, 2, 3},
					QC:       &hotstuff.QuorumCert{},
					View:     1,
					Command:  []byte{4, 5, 6},
					Proposer: 2,
				},
				AggQC: &hotstuff.AggQC{},
			}`,
		},
		{
			input: &hotstuff.BlockHash{
				Hash: []byte{1, 2, 3, 4, 5},
			},
			want: `&hotstuff.BlockHash{
				Hash: []byte{1, 2, 3, 4, 5},
			}`,
		},
		{
			input: &hotstuff.AggQC{
				QCs: map[uint32]*hotstuff.QuorumCert{
					1: {},
					2: {},
				},
				Sig:  &hotstuff.QuorumSignature{},
				View: 3,
			},
			want: `&hotstuff.AggQC{
				QCs: map[uint32]*hotstuff.QuorumCert{
					1: {},
					2: {},
				},
				Sig:  &hotstuff.QuorumSignature{},
				View: 3,
			}`,
		},
		{
			input: &hotstuff.ECDSAMultiSignature{
				Sigs: []*hotstuff.ECDSASignature{
					{
						Signer: 1,
						R:      []byte{1, 2, 3},
						S:      []byte{4, 5, 6},
					},
					{
						Signer: 2,
						R:      []byte{7, 8, 9},
						S:      []byte{10, 11, 12},
					},
				},
			},
			want: `&hotstuff.ECDSAMultiSignature{
				Sigs: []*hotstuff.ECDSASignature{
					{
						Signer: 1,
						R:      []byte{1, 2, 3},
						S:      []byte{4, 5, 6},
					},
					{
						Signer: 2,
						R:      []byte{7, 8, 9},
						S:      []byte{10, 11, 12},
					},
				},
			}`,
		},
		{
			input: &hotstuff.QuorumSignature{
				Sig: &hotstuff.QuorumSignature_ECDSASigs{
					ECDSASigs: &hotstuff.ECDSAMultiSignature{
						Sigs: []*hotstuff.ECDSASignature{
							{
								Signer: 1,
								R:      []byte{1, 2, 3},
								S:      []byte{4, 5, 6},
							},
						},
					},
				},
			},
			want: `&hotstuff.QuorumSignature{
				Sig: &hotstuff.QuorumSignature_ECDSASigs{
					ECDSASigs: &hotstuff.ECDSAMultiSignature{
						Sigs: []*hotstuff.ECDSASignature{
							{
								Signer: 1,
								R:      []byte{1, 2, 3},
								S:      []byte{4, 5, 6},
							},
						},
					},
				},
			}`,
		},
	}

	for _, testCase := range testCases {
		got := proto2.GoString(testCase.input)
		formatted, err := format.Source([]byte(testCase.want))
		if err != nil {
			t.Fatal(err)
		}
		want := string(formatted) + "\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("proto2.GoString() mismatch (-want, +got):\n%s", diff)
		}
	}
}
