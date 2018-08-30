package dxda_test

import (
	"testing"

	"github.com/geetduggal/dxda"
)

func TestHello(t *testing.T) {
	if dxda.Hello() != "Hello!" {
		t.Errorf("dxda.Hello() != 'Hello!'")
	}
}
