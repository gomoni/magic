package magic

/*
Package magis is opinionated libmagic(1) gco wrapper returning mime type of a
data. It does provide golang native interface and does not those, who
are not interesting.
*/

/*
#include<stdlib.h>
#include<magic.h>
#cgo LDFLAGS: -lmagic
*/
import "C"

import (
    "io"
    "io/ioutil"
    "fmt"
    "unsafe"
)

type Flag int

const (
    // print mime type; encoding
    MAGIC_MIME Flag = C.MAGIC_MIME
)

var (
// FIXME!!! why wrapped magic did not find the default DB??
    emptyCString *C.char = C.CString("/usr/share/misc/magic")
)


type Magic struct {
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
func New() (*Magic, error) {
    c := C.magic_open(C.int(0))
    if c == nil {
        return nil, magicError(c)
    }

    r := C.magic_setflags(c, C.int(MAGIC_MIME))
    if r == -1 {
        goto err
    }

    r = C.magic_load(c, emptyCString)
    if r == -1 {
        goto err
    }

    return &Magic{cookie: c}, nil

err:
    err := magicError(c)
    C.magic_close(c)
    return nil, err
}

// Close - deallocate the cookie
func (m *Magic) Close() {
    C.magic_close(m.cookie)
}

// Mime - reads the content of a buffer and reports `mime; encoding`
// or an error if that happen
func (m *Magic) Mime(r io.Reader) (mime string, err error) {
    var b []byte
    b, err = ioutil.ReadAll(r)
    if err != nil {
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
