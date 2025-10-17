package algorithms_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jasoncolburne/verifiable-storage-go/pkg/algorithms"
	"github.com/jasoncolburne/verifiable-storage-go/pkg/primitives"
)

func TestSelfAddressing(t *testing.T) {
	if err := testSelfAddressing(); err != nil {
		fmt.Printf("%s\n", err)
		t.FailNow()
	}
}

func testSelfAddressing() error {
	addresser := &primitives.SelfAddresser{}

	if err := algorithms.SelfAddress(addresser); err != nil {
		return err
	}

	expectedId := `EIuB8-qRNMMGsLpJQFMgeJxWr_ppYahDfQh6mgvkdD2R`
	badId := `EIuB8-qRNMMGsLpJQFMgeJxWr_ppYahDfQh6mgvkdD2S`

	if !strings.EqualFold(addresser.Id, expectedId) {
		return fmt.Errorf("unexpected id: %s", addresser.Id)
	}

	if err := algorithms.VerifyAddressAndData(addresser); err != nil {
		return err
	}

	addresser.Id = badId

	if err := algorithms.VerifyAddressAndData(addresser); err == nil {
		return fmt.Errorf("unexpected verification success with bad id")
	}

	return nil
}
