package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madz-lab/go-ibft/message/types"
)

type sigRecoverFn func([]byte, []byte) []byte

func (s sigRecoverFn) From(data, sig []byte) []byte {
	return s(data, sig)
}

func TestFeed_MsgProposal(t *testing.T) {
	t.Parallel()

	sigRecover := sigRecoverFn(func(_ []byte, _ []byte) []byte { return nil })

	t.Run("msg received", func(t *testing.T) {
		t.Parallel()

		var (
			view = &types.View{Sequence: 101, Round: 0}
			msg  = &types.MsgProposal{
				View:      view,
				Signature: []byte("sig"),
			}
		)

		store := New(sigRecover)
		require.NoError(t, AddMessage[types.MsgProposal](msg, store))

		sub, cancelSub := Feed{store}.Proposal(view, false)
		defer cancelSub()

		unwrap := <-sub
		messages := unwrap()

		assert.Equal(t, msg, messages[0])
	})

	t.Run("msg received with exact view", func(t *testing.T) {
		t.Parallel()

		var (
			view1 = &types.View{Sequence: 101, Round: 0}
			msg1  = &types.MsgProposal{
				View:      view1,
				Signature: []byte("sig"),
			}

			view2 = &types.View{Sequence: 101, Round: 5}
			msg2  = &types.MsgProposal{
				View:      view2,
				Signature: []byte("sig"),
			}
		)

		store := New(sigRecover)
		require.NoError(t, AddMessage[types.MsgProposal](msg1, store))
		require.NoError(t, AddMessage[types.MsgProposal](msg2, store))

		sub, cancelSub := Feed{store}.Proposal(view1, false)
		defer cancelSub()

		unwrap := <-sub
		messages := unwrap()

		assert.Equal(t, msg1, messages[0])
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

		store := New(sigRecover)
		require.NoError(t, AddMessage[types.MsgProposal](msg, store))
		require.Len(t, GetMessages[types.MsgProposal](view, store), 1)

		previousView := &types.View{Sequence: view.Sequence, Round: view.Round - 1}
		sub, cancelSub := Feed{store}.Proposal(previousView, true)
		defer cancelSub()

		unwrap := <-sub
		messages := unwrap()

		assert.Equal(t, msg, messages[0])
	})

	t.Run("highest round msg received", func(t *testing.T) {
		t.Parallel()

		store := New(sigRecover)

		sub, cancelSub := Feed{store}.Proposal(&types.View{
			Sequence: 101,
			Round:    6,
		},
			true,
		)
		defer cancelSub()

		unwrap := <-sub
		assert.Len(t, unwrap(), 0)

		var (
			view1 = &types.View{Sequence: 101, Round: 1}
			msg1  = &types.MsgProposal{
				View:      view1,
				Signature: []byte("signature"),
			}

			view3 = &types.View{Sequence: 101, Round: 10}
			msg3  = &types.MsgProposal{
				View:      view3,
				Signature: []byte("signature"),
			}
		)

		require.NoError(t, AddMessage[types.MsgProposal](msg3, store))
		require.NoError(t, AddMessage[types.MsgProposal](msg1, store))

		unwrap = <-sub
		msgs := unwrap()

		require.Len(t, msgs, 1)
		assert.Equal(t, msg3, msgs[0])
	})

	t.Run("subscription not notified", func(t *testing.T) {
		t.Parallel()

		store := New(sigRecover)

		view1 := &types.View{Sequence: 101, Round: 1}
		view2 := &types.View{Sequence: 102, Round: 1}

		// two subscriptions, same view
		sub, cancelSub := Feed{store}.Proposal(view1, true)

		unwrap := <-sub
		require.Len(t, unwrap(), 0)

		msg := &types.MsgProposal{
			View:      view2,
			Signature: []byte("signature"),
		}

		require.NoError(t, AddMessage[types.MsgProposal](msg, store))

		cancelSub() // close the sub so the channel can be read
		_, ok := <-sub
		assert.False(t, ok)
	})

	t.Run("subscription gets latest notification", func(t *testing.T) {
		t.Parallel()

		store := New(sigRecover)

		view1 := &types.View{Sequence: 101, Round: 1}
		view2 := &types.View{Sequence: 101, Round: 2}

		// two subscriptions, same view
		sub, cancelSub := Feed{store}.Proposal(view1, true)
		defer cancelSub()

		var (
			msg1 = &types.MsgProposal{
				View:      view1,
				Signature: []byte("signature"),
			}

			msg2 = &types.MsgProposal{
				View:      view2,
				Signature: []byte("signature"),
			}
		)

		require.NoError(t, AddMessage[types.MsgProposal](msg1, store))
		require.NoError(t, AddMessage[types.MsgProposal](msg2, store))

		unwrap := <-sub
		messages := unwrap()
		require.Len(t, messages, 1)
		assert.Equal(t, msg2, messages[0])
	})
}
