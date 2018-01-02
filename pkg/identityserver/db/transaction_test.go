// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package db

import "testing"

func TestExecTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testExec(t, tx)
		return nil
	})
}

func TestNamedExecTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testNamedExec(t, tx)
		return nil
	})
}

func TestSelectTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testSelect(t, tx)
		return nil
	})
}

func TestNamedSelectTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testNamedSelect(t, tx)
		return nil
	})
}

func TestSelectOneTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testSelectOne(t, tx)
		return nil
	})
}

func TestNamedSelectOneTx(t *testing.T) {
	getInstance(t).Transact(func(tx *Tx) error {
		testNamedSelectOne(t, tx)
		return nil
	})
}
