package objects

import "testing"

func TestIntValue(t *testing.T) {
	var value Value = Int(1712)

	if value.Type() != IntType {
		t.Errorf("got %s, expected %s", value.Type(), IntType)
	}

	if value.String() != "1712" {
		t.Errorf("got %q, expected %q", value.String(), "1712")
	}
}

func TestBoolValue(t *testing.T) {
	var value Value = Bool(true)

	if value.Type() != BoolType {
		t.Errorf("got %s, expected %s", value.Type(), BoolType)
	}

	if value.String() != "true" {
		t.Errorf("got %q, expected %q", value.String(), "true")
	}
}
