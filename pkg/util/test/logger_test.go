// Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package test

import "testing"

func TestGetLogger(t *testing.T) {
	logger := GetLogger(t)
	logger.Debug("abcabcabc - Hi!")
	logger.Info("Fooz")
	logger.Errorf("Nope %d", 1234)
}
