// Example int-calc demonstrates a simple integer calculator with
// commands defined using Gribble. It supports basic arithmetic operations
// on integers.
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
// Note that these may also be specified as pointers.
var env *gribble.Environment = gribble.New([]gribble.Command{
	Add{},
	Subtract{},
	Multiply{},
	Divide{},
})

// Add implements the arithmetic '+' operation on integers.
type Add struct {
	// We can override the name of the command to use 'add' instead of 'Add'.
	// If the 'name' field is absent, then the name of the command is the same
	// as the name of type. (In this case, 'Add'.)
	name string `add`
	Op1  int    `param:"1"`
	Op2  int    `param:"2"`
}

// Run simply adds Op1 and Op2.
func (c Add) Run() gribble.Value {
	return c.Op1 + c.Op2
}

type Subtract struct {
	name string `sub`
	Op1  int    `param:"1"`
	Op2  int    `param:"2"`
}

func (c Subtract) Run() gribble.Value {
	return c.Op1 - c.Op2
}

type Multiply struct {
	name string `mul`
	Op1  int    `param:"1"`
	Op2  int    `param:"2"`
}

func (c Multiply) Run() gribble.Value {
	return c.Op1 * c.Op2
}

type Divide struct {
	name string `div`
	Op1  int    `param:"1"`
	Op2  int    `param:"2"`
}

func (c Divide) Run() gribble.Value {
	return c.Op1 / c.Op2
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
