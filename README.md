# magic

Opiniated `libmagic` Go wrapper, does not allow any configuration and operates on `io.Reader` only.

## TODO

Memory sanitizer reports unitialized memory error

```
CC=clang go test -msan
Uninitialized bytes in __interceptor_memchr at offset 0 inside [0x729000000000, 60)
==13338==WARNING: MemorySanitizer: use-of-uninitialized-value
    #0 0x7f38f36b31ea  (/usr/lib64/libmagic.so.1+0xb1ea)
    #1 0x7f38f36b44e7  (/usr/lib64/libmagic.so.1+0xc4e7)
    #2 0x7f38f36b549b  (/usr/lib64/libmagic.so.1+0xd49b)
    #3 0x7f38f36b4391  (/usr/lib64/libmagic.so.1+0xc391)
    #4 0x7f38f36b6131  (/usr/lib64/libmagic.so.1+0xe131)
    #5 0x7f38f36b6653  (/usr/lib64/libmagic.so.1+0xe653)
    #6 0x7f38f36bea7e  (/usr/lib64/libmagic.so.1+0x16a7e)
    #7 0x7f38f36ad215 in magic_buffer (/usr/lib64/libmagic.so.1+0x5215)
    #8 0x5e854e in _cgo_dee611482f22_Cfunc_magic_buffer (/tmp/go-build583825240/b001/magic.test+0x5e854e)
    #9 0x4e9f1f  (/tmp/go-build583825240/b001/magic.test+0x4e9f1f)

SUMMARY: MemorySanitizer: use-of-uninitialized-value (/usr/lib64/libmagic.so.1+0xb1ea) 
Exiting
exit status 77
FAIL    _/home/mvyskocil/projects/vyskocilm/gomoni/magic        0.031s
mvyskocil@linux-33zu:~/projects/vyskocilm/gomoni/magic> 
```
