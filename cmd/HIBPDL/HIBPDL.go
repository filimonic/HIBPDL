//go:build ntlm || !ntlm

package main

import (
	"flag"
	"fmt"
	"hibpdl/internal/hibpdownloader"
	"log"
	"runtime"
)

const fileSha1Default = "pwnedpasswords.sha1.txt"
const fileNtlmDefault = "pwnedpasswords.ntlm.txt"

type Config struct {
	ntlm        bool
	overwrite   bool
	singleFile  bool
	parallelism uint64
	progress    bool
	file        string
}

var config Config

func init() {

	var numCpu = uint64(runtime.NumCPU())

	const ntlmDefault = false
	const overwriteDefault = false
	const noProgressDefault = false
	const singleFileDefault = false

	var parallelismDefault = numCpu
	const parallelismUsage = "number of parallel queries"

	var noProgress bool

	flag.BoolVar(&config.ntlm, "ntlm", ntlmDefault, "")
	flag.BoolVar(&config.ntlm, "n", ntlmDefault, "")

	flag.BoolVar(&config.overwrite, "overwrite", singleFileDefault, "")
	flag.BoolVar(&config.overwrite, "o", singleFileDefault, "")

	flag.BoolVar(&config.singleFile, "single", overwriteDefault, "")
	flag.BoolVar(&config.singleFile, "s", overwriteDefault, "")

	flag.BoolVar(&noProgress, "no-progress", noProgressDefault, "")
	flag.BoolVar(&noProgress, "q", noProgressDefault, "")

	flag.Uint64Var(&config.parallelism, "parallelism", parallelismDefault, parallelismUsage)
	flag.Uint64Var(&config.parallelism, "p", parallelismDefault, parallelismUsage+" (shorthand)")

	flag.Usage = usage

	header()
	flag.Parse()

	if config.parallelism < 1 {
		log.Fatalln("parallelism must be at least one")
	}

	if flag.NArg() == 1 {
		config.file = flag.Args()[0]
	} else if flag.NArg() == 0 {
		if config.ntlm {
			config.file = fileNtlmDefault
		} else {
			config.file = fileSha1Default
		}

	}

	config.progress = !noProgress
}

func main() {
	hibpdownloader.Download(config.file, config.parallelism, config.overwrite, config.ntlm, config.progress)
}

func usage() {

	fmt.Printf("Usage: HIBPDL [options] [file]\n")
	fmt.Printf("\n")

	fmt.Printf("Options:\n")

	fmt.Printf("  -n, --ntlm                 Download NTLM hashes instead of SHA1.\n")
	fmt.Printf("                               Default: download SHA1 \n")
	fmt.Printf("  -o, --overwrite            Overwrite output file if exists\n")
	fmt.Printf("  -q, --no-progress          Do not output progress bar.\n")
	fmt.Printf("  -p=N, --parallelism=N      Use N parallel jobs.\n")
	fmt.Printf("                               Default: number of CPUs.\n")
	fmt.Printf("  file                       output file name or full path.\n")
	fmt.Printf("                               Default: %s for SHA1\n", fileSha1Default)
	fmt.Printf("                               Default: %s for NTLM\n", fileNtlmDefault)
	fmt.Printf("\n")

}

func header() {
	fmt.Printf("\n")
	fmt.Printf("  HIBPDL, Pwned password hash download tool\n")
	fmt.Printf("      GitHub:  %s\n", "https://github.com/filimonic/HIBPDL")
	fmt.Printf("      Version: %s\n", hibpdownloader.Version())
	fmt.Printf("\n")
}
