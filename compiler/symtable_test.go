package compiler

import (
	"testing"
)

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Scope: SCOPE_GLOBAL, Index: 0},
		"b": {Name: "b", Scope: SCOPE_GLOBAL, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{Name: "a", Scope: SCOPE_GLOBAL, Index: 0},
		{Name: "b", Scope: SCOPE_GLOBAL, Index: 1},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	expected := []Symbol{
		{Name: "a", Scope: SCOPE_GLOBAL, Index: 0},
		{Name: "b", Scope: SCOPE_GLOBAL, Index: 1},
		{Name: "c", Scope: SCOPE_LOCAL, Index: 0},
		{Name: "d", Scope: SCOPE_LOCAL, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c")
	firstLocal.Define("d")

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e")
	secondLocal.Define("f")

	tests := []struct {
		table    *SymbolTable
		expected []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Scope: SCOPE_GLOBAL, Index: 0},
				{Name: "b", Scope: SCOPE_GLOBAL, Index: 1},
				{Name: "c", Scope: SCOPE_LOCAL, Index: 0},
				{Name: "d", Scope: SCOPE_LOCAL, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Scope: SCOPE_GLOBAL, Index: 0},
				{Name: "b", Scope: SCOPE_GLOBAL, Index: 1},
				{Name: "e", Scope: SCOPE_LOCAL, Index: 0},
				{Name: "f", Scope: SCOPE_LOCAL, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expected {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()

	expected := []Symbol{
		{Name: "yarn", Scope: SCOPE_BUILTIN, Index: 0},
		{Name: "len", Scope: SCOPE_BUILTIN, Index: 1},
		{Name: "type", Scope: SCOPE_BUILTIN, Index: 2},
	}

	for i, sym := range expected {
		global.DefineBuiltin(i, sym.Name)
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got=%+v", sym.Name, sym, result)
		}
	}
}

func TestResolveUnresolvable(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	_, ok := global.Resolve("b")
	if ok {
		t.Error("expected 'b' to be unresolvable, but it was resolved")
	}

	local := NewEnclosedSymbolTable(global)
	local.Define("c")

	_, ok = local.Resolve("d")
	if ok {
		t.Error("expected 'd' to be unresolvable, but it was resolved")
	}
}

func TestNumDefinitions(t *testing.T) {
	global := NewSymbolTable()

	if n := global.NumDefinitions(); n != 0 {
		t.Errorf("expected 0 definitions, got %d", n)
	}

	global.Define("a")
	if n := global.NumDefinitions(); n != 1 {
		t.Errorf("expected 1 definition, got %d", n)
	}

	global.Define("b")
	global.Define("c")
	if n := global.NumDefinitions(); n != 3 {
		t.Errorf("expected 3 definitions, got %d", n)
	}
}
