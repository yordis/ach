// Copyright 2019 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/moov-io/ach"
)

var (
	flagVerbose = flag.Bool("v", false, "Print verbose details about each ACH file")
	flagVersion = flag.Bool("version", false, "Print moov-io/ach cli version")

	flagDiff     = flag.Bool("diff", false, "Compare two files against each other")
	flagFlatten  = flag.Bool("flatten", false, "Flatten batches in each file")
	flagMerge    = flag.Bool("merge", false, "Merge files before describing")
	flagReformat = flag.String("reformat", "", "Reformat an incoming ACH file to another format")

	flagMask              = flag.Bool("mask", false, "Mask/hide full account numbers and individual names")
	flagMaskAccounts      = flag.Bool("mask.accounts", false, "Mask/hide full account numbers")
	flagMaskCorrectedData = flag.Bool("mask.corrections", false, "Mask/Hide Corrected Data in Addenda98 records")
	flagMaskNames         = flag.Bool("mask.names", false, "Mask/hide full individual names")

	flagPretty        = flag.Bool("pretty", false, "Display all values in their human readable format")
	flagPrettyAmounts = flag.Bool("pretty.amounts", false, "Display human readable amounts instead of exact values")
)

func main() {
	flag.Usage = help
	flag.Parse()

	switch {
	case *flagVersion:
		fmt.Printf("moov-io/ach:%s cli tool\n", ach.Version)
		return
	case *flagVerbose:
		fmt.Printf("moov-io/ach:%s cli tool\n", ach.Version)
	}

	args := flag.Args()

	// error conditions, verify we're okay for whatever the task at hand is
	switch {
	case *flagDiff && len(args) != 2:
		fmt.Printf("with -diff exactly two files are expected, found %d files\n", len(args))
		os.Exit(1)
	}

	// minor debugging
	if *flagVerbose {
		fmt.Printf("found %d ACH files to describe: %s\n", len(args), strings.Join(args, ", "))
	}

	// pick our command to do
	switch {
	case *flagDiff:
		if err := diffFiles(args); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

	case *flagReformat != "" && len(args) == 1:
		if err := reformat(*flagReformat, args[0]); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}

	default:
		if err := dumpFiles(args); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			os.Exit(1)
		}
	}
}
