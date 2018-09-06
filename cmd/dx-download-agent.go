package main

import (
	"fmt"
	"os"

	// The dxda package should contain all core functionality
	"github.com/geetduggal/dxda"
)

// The CLI is simply a wrapper around the dxda package
func main() {
	token, method := dxda.GetToken()
	fmt.Printf("Obtained token using %s\n", method)
	fname := "../test_files/single_file.manifest.json.bz2"
	// manifest := dxda.ReadManifest()
	// dxda.DownloadManifest(manifest, token)
	var opts dxda.Opts
	opts.NumThreads = 8
	if _, err := os.Stat(fname + ".stats.db"); os.IsNotExist(err) {
		dxda.CreateManifestDB(fname)
	}
	dxda.DownloadManifestDB(fname, token, opts)
}
