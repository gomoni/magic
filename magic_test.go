package magic

import (
    "testing"
    "fmt"
    "os"
)

func TestMagic(t *testing.T) {

    m, err := New()
    if err != nil {
        t.Fatal(err)
    }
    defer m.Close()

    f, err := os.Open("/etc/passwd")
    if err != nil {
        t.Fatal(err)
    }
    defer f.Close()

    mime, err := m.Mime(f)
    if err != nil {
        t.Fatal(err)
    }

    if mime != "text/plain; charset=us-ascii" {
        t.Fatalf("Wrong mime type, expected 'text/plain; charset-us-ascii', got '%s'", mime)
    }
    fmt.Printf("magic /etc/passwd: %s\n", mime)

}
