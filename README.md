[![Build Status](https://travis-ci.com/palette-software/insight-tester.svg?branch=master&token=qWG5FJDvsjLrsJpXgxSJ)](https://travis-ci.com/palette-software/insight-tester)

# insight-tester
Testing tools for Palette Insight. These components are all written in Go. This toolset consists of:
* [DB connector module]
* [CSV forker](csv_forker)
* Insight [sanity checks](dbcheck)
* [fake Insight Server](fake-insight-server)

#### CSV forker
This module can replicate incoming CSVs to test databases. This module might come in handy when one would like to implement a more efficient way of Insight [LoadTables](https://github.com/palette-software/insight-gp-import) or [Reporting](https://github.com/palette-software/insight-data-model), where you would like to make sure that the result is the same as the current calculations.

#### Insight sanity checks
This is the only package in this repository, from which an RPM package is created and it is actually part of the [Palette Insight] product. The name of this package is `palette-insight-sanity-check`. This module can run test queries via [DB connector module], and based on the results of those tests it can tell whether [Palette Insight] is operating normally or not. The test queries are defined in [dbcheck/tests/sanity_checks.yml](dbcheck/tests/sanity_checks.yml) file.

#### Fake Insight Server
Currently, a far from up-to-date Insight Server stub which could come in handy for functional testing of [Palette Insight Agent](https://github.com/palette-software/PaletteInsightAgent)

## Contribution

### Building locally

```
go get ./...
go install -v ./...
```

### Testing

```
go get -t ./...
go test ./... -v
```

[DB connector module]: common/db-connector
[Palette Insight]: https://github.com/palette-software/palette-insight
