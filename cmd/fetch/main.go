package main

import (
	"fmt"
	"os"

	"gihub.com/mrmarble/se-vende/providers"
	"github.com/k0kubun/pp/v3"
)

func main() {
	arg := os.Args[1]

	if url := providers.GetURL(arg); url != "" {
		item, err := providers.NewItem(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		pp.Println(item)
	} else {
		fmt.Println("URL not supported")
	}
}
