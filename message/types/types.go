package types

// FinalizedSeal is proof that a specific validator committed to a block
type FinalizedSeal struct {
	From, CommitSeal []byte
}

// FinalizedBlock is the consensus verified block
type FinalizedBlock struct {
	// block that was finalized
	Block []byte

	// seals of validators who committed to this block
	Seals []FinalizedSeal

	// round in which the block was finalizeds
	Round uint64
}

// Msg is a convenience wrapper for the consensus messages
type Msg interface {
	// GetFrom returns the address associated with this Msg
	GetFrom() []byte

	// GetSignature returns the signature of this Msg
	GetSignature() []byte

	// Payload returns the byte content of this Msg (signature excluded)
	Payload() []byte
}

// MsgFeed provides an asynchronous way to receive consensus messages. In addition
// to listen for any type of message for any particular view, the higherRounds flag provides an option
// to receive messages from rounds higher than the round in provided view.
//
// CONTRACT: messages received by consuming the channel's callback are assumed to be valid:
//
// - any message has a valid view (matches the one provided)
//
// - any message has a valid signature (validator signed the message payload)
//
// - all messages are considered unique (there cannot be 2 or more messages with identical From fields)
type MsgFeed interface {
	// Proposal returns the MsgProposal subscription for given view(s)
	Proposal(view *View, higherRounds bool) (<-chan func() []*MsgProposal, func())

	// Prepare returns the MsgPrepare subscription for given view(s)
	Prepare(view *View, higherRounds bool) (<-chan func() []*MsgPrepare, func())

	// Commit returns the MsgCommit subscription for given view(s)
	Commit(view *View, higherRounds bool) (<-chan func() []*MsgCommit, func())

	// RoundChange returns the MsgRoundChange subscription for given view(s)
	RoundChange(view *View, higherRounds bool) (<-chan func() []*MsgRoundChange, func())
}

// Signer signs data with its private key
type Signer interface {
	// Sign returns the signature generated by signing the provided input
	Sign([]byte) []byte
}

// SigRecover is used to validate consensus messages are properly signed
type SigRecover interface {
	// From returns the address associated with this signature
	From(data []byte, sig []byte) []byte
}