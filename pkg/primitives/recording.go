package primitives

type VerifiableAndRecordable interface {
	Prefixable // this depends on SelfAddressable
	Sequenceable
	Chainable
	Nonceable
	Timestampable
}

type SignableAndRecordable interface {
	VerifiableAndRecordable
	Signable
}

type VerifiableRecorder struct {
	Prefixer    // [id from SelfAddresser and] prefix
	Sequencer   // sequenceNumber
	Chainer     // previous
	Noncer      // nonce
	Timestamper // created_at
}

type SignableRecorder struct {
	VerifiableRecorder
	Signer // signingIdentity, signature
}
