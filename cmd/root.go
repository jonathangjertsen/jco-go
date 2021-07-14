package cmd

import (
	"fmt"
	"github.com/jonathangjertsen/jco-go/ops"
	"github.com/jonathangjertsen/jco-go/table"
	"github.com/spf13/cobra"
	"math/big"
	"strconv"
)

const (
	version = "v0.1.0"
)

type Flags struct {
	bits             uint
	help             bool
	version          bool
	numbers          [][]byte
	numbersAsWritten []string
}

var RootCmd = &cobra.Command{
	Use:   "jco",
	Short: fmt.Sprintf("jco %s", version),
	Long:  fmt.Sprintf("jco (Jonathan's converter) %s", version),
	Run: func(cmd *cobra.Command, args []string) {
		flags := parseFlags(args)
		if flags.version {
			fmt.Printf(cmd.Short)
			return
		}
		if flags.help {
			cmd.Usage()
			return
		}
		t := table.NewTable(flags.bits)

		switch len(flags.numbers) {
		case 0:
			cmd.Usage()
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
	},
	Version:           version,
	DisableAutoGenTag: true,
}

func Fatal(message string) {
	panic(message)
}

func Interactive() {
	fmt.Printf("Interactive")
}

/**
 * Custom flag parsing because Cobra does not want to support negative numbers:
 * 	 https://github.com/spf13/cobra/issues/124
 */
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

func init() {
	RootCmd.Flags().UintP("bits", "b", 32, "Number of bits")
}

func Execute() {
	cobra.MousetrapHelpText = ""
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.DisableFlagParsing = true
	if err := RootCmd.Execute(); err != nil {
		Fatal(err.Error())
	}
}
