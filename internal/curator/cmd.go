package curator

/*
	---------------------------------------------------------------------
	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
	----------------------------------------------------------------------
*/

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	jsonFlagStr         = "json"
	indexFlagStr        = "index"
	outputFlagStr       = "output"
	outputAppendFlagStr = "append"

	// Version flags
	detailsFlagStr = "details"

	// Auto flags
	configFlagStr       = "conf"
	generateFlagStr     = "generate"
	keysFlagStr         = "keys"
	bloomWorkersFlagStr = "workers-bloom"
	indexWorkersFlagStr = "workers-index"
	sortWorkersFlagStr  = "workers-sort"

	// Normalize flags
	targetFlagStr     = "target"
	formatFlagStr     = "format"
	recursiveFlagStr  = "recursive"
	skipPrefixFlagStr = "skip-prefix"
	skipSuffixFlagStr = "skip-suffix"

	// Filter flags
	workersFlagStr      = "workers"
	filterSizeFlagStr   = "filter-size"
	filterHashesFlagStr = "filter-hashes"
	filterLoadFlagStr   = "filter-load"
	filterSaveFlagStr   = "filter-save"

	// Index flags
	keyFlagStr       = "key"
	noCleanupFlagStr = "no-cleanup"

	tempDirFlagStr = "temp"

	// Sort flags
	maxMemoryFlagStr     = "max-memory"
	checkFlagStr         = "check"
	maxGoRoutinesFlagStr = "max-goroutines"

	// Search flags
	valueFlagStr   = "value"
	verboseFlagStr = "verbose"

	defaultMaxMemory = 1024

	// ANSI Colors
	normal    = "\033[0m"
	black     = "\033[30m"
	red       = "\033[31m"
	green     = "\033[32m"
	orange    = "\033[33m"
	blue      = "\033[34m"
	purple    = "\033[35m"
	cyan      = "\033[36m"
	gray      = "\033[37m"
	bold      = "\033[1m"
	clearln   = "\r\x1b[2K"
	upN       = "\033[%dA"
	downN     = "\033[%dB"
	underline = "\033[4m"

	// Info - Display colorful information
	Info = bold + cyan + "[*] " + normal
	// Warn - Warn a user
	Warn = bold + red + "[!] " + normal
	// Debug - Display debug information
	Debug = bold + purple + "[-] " + normal
	// Woot - Display success
	Woot = bold + green + "[$] " + normal

	kb = 1024
	mb = kb * 1024
	gb = mb * 1024
)

var rootCmd = &cobra.Command{
	Use:   "leakdb-curator",
	Short: "Curate data sets for use with LeakDB",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		autoParseFlags(cmd, args)
	},
}

func init() {

	// Version
	versionCmd.Flags().BoolP(detailsFlagStr, "d", false, "show additional version details")
	rootCmd.AddCommand(versionCmd)

	// Main
	rootCmd.Flags().StringSliceP(keysFlagStr, "k", []string{"user", "email"}, "Comma separated list of key(s): email, user, domain")
	rootCmd.Flags().StringP(tempDirFlagStr, "T", "", "directory for temp files (default: cwd)")
	rootCmd.Flags().StringP(jsonFlagStr, "j", "", "input file/directory of normalized json file(s)")
	rootCmd.Flags().StringP(outputFlagStr, "o", "", "output directory")
	rootCmd.Flags().UintP(bloomWorkersFlagStr, "W", uint(1), "max number of bloom filter workers")
	rootCmd.Flags().UintP(indexWorkersFlagStr, "w", uint(runtime.NumCPU()), "max number of index workers")
	rootCmd.Flags().UintP(sortWorkersFlagStr, "s", uint(runtime.NumCPU()), "max number of sort workers")
	rootCmd.Flags().UintP(filterSizeFlagStr, "F", 8, "bloom filter size in GBs")
	rootCmd.Flags().UintP(filterHashesFlagStr, "f", 14, "number of bloom filter hash functions")
	rootCmd.Flags().StringP(filterLoadFlagStr, "L", "", "load existing bloom filter from saved file")
	rootCmd.Flags().StringP(filterSaveFlagStr, "S", "", "save bloom filter to file when complete")
	rootCmd.Flags().UintP(maxMemoryFlagStr, "m", defaultMaxMemory, "max memory in MBs, this is not exact! See detailed --help")

	// Normalize
	normalizeCmd.Flags().StringP(targetFlagStr, "t", "", "target directory of files")
	normalizeCmd.Flags().StringP(formatFlagStr, "f", "", "target format (see detailed help)")
	normalizeCmd.Flags().StringP(outputFlagStr, "o", "", "output json file of normalized data")
	normalizeCmd.Flags().BoolP(recursiveFlagStr, "r", false, "recursively scan directory")
	normalizeCmd.Flags().StringP(skipPrefixFlagStr, "p", "", "skip files with prefix")
	normalizeCmd.Flags().StringP(skipSuffixFlagStr, "s", "", "skip files with suffix")
	rootCmd.AddCommand(normalizeCmd)

	// Bloom
	bloomCmd.Flags().StringP(jsonFlagStr, "j", "", "input directory of normalized json file(s)")
	bloomCmd.Flags().StringP(outputFlagStr, "o", "", "output json file")
	bloomCmd.Flags().BoolP(outputAppendFlagStr, "a", false, "append output file")
	bloomCmd.Flags().UintP(workersFlagStr, "w", uint(1), "number of worker threads")
	bloomCmd.Flags().UintP(filterSizeFlagStr, "s", 8, "bloom filter size in GBs")
	bloomCmd.Flags().UintP(filterHashesFlagStr, "f", 14, "number of bloom filter hash functions")
	bloomCmd.Flags().StringP(filterLoadFlagStr, "L", "", "load existing bloom filter from saved file")
	bloomCmd.Flags().StringP(filterSaveFlagStr, "S", "", "save bloom filter to file when complete")
	rootCmd.AddCommand(bloomCmd)

	// Indexer
	indexCmd.Flags().StringP(jsonFlagStr, "j", "", "json input file")
	indexCmd.Flags().StringP(outputFlagStr, "o", "leakdb.idx", "output index file")
	indexCmd.Flags().UintP(workersFlagStr, "w", uint(runtime.NumCPU()), "number of worker threads")
	indexCmd.Flags().StringP(keyFlagStr, "k", "email", "index key can be: email, user, or domain")
	indexCmd.Flags().BoolP(noCleanupFlagStr, "N", false, "skip cleanup of temp file(s)")
	indexCmd.Flags().StringP(tempDirFlagStr, "T", "", "directory for temp files (default: cwd)")
	rootCmd.AddCommand(indexCmd)

	// Sorter
	sortCmd.Flags().StringP(indexFlagStr, "i", "", "index file to sort")
	sortCmd.Flags().StringP(outputFlagStr, "o", "", "output index file")
	sortCmd.Flags().UintP(workersFlagStr, "w", uint(runtime.NumCPU()), "number of worker threads")
	sortCmd.Flags().UintP(maxMemoryFlagStr, "m", defaultMaxMemory, "max memory in MBs, this is not exact! See detailed --help")
	sortCmd.Flags().StringP(tempDirFlagStr, "T", "", "directory for temp files (default: cwd)")
	sortCmd.Flags().BoolP(noCleanupFlagStr, "N", false, "skip cleanup temp file(s)")
	rootCmd.AddCommand(sortCmd)

	// Search
	searchCmd.Flags().StringP(indexFlagStr, "i", "", "index file to search")
	searchCmd.Flags().StringP(jsonFlagStr, "j", "", "original json file")
	searchCmd.Flags().StringP(valueFlagStr, "v", "", "value to search for")
	rootCmd.AddCommand(searchCmd)

	// Parquet
	parquetCmd.Flags().StringP(targetFlagStr, "t", "", "target JSON file to convert")
	parquetCmd.Flags().StringP(outputFlagStr, "o", "leakdb-parquet", "output directory name")
	rootCmd.AddCommand(parquetCmd)
}

// Execute - Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
