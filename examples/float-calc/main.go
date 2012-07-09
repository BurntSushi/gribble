// Example float-calc demonstrates an extremely simple usage scenario for
// Gribble.Namely, float-calc is a floating point calculator that supports only 
// basic arithmetic operations on both integers and floating point numbers.
//
// The key point of interest in this example is to show how parameters can
// allow values of more than one type.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/gribble"
)

// A Gribble environment is composed of Go struct types. Since the environment
// does not care about values, zero values of each command struct can be
// used in the gribble.Command slice.
//
// Note that these may also be specified as non-pointers.
var env *gribble.Environment = gribble.New([]gribble.Command{
	&Add{},
	&Subtract{},
	&Multiply{},
	&Divide{},
})

// Add implements the arithmetic '+' operation on decimals.
type Add struct {
	// We can override the name of the command to use 'add' instead of 'Add'.
	// If the 'name' field is absent, then the name of the command is the same
	// as the name of type. (In this case, 'Add'.)
	name string `add`

	// In the 'int-calc' example, the operators are defined using a conrete
	// type like 'Op1 int `param:"1"`'. In this case, we want to allow the
	// operators to be either floats or integers. To accomplish this, we
	// make the operators have type gribble.Any (which is just an empty
	// interface), and specify the allowable types using the 'types' struct
	// tag.
	//
	// The types available to you are int, float or string.
	Op1 gribble.Any `param:"1" types:"int,float"`
	Op2 gribble.Any `param:"2" types:"int,float"`
}

// Run converts each operator to a float, and adds them.
func (c *Add) Run() gribble.Value {
	return float(c.Op1) + float(c.Op2)
}

type Subtract struct {
	name string      `sub`
	Op1  gribble.Any `param:"1" types:"int,float"`
	Op2  gribble.Any `param:"2" types:"int,float"`
}

func (c *Subtract) Run() gribble.Value {
	return float(c.Op1) - float(c.Op2)
}

type Multiply struct {
	name string      `mul`
	Op1  gribble.Any `param:"1" types:"int,float"`
	Op2  gribble.Any `param:"2" types:"int,float"`
}

func (c *Multiply) Run() gribble.Value {
	return float(c.Op1) * float(c.Op2)
}

type Divide struct {
	name string      `div`
	Op1  gribble.Any `param:"1" types:"int,float"`
	Op2  gribble.Any `param:"2" types:"int,float"`
}

func (c *Divide) Run() gribble.Value {
	return float(c.Op1) / float(c.Op2)
}

// float is a convenience method that accepts a parameter value of type
// gribble.Any, and returns a float64 of that value. This is accomplished
// with a type switch and a type conversion if the underlying value is an
// integer.
func float(val gribble.Any) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	// This is unreachable because the values allowed in each of 'Op1' and 'Op2'
	// are allowed to be 'int' or 'float'. This invariant is enforced by
	// Gribble. That is, if this panic is hit, Gribble has a bug.
	panic("unreachable")
}

// usage overrides the default flag.Usage to output a list of all available
// commands and each command's parameter type list.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s command\n", path.Base(os.Args[0]))
	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr, "\nAvailable commands:")
	fmt.Fprintln(os.Stderr, env.StringTypes())
	os.Exit(1)
}

// main uses all of flag.Args as a command, and prints either the output
// of the command or an error if one occurs.
func main() {
	flag.Usage = usage
	flag.Parse()
	cmd := strings.Join(flag.Args(), " ")

	val, err := env.Run(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		if name := env.CommandName(cmd); len(name) > 0 {
			fmt.Printf("Usage: %s\n", env.UsageTypes(name))
		}
		return
	}
	fmt.Println(val)
}
