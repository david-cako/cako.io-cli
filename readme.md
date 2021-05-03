# cako.io cli
A command line interface for archiving and serving pages from cako.io.

## Usage:
```
  -outDir string
        Output directory to save files (default "./saved/")
  -page string
        Crawl specified page name only
  -serve
        Serve locally saved files
  -skipAssets
        Only crawl html files
  -skipExisting
        Skip crawling pages already in output directory
```

## Example:

```
cako -skipExisting
```

After archiving, run with `-serve` to start a local server for the given output directory.

```
cako -serve
```