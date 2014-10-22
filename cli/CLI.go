package cli

import (
	"io"
	"os"
)

import "strings"
import "github.com/jwaldrip/odin/cli/values"

// ShowUsage is a function to show usage
var ShowUsage = func(cmd Command) { cmd.Usage() }

// CLI represents a set of defined flags.  The zero value of a FlagSet
// has no name and has ContinueOnError error handling.
type CLI struct {
	ErrorHandling ErrorHandling

	aliases          map[rune]*Flag
	description      string
	errOutput        io.Writer
	flags            flagMap
	flagsTerminated  bool
	flagValues       map[*Flag]values.Value
	fn               func(Command)
	inheritedFlags   flagMap
	name             string
	params           paramsList
	paramValues      map[*Param]values.Value
	paramsParsed     bool
	parent           Command
	propogatingFlags flagMap
	stdOutput        io.Writer
	subCommands      map[string]*SubCommand
	unparsedArgs     values.List
	usage            func()
	version          string
}

func (cmd *CLI) init(name, desc string, fn func(Command), paramNames ...string) {
	cmd.name = name
	cmd.fn = fn
	cmd.description = desc
	cmd.DefineParams(paramNames...)
	cmd.ErrorHandling = ExitOnError
}

// New returns a new cli with the specified name and
// error handling property.
func New(version, desc string, fn func(Command), paramNames ...string) *CLI {
	nameParts := strings.Split(os.Args[0], "/")
	cli := new(CLI)
	cli.init(nameParts[len(nameParts)-1], desc, fn, paramNames...)
	cli.version = version
	cli.description = desc
	return cli
}

// Alias for New
var NewCLI = New

// Start starts the command with args, arg[0] is ignored
func (cmd *CLI) Start(args ...string) {
	if args == nil {
		args = os.Args
	}

	if len(args) > 1 {
		args = args[1:]
	} else {
		args = []string{}
	}

	// parse flags and params
	args = cmd.parse(args)

	var subCommandsParsed bool
	if args, subCommandsParsed = cmd.parseSubCommands(args); subCommandsParsed {
		return
	}

	cmd.assignUnparsedArgs(args)

	// Run the function
	cmd.fn(cmd)
}
