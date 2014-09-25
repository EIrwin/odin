package cli

import "fmt"
import "time"
import "strings"
import "bytes"

type flagable struct {
	*writer
	flags           flagMap
	aliases         map[rune]*Flag
	flagValues      map[*Flag]Value
	flagsTerminated bool
	parsed          bool
	version         string
}

func (cmd *flagable) AliasFlag(alias rune, flagname string) {
	flag, ok := cmd.flags[flagname]
	if !ok {
		panic(fmt.Sprintf("flag not defined %v", flagname))
	}
	if cmd.aliases == nil {
		cmd.aliases = make(map[rune]*Flag)
	}
	cmd.aliases[alias] = flag
}

// Bool defines a bool flag with specified name, default value, and usage string.
// The return value is the address of a bool variable that stores the value of the flag.
func (cmd *flagable) DefineBoolFlag(name string, value bool, usage string) *bool {
	p := new(bool)
	cmd.DefineBoolFlagVar(p, name, value, usage)
	return p
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func (cmd *flagable) DefineBoolFlagVar(p *bool, name string, value bool, usage string) {
	cmd.DefineFlag(newBoolValue(value, p), name, usage)
}

// Duration defines a time.Duration flag with specified name, default value, and usage string.
// The return value is the address of a time.Duration variable that stores the value of the flag.
func (cmd *flagable) DefineDurationFlag(name string, value time.Duration, usage string) *time.Duration {
	p := new(time.Duration)
	cmd.DefineDurationFlagVar(p, name, value, usage)
	return p
}

// DurationVar defines a time.Duration flag with specified name, default value, and usage string.
// The argument p points to a time.Duration variable in which to store the value of the flag.
func (cmd *flagable) DefineDurationFlagVar(p *time.Duration, name string, value time.Duration, usage string) {
	cmd.DefineFlag(newDurationValue(value, p), name, usage)
}

// Float64 defines a float64 flag with specified name, default value, and usage string.
// The return value is the address of a float64 variable that stores the value of the flag.
func (cmd *flagable) DefineFloat64Flag(name string, value float64, usage string) *float64 {
	p := new(float64)
	cmd.DefineFloat64FlagVar(p, name, value, usage)
	return p
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func (cmd *flagable) DefineFloat64FlagVar(p *float64, name string, value float64, usage string) {
	cmd.DefineFlag(newFloat64Value(value, p), name, usage)
}

// Int64 defines an int64 flag with specified name, default value, and usage string.
// The return value is the address of an int64 variable that stores the value of the flag.
func (cmd *flagable) DefineInt64Flag(name string, value int64, usage string) *int64 {
	p := new(int64)
	cmd.DefineInt64FlagVar(p, name, value, usage)
	return p
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func (cmd *flagable) DefineInt64FlagVar(p *int64, name string, value int64, usage string) {
	cmd.DefineFlag(newInt64Value(value, p), name, usage)
}

// Int defines an int flag with specified name, default value, and usage string.
// The return value is the address of an int variable that stores the value of the flag.
func (cmd *flagable) DefineIntFlag(name string, value int, usage string) *int {
	p := new(int)
	cmd.DefineIntFlagVar(p, name, value, usage)
	return p
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func (cmd *flagable) DefineIntFlagVar(p *int, name string, value int, usage string) {
	cmd.DefineFlag(newIntValue(value, p), name, usage)
}

// String defines a string flag with specified name, default value, and usage string.
// The return value is the address of a string variable that stores the value of the flag.
func (cmd *flagable) DefineStringFlag(name string, value string, usage string) *string {
	p := new(string)
	cmd.DefineStringFlagVar(p, name, value, usage)
	return p
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func (cmd *flagable) DefineStringFlagVar(p *string, name string, value string, usage string) {
	cmd.DefineFlag(newStringValue(value, p), name, usage)
}

// Uint64 defines a uint64 flag with specified name, default value, and usage string.
// The return value is the address of a uint64 variable that stores the value of the flag.
func (cmd *flagable) DefineUint64Flag(name string, value uint64, usage string) *uint64 {
	p := new(uint64)
	cmd.DefineUint64FlagVar(p, name, value, usage)
	return p
}

// Uint64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the flag.
func (cmd *flagable) DefineUint64FlagVar(p *uint64, name string, value uint64, usage string) {
	cmd.DefineFlag(newUint64Value(value, p), name, usage)
}

// Uint defines a uint flag with specified name, default value, and usage string.
// The return value is the address of a uint  variable that stores the value of the flag.
func (cmd *flagable) DefineUintFlag(name string, value uint, usage string) *uint {
	p := new(uint)
	cmd.DefineUintFlagVar(p, name, value, usage)
	return p
}

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func (cmd *flagable) DefineUintFlagVar(p *uint, name string, value uint, usage string) {
	cmd.DefineFlag(newUintValue(value, p), name, usage)
}

// DefineFlag defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type Value, which
// typically holds a user-defined implementation of Value. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of Value; in particular, Set would
// decompose the comma-separated string into the slice.
func (cmd *flagable) DefineFlag(value Value, name string, usage string) {
	// Remember the default value as a string; it won't change.
	flag := &Flag{
		Name:     name,
		Usage:    usage,
		value:    value,
		DefValue: value.String(),
	}
	_, alreadythere := cmd.flags[name]
	if alreadythere {
		cmd.panicf("flag redefined: %s", name)
	}
	if cmd.flags == nil {
		cmd.flags = make(map[string]*Flag)
	}
	cmd.flags[name] = flag
}

// Flag returns the Value interface to the value of the named flag,
// returning nil if none exists.
func (cmd *flagable) Flag(name string) Value {
	flag, ok := cmd.flags[name]
	if !ok {
		panic(fmt.Sprintf("flag not defined %v", name))
	}
	value, ok := cmd.flagValues[flag]
	if ok {
		return value
	}
	return nil
}

func (cmd *flagable) Flags() map[string]Value {
	flags := make(map[string]Value)
	for name := range cmd.flags {
		flags[name] = cmd.Flag(name)
	}
	return flags
}

// FlagCount returns the number of flags that have been set.
func (cmd *flagable) FlagCount() int { return len(cmd.flagValues) }

// Parsed returns if the flags have been parsed
func (cmd *flagable) Parsed() bool {
	return cmd.parsed
}

// UsageString returns the flags usage as a string
func (cmd *flagable) UsageString() string {
	var maxBufferLen int
	flagsUsages := make(map[*Flag]*bytes.Buffer)

	// init the map for each flag
	for _, flag := range cmd.aliases {
		flagsUsages[flag] = bytes.NewBufferString("")
	}

	// Get each flags aliases
	for r, flag := range cmd.aliases {
		alias := string(r)
		buffer := flagsUsages[flag]
		var err error
		if buffer.Len() == 0 {
			_, err = buffer.WriteString(fmt.Sprintf("-%s", alias))
		} else {
			_, err = buffer.WriteString(fmt.Sprintf(", -%s", alias))
		}
		exitIfError(err)
		buffLen := len(buffer.String())
		if buffLen > maxBufferLen {
			maxBufferLen = buffLen
		}
	}

	// Get each flags names
	for name, flag := range cmd.flags {
		buffer := flagsUsages[flag]
		if buffer == nil {
			flagsUsages[flag] = new(bytes.Buffer)
			buffer = flagsUsages[flag]
		}
		var err error
		if buffer.Len() == 0 {
			_, err = buffer.WriteString(fmt.Sprintf("--%s", name))
		} else {
			_, err = buffer.WriteString(fmt.Sprintf(", --%s", name))
		}
		if _, ok := flag.value.(boolFlag); !ok {
			buffer.WriteString(fmt.Sprintf("=\"%s\"", flag.DefValue))
		}
		exitIfError(err)
		buffLen := len(buffer.String())
		if buffLen > maxBufferLen {
			maxBufferLen = buffLen
		}
	}

	// get the flag strings and append the usage info
	var outputLines []string
	for i := 0; i < len(cmd.flags); i++ {
		flag := cmd.flags.Sort()[i]
		buffer := flagsUsages[flag]
		for {
			buffLen := len(buffer.String())
			if buffLen > maxBufferLen {
				break
			}
			buffer.WriteString(" ")
		}
		outputLines = append(outputLines, fmt.Sprintf("  %s # %s", buffer.String(), flag.Usage))
	}

	return strings.Join(outputLines, "\n")
}

func (cmd *flagable) Version() string {
	return cmd.version
}

// defineHelp defines a help function and alias if they are not present
func (cmd *flagable) defineHelp() {
	if _, ok := cmd.flags["help"]; !ok {
		cmd.DefineBoolFlag("help", false, "show help and exit")
		if _, ok := cmd.aliases['h']; !ok {
			cmd.AliasFlag('h', "help")
		}
	}
}

// defineVersion defines a version if one has been set
func (cmd *flagable) defineVersion() {
	if _, ok := cmd.flags["version"]; !ok && len(cmd.Version()) > 0 {
		cmd.DefineBoolFlag("version", false, "show version and exit")
		if _, ok := cmd.aliases['v']; !ok {
			cmd.AliasFlag('v', "version")
		}
	}
}

// flagFromArg determines the flags from an argument
func (cmd *flagable) flagFromArg(arg string) (bool, []*Flag) {
	var flags []*Flag

	// Do nothing if flags terminated
	if cmd.flagsTerminated {
		return false, flags
	}
	if arg[len(arg)-1] == '=' {
		cmd.errf("invalid flag format")
	}
	arg = strings.Split(arg, "=")[0]

	// Determine if we need to terminate flags
	isFlag := arg[0] == '-'
	areAliases := isFlag && arg[1] != '-'
	isTerminator := !areAliases && len(arg) == 2

	if !isFlag || isTerminator {
		cmd.flagsTerminated = true
		return false, flags
	}

	// Determine if name or alias
	if areAliases {
		aliases := arg[1:]
		for _, c := range aliases {
			flag, ok := cmd.aliases[c]
			if !ok {
				cmd.errf("invalid alias: %v", string(c))
			}
			flags = append(flags, flag)
		}
	} else {
		name := arg[2:]
		flag, ok := cmd.flags[name]
		if !ok {
			cmd.errf("invalid flag")
		}
		flags = append(flags, flag)
	}
	return areAliases, flags
}

// parseFlags flag definitions from the argument list, returns any left over
// arguments after flags have been parsed.
func (cmd *flagable) parse(args []string) []string {
	cmd.defineHelp()
	cmd.defineVersion()
	cmd.parsed = true
	i := 0
	for i < len(args) {
		isAlias, flags := cmd.flagFromArg(args[i])
		if cmd.flagsTerminated {
			break
		}
		if isAlias {
			cmd.setAliasValues(flags, args[i])
		} else {
			cmd.setFlagValue(flags[0], args[i])
		}
		i++
	}
	// Set the remaining flags to defaults
	cmd.setFlagDefaults()
	// return the remaining unused args
	return args[i:]
}

// setAliasValues sets the values of flags from thier aliases
func (cmd *flagable) setAliasValues(flags []*Flag, arg string) {
	args := strings.Split(arg, "=")
	hasvalue := len(args) > 1
	var lastflag *Flag

	// If a value is provided, set the last flag
	if hasvalue {
		lastflag = flags[len(flags)-1]
		flags = flags[:len(flags)-1]
		cmd.setFlag(lastflag, args[1])
	}

	for i := 0; i < len(flags); i++ {
		flag := flags[i]
		if fv, ok := flag.value.(boolFlag); ok && fv.IsBoolFlag() {
			cmd.setFlag(flag, "true")
		} else {
			cmd.errf("flag \"--%v\" is missing a value", flag.Name)
		}
	}
}

// setFlagDefaults sets the default values of all flags
func (cmd *flagable) setFlagDefaults() {
	for name, flag := range cmd.flags {
		if cmd.Flag(name) == nil {
			cmd.setFlag(flag, flag.DefValue)
		}
	}
}

// setFlag sets the value of the named flag.
func (cmd *flagable) setFlag(flag *Flag, value string) error {
	// Verify the flag is a flag for f set
	flag, ok := cmd.flags[flag.Name]
	if !ok {
		return fmt.Errorf("no such flag -%v", flag.Name)
	}
	err := flag.value.Set(value)
	if err != nil {
		return err
	}
	if cmd.flagValues == nil {
		cmd.flagValues = make(map[*Flag]Value)
	}
	cmd.flagValues[flag] = flag.value
	return nil
}

// setFlagValue sets the value of a given flag
func (cmd *flagable) setFlagValue(flag *Flag, arg string) {
	args := strings.Split(arg, "=")
	hasvalue := len(args) > 1
	if hasvalue {
		value := args[1]
		cmd.setFlag(flag, value)
	} else {
		if fv, ok := flag.value.(boolFlag); ok && fv.IsBoolFlag() {
			cmd.setFlag(flag, "true")
		} else {
			cmd.errf("flag \"--%v\" is missing a value", flag.Name)
		}
	}
}
