package gribble

import (
	"testing"
)

var env = New([]Command{
	add{},
})

type add struct {
	Op1 int `param:"1"`
	Op2 int `param:"2"`
}

func (cmd add) Run() Value {
	return cmd.Op1 + cmd.Op2
}

// TestAdd tests whether a simple add command is working.
func TestAdd(t *testing.T) {
	val, err := env.Run("add 35 7")
	if err != nil {
		t.Fatalf("Using the 'add' command failed with an error: %s", err)
	}
	ival, ok := val.(int)
	if !ok {
		t.Fatalf("The 'add' command should return an int, but it returned %T.",
			val)
	}
	if ival != 42 {
		t.Fatalf("The 'add 35 7' command should return 42, but returned %d.",
			ival)
	}
}

// TestSubCommand tests whether sub-commands are working.
func TestSubCommand(t *testing.T) {
	val, err := env.Run("add (add 30 5) 7")
	if err != nil {
		t.Fatalf("Using the 'add' command failed with an error: %s", err)
	}
	ival, ok := val.(int)
	if !ok {
		t.Fatalf("The 'add' command should return an int, but it returned %T.",
			val)
	}
	if ival != 42 {
		t.Fatalf("The 'add (add 30 5) 7' command should return 42, but "+
			"returned %d.", ival)
	}
}

// TestMalformed tests whether Gribble returns an error if given a malformed
// command. (This includes commands that don't exist, parse errors, wrongly
// typed parameters, too many/few parameters, etc.)
func TestMalformed(t *testing.T) {
	cmds := []string{
		"add (2 3)", "add (2 3", "dne 1", "", "add 1 2 3", "add 1.0 2.0",
		"add 1",
	}
	for _, cmd := range cmds {
		val, err := env.Run(cmd)
		if err == nil {
			t.Fatalf("Gribble should have choked on '%s', but instead "+
				"returned a value: '%v'.", cmd, val)
		}
	}
}

type badNonContiguousParams struct {
	Op1 int `param:"1"`
	Op2 int `param:"3"`
}

func (cmd badNonContiguousParams) Run() Value { return nil }

// TestNonContiguousParams tests whether Gribble catches inconsistent parameter
// numbering.
func TestNonContiguousParams(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badNonContiguousParams{}})
	t.Fatalf("Gribble did not panic when given a command with two parameters "+
		"numbered '1' and '3' consecutively.")
}

type badType struct {
	Op uint `param:"1"`
}

func (cmd badType) Run() Value { return nil }

// TestBadType tests whether Gribble catches an invalid concrete type.
func TestBadType(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badType{}})
	t.Fatalf("Gribble did not panic when given a parameter with type 'uint'.")
}

type badInterface struct {
	Op interface{} `param:"1" types:"int,float"`
}

func (cmd badInterface) Run() Value { return nil }

// TestBadInterface tests whether Gribble catches an invalid interface type.
// (The only interface type allowed for parameters is 'Any'.)
func TestBadInterface(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badInterface{}})
	t.Fatalf("Gribble did not panic when given a parameter with type "+
		"'interface{}'.")
}

type badNumberOfTypes struct {
	Op Any `param:"1" types:"int"`
}

func (cmd badNumberOfTypes) Run() Value { return nil }

// TestBadNumberOfTypes tests whether Gribble catches an invalid number of types
// specified using the 'types' struct tag. Namely, there must be at least two
// types specified (otherwise, a regular concrete type should be used).
func TestBadNumberOfTypes(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badNumberOfTypes{}})
	t.Fatalf("Gribble did not panic when given a parameter with type 'Any' "+
		"that has only specified a single type in the 'types' struct tag.")
}

type badAnyTypes struct {
	Op Any `param:"1" types:"int,uint"`
}

func (cmd badAnyTypes) Run() Value { return nil }

// TestBadAnyTypes tests whether Gribble catches an invalid type specified
// in the 'types' struct tag.
func TestBadAnyTypes(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badAnyTypes{}})
	t.Fatalf("Gribble did not panic when given a parameter with type 'Any' "+
		"that has an invalid type 'uint' in the 'types' struct tag.")
}

type badRepeatedTypes struct {
	Op Any `param:"1" types:"int,int"`
}

func (cmd badRepeatedTypes) Run() Value { return nil }

// TestBadRepeatedTypes tests whether Gribble catches an invalid type specified
// in the 'types' struct tag.
func TestBadRepeatedTypes(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badRepeatedTypes{}})
	t.Fatalf("Gribble did not panic when given a parameter with type 'Any' "+
		"that has repeated the 'int' type in the 'types' struct tag.")
}

type badReturnValue struct {
	Op int `param:"1"`
}

func (cmd badReturnValue) Run() Value { return false }

// TestBadReturnValue tests whether Gribble panics when a sub-command returns
// a Value whose concrete type is not an int, float64 or a string.
func TestBadReturnValue(t *testing.T) {
	defer func() {
		recover()
	}()
	New([]Command{badReturnValue{}}).Run("badReturnValue (badReturnValue 0)")
	t.Fatalf("Gribble did not panic when a sub-command returned a value "+
		"with type 'bool', which is not one of the allowed return types: "+
		"int, float64 or string.")
}

