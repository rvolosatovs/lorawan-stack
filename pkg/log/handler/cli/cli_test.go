// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package cli

import (
	"bufio"
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/TheThingsNetwork/ttn/pkg/log"
	"github.com/TheThingsNetwork/ttn/pkg/log/test"
	"github.com/smartystreets/assertions"
)

func TestHandlerNewColors(t *testing.T) {
	a := assertions.New(t)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	// COLORTERM= TERM= does not enable colors
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM", "")

	a.So(New(w).UseColor, assertions.ShouldBeFalse)

	// COLORTERM=0 forces colors off
	os.Setenv("COLORTERM", "0")
	os.Setenv("TERM", "colorterm")

	a.So(New(w).UseColor, assertions.ShouldBeFalse)

	// TERM with correct substring turns colors on
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM", "colorterm")

	a.So(New(w).UseColor, assertions.ShouldBeTrue)

	// TERM with correct substring turns colors on
	os.Setenv("COLORTERM", "")
	os.Setenv("TERM", "xterm")

	a.So(New(w).UseColor, assertions.ShouldBeTrue)

	// COLORTERM=1 turns colors on
	os.Setenv("COLORTERM", "1")
	os.Setenv("TERM", "")

	a.So(New(w).UseColor, assertions.ShouldBeTrue)

	// COLORTERM=1 turns colors on
	os.Setenv("COLORTERM", "1")
	os.Setenv("TERM", "")

	// but UseColor(false) turns it off again
	a.So(New(w, UseColor(false)).UseColor, assertions.ShouldBeFalse)

	// COLORTERM=1 turns colors off
	os.Setenv("COLORTERM", "0")
	os.Setenv("TERM", "")

	// but UseColor(true) turns it off again
	a.So(New(w, UseColor(true)).UseColor, assertions.ShouldBeTrue)
}

func TestHandlerHandleLog(t *testing.T) {
	a := assertions.New(t)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	handler := New(w, UseColor(false))

	err := handler.HandleLog(&test.Entry{
		M: "Foo",
		L: log.Debug,
		T: time.Now(),
		F: log.Fields("a", 10, "b", "bar", "c", false, "d", 33.4),
	})
	a.So(err, assertions.ShouldBeNil)

	str := " DEBUG Foo                                      a=10 b=bar c=false d=33.4\n"

	err = w.Flush()
	a.So(err, assertions.ShouldBeNil)
	a.So(b.String(), assertions.ShouldEqual, str)
}
