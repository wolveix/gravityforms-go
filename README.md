# Gravity Forms Go SDK

Go SDK for interacting with the Gravity Forms API v2 (WordPress plugin).

_Note: this is very incomplete, as I only implemented what I needed for a quick project. PRs welcome!_

## Usage

You can obtain your key and secret Gravity Forms settings -> REST API page. You'll want to use `Authentication (API 
version 2)`.

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wolveix/gravityforms-go"
)

func main() {
	key := "ck_your_api_key_here"
	secret := "cs_your_api_secret_here"
	service := gravityforms.New("https://your.wordpress.domain/wp-json/gf/v2", key, secret, 15*time.Second, false)

	entries, err := service.GetEntriesByFormID(1)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		// Access field ID 1.
		fmt.Println(entry.GetField("1"))
	}
}
```

## License

BSD licensed. See the [LICENSE](LICENSE) file for details.