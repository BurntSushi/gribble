/*
Package gribble provides a command oriented language whose environment of 
commands is defined by Go structs via reflection. The primary use case for 
Gribble is to provide an easy to use command language for users to interact 
with your program.

For example, to define an environment with an "Add" command that adds two 
operators and print the output of the command "Add (Add 1 2) 42":

	package main

	import (
		"fmt"
		"github.com/BurntSushi/gribble"
	)

	type Add struct {
		Op1 int `param:"1"`
		Op2 int `param:"2"`
	}

	func (cmd Add) Run() gribble.Value {
		return cmd.Op1 + cmd.Op2
	}

	func main() {
		env := gribble.New([]gribble.Command{Add{}})
		val, _ := env.Run("Add (Add 1 2) 42")
		fmt.Println(val)
	}

How to define commands

Commands that make up a Gribble environment are struct values that implement
the 'Command' interface. A minimal command with no parameters could be
defined as follows:

	type Minimal struct {}
	func (cmd Minimal) Run() gribble.Value {
		return "Hello"
	}

Parameters are added by defining a new struct field to hold that parameter's
value, and labeling it with a parameter number using struct tags:

	type Negation struct {
		Op int `param:"1"`
	}

Struct tags MUST be used to define parameters. For example, the following
command has ZERO parameters:

	type ZeroParameters struct {
		something string
	}

Parameters may also accept more than one type of argument by defining the
parameter struct field to have type 'gribble.Any' and specifying the allowable 
types with the 'types' struct tag:

	type Negation struct {
		Op gribble.Any `param:"1" types:"int,float"`
	}

And the value of the parameter can be safely type switched on only the types
specified when running the command:

	func (cmd Negation) Run() gribble.Value {
		switch val := cmd.Op.(type) {
		case int:
			return -val
		case float64:
			return -val
		}
		panic("unreachable")
	}

Gribble enforces the invariant that 'Op' must only contain values with concrete 
type 'int' or 'float64' by returning an error if a type not specified in the 
'types' struct tag is found.

Currently, the only types allowed in the 'types' struct tag are int, float and 
string.

Errors

Gribble will panic when any of the invariants involved in defining commands are 
broken. This includes, but is not limited to: Specifying non-contiguous 
parameter numbers. Using a type for a parameter other than int, float64, string 
or gribble.Any. Specifying a type other than int, float or string in the 
'types' struct tag. Using the gribble.Any type without the 'types' struct tag 
or with a 'types' struct tag with only one type specified. Using a sub-command 
that returns a concrete value other than an int, float64 or string.

Other errors like parse errors or run-time type errors are returned as 
standard Go errors when using the 'Run' or 'Command' methods.

Gribble EBNF

What follows is an EBNF grammar of the Gribble language. (Bug reports are 
welcome!)

Names of the format 'go-NAME' refer to lexical elements called NAME in the
Go Programming Language Specification.

	program = command ;

	command = [ "(" ], go-identifier, { param }, [ ")" ] ;

	param = go-string_lit | go-int_lit | go-float_lit | "(" command ")" ;

Installation

Gribble is go gettable:

	go get github.com/BurntSushi/gribble

Quick example

To demonstrate the bundled integer calculator:

	go get github
	GO/PATH/bin/int-calc 'add 5 (mul 2 6)'

The output should be '17'.

Why

The initial motivation for Gribble was for interacting with my window manager, 
Wingo (which is where the name Gribble came from). Wingo initially had a hacked 
together set of commands, but it quickly grew out of control and incredibly 
difficult to maintain. Having a language that is definable using just Go 
structs makes everything a lot cleaner.

*/
package gribble
