# CSV Importer

The purpose of this tool is to iterate over rows within a CSV file, and create the corresponding entries within your Gravity Forms installation.

## Usage

Build the binary with:
```go
go build -o csvimport .
```

With the binary in your directory, you can then call it like so:

```shell
./csvimport --CSV_FILE_PATH='input.csv' --API_KEY=ck_your_gf_api_key --API_SECRET=cs_your_gf_api_secret --API_URL='https://your.domain.here/wp-json/gf/v2' --FORM_ID=1
```

Your CSV's first row should contain the field IDs from your Gravity Form, e.g.:

```csv
4,6,22,15
Something,Some other value,Another value,One more value
```