//nolint:forcetypeassert, gocritic
package ibft

import (
	"context"

	"github.com/madz-lab/go-ibft/message/types"
)

type (
	Signer     = types.Signer
	SigRecover = types.SigRecover
	Feed       = types.MsgFeed

	Transport interface {
		Multicast(types.Msg)
	}

	Quorum interface {
		HasQuorum(uint64, []types.Msg) bool
	}

	Keccak interface {
		Hash([]byte) []byte
	}

	Verifier interface {
		IsValidator(id []byte, sequence uint64) bool
		IsValidBlock(block []byte, sequence uint64) bool
		IsProposer(id []byte, sequence uint64, round uint64) bool
	}

	Validator interface {
		Signer

		ID() []byte
		BuildBlock(uint64) []byte
	}
)

// Context keys
type ctxKey string

const (
	transport  ctxKey = "transport"
	feed       ctxKey = "feed"
	quorum     ctxKey = "quorum"
	keccak     ctxKey = "keccak"
	sigRecover ctxKey = "sig_recover"
)

type Context struct {
	context.Context
}

func NewIBFTContext(ctx context.Context) Context {
	return Context{ctx}
}

func (c Context) WithCancel() (Context, func()) {
	subCtx, cancelFn := context.WithCancel(c)

	return Context{subCtx}, cancelFn
}

func (c Context) WithTransport(t Transport) Context {
	return Context{context.WithValue(c, transport, t)}
}

func (c Context) Transport() Transport {
	return c.Value(transport).(Transport)
}

func (c Context) WithFeed(f Feed) Context {
	return Context{context.WithValue(c, feed, f)}
}

func (c Context) Feed() Feed {
	return c.Value(feed).(Feed)
}

func (c Context) WithQuorum(q Quorum) Context {
	return Context{context.WithValue(c, quorum, q)}
}

func (c Context) Quorum() Quorum {
	return c.Value(quorum).(Quorum)
}

func (c Context) WithKeccak(k Keccak) Context {
	return Context{context.WithValue(c, keccak, k)}
}

func (c Context) Keccak() Keccak {
	return c.Value(keccak).(Keccak)
}

func (c Context) WithSigRecover(s SigRecover) Context {
	return Context{context.WithValue(c, sigRecover, s)}
}

func (c Context) SigRecover() SigRecover {
	return c.Value(sigRecover).(SigRecover)
}
