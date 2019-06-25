package main

// expand line-separated geojson in to the standard data/123/456/123456.geojson
// directory structure (20190326/thisisaaronland)

// THIS IS WORK IN PROGRESS SO ALL THE USUAL CAVEATS APPLY

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	root := flag.String("root", "", "...")

	dryrun := flag.Bool("dryrun", false, "...")
	// verbose := flag.Bool("verbose", false, "...")

	flag.Parse()

	abs_root, err := filepath.Abs(*root)

	if err != nil {
		log.Fatal(err)
	}

	abs_data := filepath.Join(abs_root, "data")

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		f, err := feature.LoadGeoJSONFeatureFromReader(fh)

		if err != nil {
			return err
		}

		str_id := f.Id()

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return err
		}

		tree, err := uri.Id2Path(id)

		if err != nil {
			return err
		}

		abs_tree := filepath.Join(abs_data, tree)

		// HOW TO: (20190326/thisisaaronland)
		// determine if this is an alt file or not
		// if it is an alt file, what its file name is
		// format geojson properly
		// is go-whosonfirst-export "good enough is perfect" yet? - for example
		// I don't think it knows about alt files...

		fname := fmt.Sprintf("%d.geojson", id) // THIS IS ONLY FOR DEBUGGING SEE ABOVE...

		abs_path := filepath.Join(abs_tree, fname)

		_, err = os.Stat(abs_tree)

		if os.IsNotExist(err) {

			err := os.MkdirAll(abs_tree, 0755)

			if err != nil {
				return err
			}
		}

		var wr io.Writer

		if *dryrun {
			wr = ioutil.Discard
		} else {

			fh, err := os.Open(abs_path)

			if err != nil {
				return err
			}

			wr = fh
		}

		_, err = wr.Write(f.Bytes()) // THIS IS ONLY FOR DEBUGGING SEE ABOVE...

		if err != nil {
			return err
		}

		log.Printf("wrote %d to %s\n", id, abs_path)
		return nil
	}

	i, err := index.NewIndexer("geojson-ls", cb)

	if err != nil {
		log.Fatal(err)
	}

	t1 := time.Now()

	for _, path := range flag.Args() {

		ta := time.Now()

		err := i.IndexPath(path)

		log.Printf("time to contract %s %v\n", path, time.Since(ta))

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("time to contract all %v\n", time.Since(t1))
}
