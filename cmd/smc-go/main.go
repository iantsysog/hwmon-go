package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	smc "github.com/iantsysog/smc-go"
)

type options struct {
	key      string
	valueHex string
}

func parseArgs(args []string) (options, error) {
	var opts options
	fs := flag.NewFlagSet("smc-go", flag.ContinueOnError)
	fs.StringVar(&opts.key, "k", "", "SMC key to read or write. Output is printed as hex.")
	fs.StringVar(&opts.valueHex, "v", "", "Hex value to write. If omitted, the key is read.")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	if opts.key == "" {
		return options{}, fmt.Errorf("missing required flag: -k (SMC key)")
	}
	if len(opts.key) != 4 {
		return options{}, fmt.Errorf("invalid -k value: SMC key must be exactly 4 characters")
	}
	return opts, nil
}

func run(opts options, c smc.Connection, out io.Writer) (err error) {
	if err := c.Open(); err != nil {
		return err
	}
	defer func() {
		if closeErr := c.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if opts.valueHex == "" {
		v, err := c.Read(opts.key)
		if err != nil {
			return err
		}
		fmt.Fprintln(out, hex.EncodeToString(v.Bytes))
		if v.DataType == "flt " {
			if f, err := v.Float32LE(); err == nil {
				fmt.Fprintln(out, f)
			}
		}
		return nil
	}

	b, err := hex.DecodeString(opts.valueHex)
	if err != nil {
		return err
	}

	if err := c.Write(opts.key, b); err != nil {
		return err
	}
	return nil
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	if err := run(opts, smc.New(), log.Writer()); err != nil {
		log.Fatal(err)
	}
}
