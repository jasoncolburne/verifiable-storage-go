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

That said, a few other direct APIs are supported (`GetById()` and `ListByPrefix()`), and some
generic APIs exist (`Get()`, `Select()`, and `ListLatestByPrefix()`)

### ListLatestByPrefix()

ListLatestByPrefix deserves some discussion. It permits this kind of thing:

```go
accountRecord := &AccountRecord{
    ...
    Active: true,
}
r.CreateVersion(accountRecord)

// later on...

accountRecord.Active = false
r.CreateVersion(accountRecord)

// the problem is that after this point, if you did a regular select for active accounts, you'd get
// versions of the record that were labeled as active, when really the most recent version has
// rendered it inactive.

// the below method accomodates this by first filtering down to the most recent records in sequences
// only, using prefix as a partitioning field. then, the supplied conditions are applied.

r.ListLatestByPrefix(
    ctx,
    &records, 
    expressions.Equal("active", true),
    nil,
    nil
)

// records will not contain any versions of `accountRecord`. if you had performed a regular select,
// you'd have all the old versions of `accountRecord` which were marked active.

```


## Concepts

- **Chains**: Like a blockchain, each record (other than the first) points to the previous record
with a hash commitment.
- **Prefixes**: A self-address derived during creation of the first record in a chain. This value is
both the id and the prefix of that record, and it is the prefix of future records in the chain.
- **Self-Addresses**: A self-address is an id (named `id`) embedded in the data itself that is derived
from the data. This provides tamper evidence, for if either the identifier or data is modified, the
verification fails.
- **Sequence Numbers**: Each record contains a sequence number that increments monotonically. This,
coupled with a unique constraint, provides a very good solution to divergence prevention (two
writes to the same chain of data based on the same record - both would have the same sequence
number).

### Optional

In the implementation, repositories come in only two flavours (verifiable and signed). These two
types can (and most often should) be configured with optional nonces and timestamps.

- **Nonces**: A record may contain a nonce to add uniqueness. In some cases this may be undesirable,
but for the majority of cases this is what you need. If you want more determinism (duplicate
prevention for instance) you can supply a nil noncer to the repository creation method and the field
will be omitted. Be sure to disable both nonces and timestamping for true determinism.
- **Timestamping**: Each record may be timestamped. If you want determinism and can tolerate the
lack of a timestamp, disable this and nonces.
- **Signing**: Records can be signed and when they are, two fields are added. One for the signature
itself, and the other to identify the signer.

It's worth noting that you'll still have a `CreatedAt` and `Nonce` field on the struct you're using
even if you disable them (as pointers). Just don't assign them, the code omits them from writes
and computations if they aren't set.

## Verification

As data is read from a repository, it is verified for id/data validity, and if signed, the
signature is also verified.
