package checker

import (
	"testing"

	"github.com/shsiddhant/gotiny/internal/objects"
)

func TestEnvironmentDefineAndGet(t *testing.T) {
	env := NewTypeEnvironment()

	err := env.Define("x", objects.IntType)
	if err != nil {
		t.Fatalf("unexpected error defining x: %v", err)
	}

	typ, err := env.Get("x")
	if err != nil {
		t.Fatalf("unexpected error getting x: %v", err)
	}
	if typ != objects.IntType {
		t.Errorf("expected IntType, got %v", typ)
	}
}

func TestEnvironmentOuterLookup(t *testing.T) {
	parent := NewTypeEnvironment()
	_ = parent.Define("a", objects.IntType)

	child := parent.NewChild()

	typ, err := child.Get("a")
	if err != nil {
		t.Fatalf("failed to look up variable from outer scope: %v", err)
	}
	if typ != objects.IntType {
		t.Errorf("expected IntType, got %v", typ)
	}
}

func TestEnvironmentShadowing(t *testing.T) {
	parent := NewTypeEnvironment()
	_ = parent.Define("x", objects.IntType)

	child := parent.NewChild()
	_ = child.Define("x", objects.BoolType)

	childType, _ := child.Get("x")
	if childType != objects.BoolType {
		t.Errorf("expected child to shadow x as BoolType, got %v", childType)
	}

	parentType, _ := parent.Get("x")
	if parentType != objects.IntType {
		t.Errorf("expected parent x to remain IntType, got %v", parentType)
	}
}

func TestEnvironmentDuplicateDefinitionInSameScope(t *testing.T) {
	env := NewTypeEnvironment()
	_ = env.Define("y", objects.IntType)

	err := env.Define("y", objects.IntType)
	if err == nil {
		t.Error("expected error for duplicate definition in the same scope, got nil")
	}
}

func TestEnvironmentScopeIsolation(t *testing.T) {
	parent := NewTypeEnvironment()

	child := parent.NewChild()
	_ = child.Define("innerOnly", objects.IntType)

	_, err := parent.Get("innerOnly")
	if err == nil {
		t.Error("expected error when looking up inner scope variable from parent, got nil")
	}
}
