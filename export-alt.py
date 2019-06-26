#!/usr/bin/env python

import mapzen.whosonfirst.utils
import mapzen.whosonfirst.export

import os
import logging
logging.basicConfig(level=logging.FATAL)

def crawl_filelist(root, filelist):

    # this assumes something like:
    # git status --porcelain --untracked-files=all | grep '.geojson' | awk '{ print $2 }' > new.txt
    
    filelist = os.path.abspath(filelist)
    fh = open(filelist, "r")

    for ln in fh:
        ln = ln.strip()
        path = os.path.join(root, ln)
        feature = mapzen.whosonfirst.utils.load_file(path)
        yield feature
        
if __name__ == "__main__":

    import optparse
    opt_parser = optparse.OptionParser()

    opt_parser.add_option('-r', '--root', dest='root', action='store', default='.', help='...')
    opt_parser.add_option('-f', '--filelist', dest='filelist', action='store', default='None', help='...')            

    opt_parser.add_option('-d', '--debug', dest='debug', action='store_true', default=False, help='Go through the motions but don\'t export anything')    
    opt_parser.add_option('-v', '--verbose', dest='verbose', action='store_true', default=False, help='Be chatty (default is false)')
    options, args = opt_parser.parse_args()

    if options.verbose:	
        logging.basicConfig(level=logging.DEBUG)
    else:
        logging.basicConfig(level=logging.INFO)

    root = os.path.abspath(options.root)
    data = os.path.join(root, "data")

    exporter = mapzen.whosonfirst.export.flatfile(data)

    crawl = None

    if options.filelist:
        crawl = crawl_filelist(root, options.filelist)
    else:
        crawl = mapzen.whosonfirst.utils.crawl(data, inflate=True)

    for feature in crawl:

        props = feature["properties"]
        alt_label = props["wof:alt_label"]
        alt_parts = alt_label.split("-")

        kwargs = {
            "alt": True,
            "source": alt_parts[0],
        }

        if len(alt_parts) > 1:
            kwargs["function"] = alt_parts[1]

        if len(alt_parts) > 2:
            kwargs["extras"] = alt_parts[2:]

        try:
            exporter.export_alt_feature(feature, **kwargs)
        except Exception, e:
            print "ERROR", props["wof:id"], e
