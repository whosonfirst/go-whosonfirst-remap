package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {

	target := flag.String("target", "/usr/local/data", "where to write new (remapped) WOF files")
	mode := flag.String("mode", "repo", "...")
	dryrun := flag.Bool("dryrun", false, "...")

	flag.Parse()

	abs_target, err := filepath.Abs(*target)

	if err != nil {
		log.Fatal(err)
	}

	cb := func(fh io.Reader, ctx context.Context, args ...interface{}) error {

		path, err := index.PathForContext(ctx)

		if err != nil {
			return err
		}

		is_principal, err := utils.IsPrincipalWOFRecord(fh, ctx)

		if err != nil {
			return err
		}

		new_path := ""
		new_repo := ""

		if !is_principal {

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

			new_repo = "whosonfirst-data-alt"

			new_root := filepath.Join(abs_target, new_repo)
			new_data := filepath.Join(new_root, "data")

			new_parent := filepath.Join(new_data, tree)
			new_path = filepath.Join(new_parent, filepath.Base(path))

		} else {

			f, err := feature.LoadFeatureFromReader(fh)

			if err != nil {
				return err
			}

			id := whosonfirst.Id(f)

			country := whosonfirst.Country(f)
			country = strings.ToLower(country)

			if country == "" {
				country = "xy" // xx ?
			}

			new_repo = fmt.Sprintf("whosonfirst-data-%s", country)

			new_root := filepath.Join(abs_target, new_repo)
			new_data := filepath.Join(new_root, "data")

			p, err := uri.Id2AbsPath(new_data, id)

			if err != nil {
				return err
			}

			new_path = p
		}

		if *dryrun {
			log.Printf("move %s to %s\n", path, new_path)
			return nil
		}

		path_root := filepath.Dir(new_path)

		_, err = os.Stat(path_root)

		if os.IsNotExist(err) {

			err = os.MkdirAll(path_root, 0755)

			if err != nil {
				return err
			}
		}

		in, err := os.Open(path)

		if err != nil {
			return err
		}

		defer in.Close()

		out, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			return err
		}

		defer out.Close()

		_, err = io.Copy(out, in)

		if err != nil {
			return err
		}

		// update wof:repo in new_path here...

		return nil
	}

	i, err := index.NewIndexer(*mode, cb)

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		err := i.IndexPath(path)

		if err != nil {
			log.Fatal(err)
		}
	}
}