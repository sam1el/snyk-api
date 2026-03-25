package output

import (
	"bytes"
	"testing"
)

func TestFormatterTemplate(t *testing.T) {
	data := map[string]string{"name": "alice"}
	var buf bytes.Buffer

	f := New("json").WithWriter(&buf).WithTemplate("{{.name}}")
	if err := f.Print(data); err != nil {
		t.Fatalf("template print failed: %v", err)
	}

	if got := buf.String(); got != "alice" {
		t.Fatalf("expected 'alice', got %q", got)
	}
}

func TestFormatterJQ(t *testing.T) {
	data := map[string]interface{}{"name": "alice", "age": 30}
	var buf bytes.Buffer

	f := New("json").WithWriter(&buf).WithJQ(".name")
	if err := f.Print(data); err != nil {
		t.Fatalf("jq print failed: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("alice")) {
		t.Fatalf("expected jq output to contain alice, got %q", buf.String())
	}
}
