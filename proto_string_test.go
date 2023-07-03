package proto2_test

import (
	"go/format"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/meling/proto2"
	"github.com/meling/proto2/internal/testprotos/hotstuff"
	"google.golang.org/protobuf/proto"
)

func TestGoStruct(t *testing.T) {
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
		{
			input: &hotstuff.Proposal{
				Block: &hotstuff.Block{
					Parent:   []byte{136, 111, 248, 42, 223, 172},
					Command:  []byte{136, 111, 248, 42, 223, 172},
					View:     3877690086,
					Proposer: 21312,
					QC: &hotstuff.QuorumCert{
						Sig: &hotstuff.QuorumSignature{
							Sig: nil,
						},
						Hash: []byte{136, 111, 248, 42, 223, 172},
					},
				},
				AggQC: &hotstuff.AggQC{
					QCs: map[uint32]*hotstuff.QuorumCert{
						3877690086: {
							Sig: &hotstuff.QuorumSignature{
								Sig: nil,
							},
							Hash: []byte{136, 111, 248, 42, 223, 172},
							View: 3877690086,
						},
					},
				},
			},
			want: `&hotstuff.Proposal{
				Block: &hotstuff.Block{
					Parent:   []byte{136, 111, 248, 42, 223, 172},
					QC: &hotstuff.QuorumCert{
						Sig: &hotstuff.QuorumSignature{},
						Hash: []byte{136, 111, 248, 42, 223, 172},
					},
					View:     3877690086,
					Command:  []byte{136, 111, 248, 42, 223, 172},
					Proposer: 21312,
				},
				AggQC: &hotstuff.AggQC{
					QCs: map[uint32]*hotstuff.QuorumCert{
						3877690086: {
							Sig: &hotstuff.QuorumSignature{},
							View: 3877690086,
							Hash: []byte{136, 111, 248, 42, 223, 172},
						},
					},
				},
			}`,
		},
		{
			input: &hotstuff.MsgInfo{
				Type:   hotstuff.MsgType_NEWVIEW,
				Height: 1,
				Round:  2,
				Step:   4,
			},
			want: `&hotstuff.MsgInfo{
				Type:   hotstuff.MsgType_NEWVIEW,
				Height: 1,
				Round:  2,
				Step:   4,
			}`,
		},
		{
			input: &hotstuff.MsgInfo{
				Type:   hotstuff.MsgType_PROPOSAL, // zero-value will not be printed
				Height: 1,
				Round:  2,
				Step:   4,
			},
			want: `&hotstuff.MsgInfo{
				Height: 1,
				Round:  2,
				Step:   4,
			}`,
		},
	}

	for _, testCase := range testCases {
		got := proto2.GoStruct(testCase.input)
		formatted, err := format.Source([]byte(testCase.want))
		if err != nil {
			t.Fatal(err)
		}
		want := string(formatted) + "\n"
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("proto2.GoStruct() mismatch (-want, +got):\n%s", diff)
		}
	}
}
