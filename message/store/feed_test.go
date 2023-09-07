package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madz-lab/go-ibft/message/types"
)

func TestFeed_MsgProposal(t *testing.T) {
	t.Parallel()

	codec := mockCodec{func(_, _ []byte) []byte { return nil }}

	t.Run("msg received", func(t *testing.T) {
		t.Parallel()

		var (
			view = &types.View{Sequence: 101, Round: 0}
			msg  = &types.MsgProposal{
				View:      view,
				Signature: []byte("sig"),
			}
		)

		store := New(codec)
		require.NoError(t, store.AddMsgProposal(msg))

		sub, cancelSub := Feed{store}.SubscribeToProposalMessages(view, false)
		defer cancelSub()

		unwrap := <-sub
		messages := unwrap()

		assert.Equal(t, msg, messages[0])
	})

	t.Run("future round msg received", func(t *testing.T) {
		t.Parallel()

		var (
			view = &types.View{Sequence: 101, Round: 1}
			msg  = &types.MsgProposal{
				View:      view,
				Signature: []byte("signature 2"),
			}
		)

		store := New(codec)
		require.NoError(t, store.AddMsgProposal(msg))
		require.Len(t, store.GetProposalMessages(view), 1)

		previousView := &types.View{Sequence: view.Sequence, Round: view.Round - 1}
		sub, cancelSub := Feed{store}.SubscribeToProposalMessages(previousView, true)
		defer cancelSub()

		unwrap := <-sub
		messages := unwrap()

		assert.Equal(t, msg, messages[0])
	})

	t.Run("highest round msg received", func(t *testing.T) {
		t.Parallel()

		store := New(codec)

		sub, cancelSub := Feed{store}.SubscribeToProposalMessages(&types.View{
			Sequence: 101,
			Round:    0,
		},
			true,
		)
		defer cancelSub()

		var (
			view1 = &types.View{Sequence: 101, Round: 1}
			msg1  = &types.MsgProposal{
				View:      view1,
				Signature: []byte("signature"),
			}

			view2 = &types.View{Sequence: 101, Round: 5}
			msg2  = &types.MsgProposal{
				View:      view2,
				Signature: []byte("signature"),
			}

			view3 = &types.View{Sequence: 101, Round: 10}
			msg3  = &types.MsgProposal{
				View:      view3,
				Signature: []byte("signature"),
			}
		)

		require.NoError(t, store.AddMsgProposal(msg2))
		require.NoError(t, store.AddMsgProposal(msg3))
		require.NoError(t, store.AddMsgProposal(msg1))
		require.Len(t, store.GetProposalMessages(view1), 1)
		require.Len(t, store.GetProposalMessages(view2), 1)
		require.Len(t, store.GetProposalMessages(view3), 1)

		unwrap := <-sub
		msgs := unwrap()

		assert.Equal(t, msg3, msgs[0])
	})
}