package main

import (
	"context"
	"flag"
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

func main() {

	mode := flag.String("mode", "repo", "...")
	dryrun := flag.Bool("dryrun", false, "...")
	verbose := flag.Bool("verbose", false, "...")
	stdout := flag.Bool("stdout", false, "...")	

	flag.Parse()

	writers := make([]io.Writer, 0)

	if *dryrun {
		writers = append(writers, ioutil.Discard)
	}

	if *stdout {
		writers = append(writers, os.Stdout)
	}
	
	if len(writers) == 0 {
		log.Fatal("Nowhere to write")
	}
	
	wr := io.MultiWriter(writers...)
	
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

		if *verbose {
			log.Printf("write %s\n", path)
		}
		
		if *dryrun {
			return nil
		}
		
		_, err  = wr.Write(f.Bytes())

		if err != nil{
			return err
		}
		
		return nil
	}

	i, err := index.NewIndexer(*mode, cb)

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
