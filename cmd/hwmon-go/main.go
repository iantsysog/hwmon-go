package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	_ "github.com/iantsysog/hwmon-go/backends/hid"
	_ "github.com/iantsysog/hwmon-go/backends/smc"
	"github.com/iantsysog/hwmon-go/hwmon"
	"github.com/iantsysog/hwmon-go/internal/utils"
)

type options struct {
	backendsCSV string
	json        bool
}

func parseArgs(args []string) (options, error) {
	var opts options
	fs := flag.NewFlagSet("hwmon-go", flag.ContinueOnError)
	fs.StringVar(&opts.backendsCSV, "backends", "", "Comma-separated list of backends to enable (e.g. smc,hid). If omitted, all registered backends are used.")
	fs.BoolVar(&opts.json, "json", false, "Output as JSON.")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	return opts, nil
}

func run(opts options) error {
	var copts []hwmon.Option
	if strings.TrimSpace(opts.backendsCSV) != "" {
		seen := make(map[string]struct{}, 8)
		var names []string
		for p := range strings.SplitSeq(opts.backendsCSV, ",") {
			n := strings.TrimSpace(p)
			if n == "" {
				continue
			}
			if _, ok := seen[n]; ok {
				continue
			}
			seen[n] = struct{}{}
			names = append(names, n)
		}
		if len(names) > 0 {
			copts = append(copts, hwmon.WithBackends(names...))
		}
	}

	rs, warn := hwmon.Collect(context.Background(), copts...)
	if opts.json {
		return writeJSON(os.Stdout, rs, warn)
	}
	return writeTable(os.Stdout, rs, warn)
}

func main() {
	opts, err := parseArgs(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		log.Fatal(err)
	}
	if err := run(opts); err != nil {
		log.Fatal(err)
	}
}

func writeTable(w io.Writer, rs []hwmon.Reading, warn error) error {
	if w == nil {
		return fmt.Errorf("nil writer")
	}

	tw := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	for _, r := range rs {
		value := formatValue(r.Value)
		if r.Unit != "" && value != "" {
			value += " " + r.Unit
		} else if r.Unit != "" && value == "" {
			value = r.Unit
		}

		cols := []string{string(r.Kind), r.Name, value, r.Source, r.KeyOrID}
		if r.DataType != "" {
			cols = append(cols, r.DataType)
		}
		utils.WriteTSVRow(tw, cols...)
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	utils.PrintWarning(warn)
	return nil
}

func writeJSON(w io.Writer, rs []hwmon.Reading, warn error) error {
	if err := utils.WriteReadingsJSON(w, rs); err != nil {
		return err
	}
	utils.PrintWarning(warn)
	return nil
}

func formatValue(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case float64:
		return strconv.FormatFloat(x, 'g', 6, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'g', 6, 64)
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case bool:
		return strconv.FormatBool(x)
	case string:
		return x
	default:
		return fmt.Sprint(x)
	}
}
