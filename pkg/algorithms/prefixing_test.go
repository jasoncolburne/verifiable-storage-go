package algorithms_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func TestPrefixing(t *testing.T) {
	if err := testPrefixing(); err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
}

func testPrefixing() error {
	prefixer := &primitives.Prefixer{}

	if err := algorithms.CreatePrefix(prefixer); err != nil {
		return err
	}

	expectedPrefix := `EE9sNjnroj4p9NQXEEhuMNu4yQPyi2whAbNE9wJjnaJM`
	badPrefix := `EE9sNjnroj4p9NQXEEhuMNu4yQPyi2whAbNE9wJjnaJN`

	if !strings.EqualFold(prefixer.Id, expectedPrefix) {
		return fmt.Errorf("unexpected id: %s", prefixer.Id)
	}

	if !strings.EqualFold(prefixer.Prefix, expectedPrefix) {
		return fmt.Errorf("unexpected prefix: %s", prefixer.Prefix)
	}

	if err := algorithms.VerifyPrefixAndData(prefixer); err != nil {
		return err
	}

	prefixer.Id = badPrefix

	if err := algorithms.VerifyPrefixAndData(prefixer); err == nil {
		return fmt.Errorf("unexpected verification success with bad id")
	}

	prefixer.Id = expectedPrefix
	prefixer.Prefix = badPrefix

	if err := algorithms.VerifyPrefixAndData(prefixer); err == nil {
		return fmt.Errorf("unexpected verification success with bad prefix")
	}

	return nil
}
