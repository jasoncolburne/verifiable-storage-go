package primitives

type VerifiableAndRecordable interface {
	Prefixable // this depends on SelfAddressable
	Sequenceable
	Chainable
	Nonceable
	Timestampable
}

type SearchableAndRecordable interface {
	VerifiableAndRecordable
	Searchable
	DeriveSearchKey() string // up to the implementer
}

type SignableAndRecordable interface {
	VerifiableAndRecordable
	Signable
}

type SignableAndSearchableAndRecordable interface {
	SearchableAndRecordable
	Signable
}

type VerifiableRecorder struct {
	Prefixer    // [id from SelfAddresser and] prefix
	Sequencer   // sequenceNumber
	Chainer     // previous
	Noncer      // nonce
	Timestamper // created_at
}

type SearchableRecorder struct {
	VerifiableRecorder
	Searcher
}

type SignableRecorder struct {
	VerifiableRecorder
	Signer // signingIdentity, signature
}

type SignableSearchableRecorder struct {
	SearchableRecorder
	Signer
}
