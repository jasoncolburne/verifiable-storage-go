# verifiable-storage-go
A prefixable, self-addressing (tamper evident), sequenced, chained data store with optional signing.

Ideas here were derived from KERI's KELs.

## Getting Started

To start, look in `pkg/repository/repository_test.go`.

This provides an interface for versioned, verifiable, optionally signed model repositories with
minimal setup.

Data is never deleted, as this is designed to support a decentralized deployment and if you release
data into the wild, it can never be undone - you can at most append to it. There are no deletes or
updates in this api.

The notion of a `prefix` is one where the first id in the chain of record versions represents the
entire chain. This prefix is embedded in each record and does not change.

The typical pattern then, is:

1. `GetLatestByPrefix()`
2. Modify data
3. `CreateVersion()`

That said, a couple other APIs are supported (`GetById()` and `ListByPrefix()`).

## Concepts

- **Chains**: Like a blockchain, each record (other than the first) points to the previous record
with a hash commitment.
- **Nonces**: Each record contains a nonce to add uniqueness. In some cases this may be undesirable
and likely warrants configuration, but for the majority of cases this is what you need.
- **Prefixes**: A self-address derived during creation of the first record in a chain. This value is
both the id and the prefix of that record, and it is the prefix of future records in the chain.
- **Self-Addresses**: A self-address is an id (named `id`) embedded in the data itself that is derived
from the data. This provides tamper evidence, for if either the identifier or data is modified, the
verification fails.
- **Sequence Numbers**: Each record contains a sequence number that increments monotonically. This,
coupled with a unique constraint, provides a very good solution to divergence prevention (two
writes to the same chain of data based on the same record - both would have the same sequence
number).
- **Timestamping**: Each record is timestamped.

### Optional

- **Signing**: Records can be signed and when they are, two fields are added. One for the signature
itself, and the other to identify the signer.

## Verification

As data is read from a repository, it is verified for id/data validity, and if signed, the
signature is also verified.
