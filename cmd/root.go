package cmd

import (
	"fmt"
	"github.com/jonathangjertsen/jco-go/ops"
	"github.com/jonathangjertsen/jco-go/table"
	"math/big"
	"os"
	"strconv"
)

const (
	VERSION = "v0.1.0"
)

type Flags struct {
	bits             uint
	help             bool
	version          bool
	numbers          [][]byte
	numbersAsWritten []string
}

func parseFlags(args []string) *Flags {
	flags := Flags{}
	opts := map[string]string{
		"-b": "32",
	}
	currentOpt := ""
	for _, arg := range args {
		if currentOpt == "" {
			num, ok := new(big.Int).SetString(arg, 0)
			if ok {
				bytes := num.Bytes()
				flags.numbers = append(flags.numbers, bytes)
				flags.numbersAsWritten = append(flags.numbersAsWritten, arg)
			} else {
				switch arg {
				case "-v", "--version":
					flags.version = true
				case "-h", "--help":
					flags.help = true
				default:
					currentOpt = arg
				}
			}
		} else {
			opts[currentOpt] = arg
			currentOpt = ""
		}
	}

	// Extracts 'bits' argument
	bitsU64, err := strconv.ParseUint(opts["-b"], 0, 6)
	bits := uint(bitsU64)
	if err != nil || bits < 1 || bits > 64 {
		Fatal(fmt.Sprintf("Invalid value for -b: %s", opts["-b"]))
	}
	if bits%8 != 0 {
		bitsRounded := (bits + 7) / 8
		fmt.Printf("Warning: -b %v will be rounded up to %v", bits, bitsRounded)
		bits = bitsRounded
	}
	flags.bits = bits

	// Pad numbers up to bytes
	for i, num := range flags.numbers {
		nBytes := flags.bits / 8
		nBytesInNum := uint(len(num))
		if nBytesInNum < nBytes {
			flags.numbers[i] = ops.PrependZeros(num, uint(nBytes-nBytesInNum))
		}
	}

	return &flags
}

func Version() {
	fmt.Printf("jco %s", VERSION)
}

func Usage() {
	fmt.Printf(`jco (Jonathan's converter) %s

Usage:
	Show information about <number>
		jco <number>

	Show information about how <number1> and <number2> relate
		jco <number1> <number2>

	Like the above, but treat numbers as 16-bit
		jco <number1> <number2> -b 16

	Show this help screen
		jco --help

	Show one-liner version
		jco --version
`, VERSION)
}

func Execute() {
	args := os.Args[1:]
	flags := parseFlags(args)
	if flags.version {
		Version()
		return
	}
	if flags.help {
		Usage()
		return
	}
	t := table.NewTable(flags.bits)

	switch len(flags.numbers) {
	case 0:
		Usage()
		return
	case 1:
		t.One(
			flags.numbers[0],
			flags.numbersAsWritten[0],
		)
	case 2:
		t.Two(
			flags.numbers[0],
			flags.numbers[1],
			flags.numbersAsWritten[0],
			flags.numbersAsWritten[1],
		)
	default:
	}
	t.Render()
}

func Fatal(message string) {
	panic(message)
}

func Interactive() {
	fmt.Printf("Interactive")
}
