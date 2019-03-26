package main

// contract anything that can be indexed with go-whosonfirst-index
// in to line-separated geojson (20190326/thisisaaronland)

// THIS IS WORK IN PROGRESS SO ALL THE USUAL CAVEATS APPLY

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-index"
	_ "github.com/whosonfirst/go-whosonfirst-index/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func CallbackForWriter(wr io.Writer) func(fh io.Reader, ctx context.Context, args ...interface{}) error {

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		ext := filepath.Ext(path)

		if ext != ".geojson" {
			return nil // eg remarks.md files
		}

		f, err := feature.LoadGeoJSONFeatureFromReader(fh)

		if err != nil {
			return err
		}

		var stub interface{}

		err = json.Unmarshal(f.Bytes(), &stub)

		if err != nil {
			return err
		}

		body, err := json.Marshal(stub)

		if err != nil {
			return err
		}

		_, err = wr.Write(body)

		if err != nil {
			return err
		}

		wr.Write([]byte("\n"))

		return nil
	}

	return cb
}

func main() {

	mode := flag.String("mode", "repo", "...")
	dryrun := flag.Bool("dryrun", false, "...")
	// verbose := flag.Bool("verbose", false, "...")
	stdout := flag.Bool("stdout", false, "...")

	flag.Parse()

	global_writers := make([]io.Writer, 0)

	if *dryrun {
		global_writers = append(global_writers, ioutil.Discard)
	}

	if *stdout {
		global_writers = append(global_writers, os.Stdout)
	}

	t1 := time.Now()

	for _, path := range flag.Args() {

		abs_path, err := filepath.Abs(path)

		if err != nil {
			log.Fatal(err)
		}

		writers := make([]io.Writer, 0)

		for _, wr := range global_writers {
			writers = append(writers, wr)
		}

		if !*dryrun {

			// sudo make all this configurable...
			// (20190326/thisisaaronland)

			// root := filepath.Dir(abs_path)
			// fname := "data.txt"

			root := "/usr/local/data/whosonfirst-data-geojson-ls"
			fname := fmt.Sprintf("%s.txt", filepath.Base(abs_path))

			data_path := filepath.Join(root, fname)
			log.Println("WRITE", data_path)

			data_fh, err := os.OpenFile(data_path, os.O_RDWR|os.O_CREATE, 0644)

			if err != nil {
				log.Fatal(err)
			}

			defer data_fh.Close()

			writers = append(writers, data_fh)
		}

		if len(writers) == 0 {
			log.Fatal("Nowhere to write")
		}

		wr := io.MultiWriter(writers...)
		cb := CallbackForWriter(wr)

		i, err := index.NewIndexer(*mode, cb)

		if err != nil {
			log.Fatal(err)
		}

		ta := time.Now()

		err = i.IndexPath(abs_path)

		log.Printf("time to contract %s %v\n", path, time.Since(ta))

		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("time to contract all %v\n", time.Since(t1))
}
