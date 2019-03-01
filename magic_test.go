package magic

import (
	"fmt"
	"os"
	"testing"
)

func TestMagic(t *testing.T) {

	if _, err := os.Stat("/etc/passwd"); os.IsNotExist(err) {
		t.Skip("/etc/passwd does not exists")
	}

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

	if testing.Verbose() {
		fmt.Printf("magic /etc/passwd: %s\n", mime)
	}

}
