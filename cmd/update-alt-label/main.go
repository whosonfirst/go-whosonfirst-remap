package main

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"	
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

func main() {

	flag.Parse()

	re_alt, err := regexp.Compile(`(?:\d+)-alt-(.*).geojson`)

	if err != nil {
		log.Fatal(err)
	}
	
	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		fname := filepath.Base(path)
		
		if !re_alt.MatchString(fname){
			return nil
		}
		
		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return err
		}

		label_rsp := gjson.GetBytes(body, "properties.wof:alt_label")

		if label_rsp.Exists(){
			return nil
		}
		
		m := re_alt.FindStringSubmatch(fname)
		alt_label := m[1]

		body, err = sjson.SetBytes(body, "properties.wof:alt_label", alt_label)

		if err != nil {
			return err
		}
		
		log.Println("UPDATE", fname)
		return nil
	}

	idx, err := index.NewIndexer("repo", cb)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		err := idx.IndexPath(path)

		if err != nil {
			log.Fatal(err)
		}
	}
	
}
