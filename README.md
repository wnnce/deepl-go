# Deepl translate api for Go

`deepl-go` provides a go integration of the Deepl translation API

## Features

- Supports all Deepl translation APIs
- Request parameters and buffer object pool
- Simple or customizable translation functions
- Sync and Async execution support

## Installation

Requires at least go1.18, then install `deepl-go`

```shell
go get github.com/wnnce/deepl-go
```

## Quickstart

Create deepl client

```go
package main

import (
	"encoding/json"
	"github.com/wnnce/deepl-go"
	"log"
)

func main() {
	client, err := deepl.NewDeepl(deepl.Config{
		AuthKey:    "<Token>",
		JSONEncode: json.Marshal,
		JSONDecode: json.Unmarshal,
	})
	if err != nil {
		log.Fatalln(err)
	}
}
```

Multi-text translate and async execution

```go
done := make(chan struct{})
client.TextsTranslate([]string{"hello", "world"}, "ZH").Async(func(ctx context.Context, result []*deepl.TextResult, err error) {
    if err != nil {
        log.Fatalln(err)
    }
    for _, item := range result {
        log.Printf("translation text: %s", item.Text)
    }
    close(done)
})
<-done
```

Text translate of custom parameters

```go
body := deepl.AcquireTextTranslateParams()
body.SourceLang = "EN"
body.TargetLang = "ZH"
body.Text = []string{"hello", "world"}
body.ShowBilledCharacters = true

results, err := client.TextTranslateWithParams(context.Background(), body).Sync()
if err != nil {
    log.Fatalln(err)
}

deepl.RecycleParams(body)
```

Check usage quota and limit

```go
result, err := client.Usage().Sync()
if err != nil {
    log.Fatalln(err)
}
log.Printf("used: %d, limit: %d", result.CharacterCount, result.CharacterLimit)
```

Simple document translate

```go
file, _ := os.Open("./input.pdf")
result, err := client.DocumentTranslate(file, "input.pdf", "ZH").Sync()
if err != nil {
    log.Fatalln(err)
}
log.Printf("documentId: %s, documentKey: %s", result.DocumentId, result.DocumentKey)
```

Download translation document

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
file, _ := os.Create("./output.pdf")
body, err := client.DownloadDocumentWithContext(ctx, "<document_id>", "<document_key>").Sync()
if err != nil {
    log.Fatalln(err)
}
_, err = io.Copy(file, bytes.NewBuffer(body))
```