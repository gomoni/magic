package magic

/*
Package magis is opinionated libmagic(1) gco wrapper returning mime type of a
data. It does provide golang native interface and does not those, who
are not interesting.

Additionally implements bounded pool of magic cookies, so that library can be
called from gorutines in a safe way.
*/

/*
#include<stdlib.h>
#include<magic.h>
#cgo LDFLAGS: -lmagic
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

type flag int

const (
	// MagicMime - return mime type; encoding
	MagicMime flag = C.MAGIC_MIME
)

var (
	// FIXME!!! why wrapped magic did not find the default DB??
	emptyCString    *C.char = C.CString("/usr/share/misc/magic")
	errEmptyRequest         = errors.New("Empty request")
)

// Cookie - wrapper for C magic_t of libmagic
type Cookie struct {
	cookie C.magic_t
}

// magicError - return output of magic_error as Go error
// panic if there was no error
func magicError(cookie C.magic_t) error {

	if cookie == nil {
		return fmt.Errorf("Can't allocate magic cookie")
	}

	return fmt.Errorf(C.GoString(C.magic_error(cookie)))
}

// New - allocate new magic cookie and set flags to MAGIC_MIME
// returns an error in a case of failure
func New() (*Cookie, error) {
	c := C.magic_open(C.int(0))
	if c == nil {
		return nil, magicError(c)
	}

	r := C.magic_setflags(c, C.int(MagicMime))
	if r == -1 {
		goto err
	}

	r = C.magic_load(c, emptyCString)
	if r == -1 {
		goto err
	}

	return &Cookie{cookie: c}, nil

err:
	err := magicError(c)
	C.magic_close(c)
	return nil, err
}

// Close - deallocate the cookie
func (m *Cookie) Close() {
	C.magic_close(m.cookie)
}

// Mime - reads the content of a buffer and reports `mime; encoding`
// or an error if that happen
// Note: it is NOT gorutine safe
func (m *Cookie) Mime(r io.Reader) (mime string, err error) {
	var b []byte
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	if len(b) == 0 {
		err = errEmptyRequest
		return
	}

	cb := unsafe.Pointer(&b[0])
	cmime := C.magic_buffer(m.cookie, cb, C.size_t(len(b)))
	if cmime == nil {
		err = magicError(m.cookie)
		return
	}
	return C.GoString(cmime), nil
}
