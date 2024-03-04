package sequencer

import (
	"bytes"

	"github.com/madz-lab/go-ibft"
	"github.com/madz-lab/go-ibft/message/types"
)

func (s *Sequencer) sendMsgPrepare(ctx Context) {
	msg := &types.MsgPrepare{
		From:      s.ID(),
		View:      s.state.View(),
		BlockHash: s.state.AcceptedBlockHash(),
	}

	msg.Signature = s.Sign(ctx.Keccak().Hash(msg.Payload()))

	ctx.MessageTransport().Prepare.Multicast(msg)
}

func (s *Sequencer) awaitPrepare(ctx Context) error {
	messages, err := s.awaitQuorumPrepares(ctx)
	if err != nil {
		return err
	}

	s.state.PrepareCertificate(messages)

	return nil
}

func (s *Sequencer) awaitQuorumPrepares(ctx Context) ([]*types.MsgPrepare, error) {
	sub, cancelSub := ctx.MessageFeed().PrepareMessages(s.state.view, false)
	defer cancelSub()

	isValidMsg := func(msg *types.MsgPrepare) bool {
		if !s.IsValidSignature(msg.GetSender(), ctx.Keccak().Hash(msg.Payload()), msg.GetSignature()) {
			return false
		}

		return s.isValidMsgPrepare(msg)
	}
	cache := newMsgCache(isValidMsg)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case notification := <-sub:
			cache = cache.Add(notification.Unwrap())

			prepares := cache.Get()
			if !ctx.Quorum().HasQuorum(ibft.WrapMessages(prepares...)) {
				continue
			}

			return prepares, nil
		}
	}
}

func (s *Sequencer) isValidMsgPrepare(msg *types.MsgPrepare) bool {
	if !s.IsValidator(msg.From, msg.View.Sequence) {
		return false
	}

	if !bytes.Equal(msg.BlockHash, s.state.AcceptedBlockHash()) {
		return false
	}

	return true
}
