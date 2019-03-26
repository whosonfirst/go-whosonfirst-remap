package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-index/utils"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func RepoFromFeature(f geojson.Feature) string {

	country := CountryFromFeature(f)
	repo := fmt.Sprintf("whosonfirst-data-admin-%s", country)

	return repo
}

func CountryFromFeature(f geojson.Feature) string {

	country := whosonfirst.Country(f)
	country = strings.ToLower(country)

	if country == "" {
		country = "xy" // xx ?
	}

	return country
}

func main() {

	target := flag.String("target", "/usr/local/data", "where to write new (remapped) WOF files")
	mode := flag.String("mode", "repo", "...")
	dryrun := flag.Bool("dryrun", false, "...")
	verbose := flag.Bool("verbose", false, "...")

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

			ext := filepath.Ext(path)

			if ext != ".geojson" {
				return nil // eg remarks.md files
			}

			f, err := feature.LoadGeoJSONFeatureFromReader(fh)

			if err != nil {
				log.Printf("FAILED TO LOAD %s %s\n", path, err)
				return nil
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

			root := filepath.Dir(path)

			principal_fname := fmt.Sprintf("%d.geojson", id)
			principal_path := filepath.Join(root, principal_fname)

			principal_f, err := feature.LoadFeatureFromFile(principal_path)

			if err != nil {
				return err
			}

			new_repo = RepoFromFeature(principal_f)

			new_root := filepath.Join(abs_target, new_repo)
			new_data := filepath.Join(new_root, "data")

			new_parent := filepath.Join(new_data, tree)
			new_path = filepath.Join(new_parent, filepath.Base(path))

		} else {

			f, err := feature.LoadFeatureFromReader(fh)

			if err != nil {
				log.Printf("FAILED TO LOAD %s %s\n", path, err)
				return err
			}

			new_repo = RepoFromFeature(f)

			new_root := filepath.Join(abs_target, new_repo)
			new_data := filepath.Join(new_root, "data")

			id := whosonfirst.Id(f)
			p, err := uri.Id2AbsPath(new_data, id)

			if err != nil {
				return err
			}

			new_path = p
		}

		in, err := os.Open(path)

		if err != nil {
			return err
		}

		defer in.Close()

		body, err := ioutil.ReadAll(in)

		if err != nil {
			return err
		}

		body, err = sjson.SetBytes(body, "properties.wof:repo", new_repo)

		if err != nil {
			return err
		}

		rsp := gjson.GetBytes(body, "properties.wof:repo")

		if !rsp.Exists() {
			msg := fmt.Sprintf("Unable to find wof:repo in %s", path)
			return errors.New(msg)
		}

		if rsp.String() != new_repo {
			msg := fmt.Sprintf("Failed to set wof:repo (%s) in %s", new_repo, path)
			return errors.New(msg)
		}

		if *verbose {
			log.Printf("move (%s) %s to %s\n", rsp.String(), path, new_path)
		}

		if *dryrun {
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

		out, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			return err
		}

		defer out.Close()

		r := bytes.NewReader(body)

		_, err = io.Copy(out, r)

		if err != nil {
			return err
		}

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
