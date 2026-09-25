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

func TestSameType(t *testing.T) {
	tests := []struct {
		name string
		a    Type
		b    Type
		want bool
	}{
		{
			name: "same primitive type",
			a:    IntType,
			b:    IntType,
			want: true,
		},
		{
			name: "different primitive types",
			a:    IntType,
			b:    BoolType,
			want: false,
		},
		{
			name: "same function signature",
			a: &FunctionType{
				ParameterTypes: []Type{IntType, BoolType},
				ReturnType:     IntType,
			},
			b: &FunctionType{
				ParameterTypes: []Type{IntType, BoolType},
				ReturnType:     IntType,
			},
			want: true,
		},
		{
			name: "different parameter type",
			a: &FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     BoolType,
			},
			b: &FunctionType{
				ParameterTypes: []Type{BoolType},
				ReturnType:     BoolType,
			},
			want: false,
		},
		{
			name: "different parameter count",
			a: &FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     BoolType,
			},
			b: &FunctionType{
				ParameterTypes: []Type{IntType, IntType},
				ReturnType:     BoolType,
			},
			want: false,
		},
		{
			name: "different return type",
			a: &FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     BoolType,
			},
			b: &FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     IntType,
			},
			want: false,
		},
		{
			name: "nested function types",
			a: &FunctionType{
				ParameterTypes: []Type{
					&FunctionType{
						ParameterTypes: []Type{IntType},
						ReturnType:     BoolType,
					},
				},
				ReturnType: IntType,
			},
			b: &FunctionType{
				ParameterTypes: []Type{
					&FunctionType{
						ParameterTypes: []Type{IntType},
						ReturnType:     BoolType,
					},
				},
				ReturnType: IntType,
			},
			want: true,
		},
		{
			name: "function versus primitive",
			a: &FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     IntType,
			},
			b:    IntType,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SameType(tt.a, tt.b); got != tt.want {
				t.Fatalf("SameType(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}

			if got := SameType(tt.b, tt.a); got != tt.want {
				t.Fatalf(
					"SameType(%v, %v) = %v, want %v (reversed)",
					tt.b,
					tt.a,
					got,
					tt.want,
				)
			}

		})
	}
}

func TestSameTypeNestedFunctionTypes(t *testing.T) {
	a := &FunctionType{
		ParameterTypes: []Type{
			&FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     BoolType,
			},
		},
		ReturnType: IntType,
	}

	b := &FunctionType{
		ParameterTypes: []Type{
			&FunctionType{
				ParameterTypes: []Type{IntType},
				ReturnType:     BoolType,
			},
		},
		ReturnType: IntType,
	}

	if !SameType(a, b) {
		t.Fatal("expected nested function types to be equal")
	}
}
