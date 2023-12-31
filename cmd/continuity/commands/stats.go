/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var StatsCmd = &cobra.Command{
	Use:   "stats <manifest>",
	Short: "display statistics about the specified manifest",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalln("please specify a manifest")
		}

		bm, err := readManifestFile(args[0])
		if err != nil {
			log.Fatalf("error reading manifest: %v", err)
		}

		var stats struct {
			resources   int
			files       int
			directories int
			totalSize   int64
			symlinks    int
		}

		for _, entry := range bm.Resource {
			stats.resources++
			stats.totalSize += int64(entry.Size)

			mode := os.FileMode(entry.Mode)
			if mode.IsRegular() {
				stats.files += len(entry.Path) // count hardlinks!
			} else if mode.IsDir() {
				stats.directories++
			} else if mode&os.ModeSymlink != 0 {
				stats.symlinks++
			}
		}

		w := newTabwriter(os.Stdout)
		defer w.Flush()

		fmt.Fprintf(w, "resources\t%v\n", stats.resources)
		fmt.Fprintf(w, "directories\t%v\n", stats.directories)
		fmt.Fprintf(w, "files\t%v\n", stats.files)
		fmt.Fprintf(w, "symlinks\t%v\n", stats.symlinks)
		fmt.Fprintf(w, "size\t%v\n", humanize.Bytes(uint64(stats.totalSize)))
	},
}
