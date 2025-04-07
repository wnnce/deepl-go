package deepl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

var client, _ = NewDeepl(Config{
	AuthKey:     "<token>",
	Timeout:     10 * time.Second,
	AccountType: FreeAccount,
})

func TestDeepl_TextTranslate_Sync(t *testing.T) {
	result, err := client.TextTranslate("hello", "ZH").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}

func TestDeepl_TextTranslate_Async(t *testing.T) {
	done := make(chan struct{})
	client.TextTranslate("hello", "ZH").Async(func(ctx context.Context, result *TextResult, err error) {
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(result)
		close(done)
	})
	log.Println("wait callback ....")
	<-done
}

func TestDeepl_TextsTranslate_Sync(t *testing.T) {
	result, err := client.TextsTranslate([]string{"hello", "world"}, "ZH").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range result {
		log.Println(item)
	}
}

func TestDeepl_TextsTranslate_Async(t *testing.T) {
	done := make(chan struct{})
	client.TextsTranslate([]string{"hello", "world"}, "ZH").Async(func(ctx context.Context, result []*TextResult, err error) {
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range result {
			log.Println(item)
		}
		close(done)
	})
	log.Println("wait callback ...")
	<-done
}

func TestDeepl_TextsTranslateWithContext(t *testing.T) {
	done := make(chan struct{})
	ctx := context.WithValue(context.Background(), "key", "value")
	client.TextsTranslateWithContext(ctx, []string{"hello", "world"}, "EN", "ZH").Async(func(ctx context.Context, result []*TextResult, err error) {
		fmt.Println(ctx.Value("key").(string))
		if err != nil {
			log.Fatalln(err)
		}
		for _, itme := range result {
			log.Println(itme)
		}
		close(done)
	})
	<-done
}

func TestDeepl_TextTranslateWithParams(t *testing.T) {
	params := AcquireTextTranslateParams()
	params.Text = []string{"hello", "world"}
	params.TargetLang = "ZH"
	done := make(chan struct{})
	client.TextTranslateWithParams(context.Background(), params).Async(func(ctx context.Context, result []*TextResult, err error) {
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range result {
			log.Println(item)
		}
		close(done)
	})
	<-done
	RecycleParams(params)
}

func TestDeepl_Usage(t *testing.T) {
	result, err := client.Usage().Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}

func TestDeepl_Languages(t *testing.T) {
	result, err := client.Languages().Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(len(result))
	for _, language := range result {
		log.Println(language)
	}
}

func TestDeepl_LanguagesWithContext(t *testing.T) {
	ctx := context.Background()
	done := make(chan struct{})
	client.LanguagesWithContext(ctx, LanguagesTypeTarget).Async(func(ctx context.Context, result []LanguageResult, err error) {
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(len(result))
		for _, language := range result {
			log.Println(language)
		}
		close(done)
	})
	<-done
}

func TestDeepl_DocumentTranslate(t *testing.T) {
	file, err := os.Open("/home/cola/Downloads/input.pdf")
	if err != nil {
		log.Fatalln(err)
	}
	result, err := client.DocumentTranslate(file, file.Name(), "EN").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}

func TestDeepl_CheckDocumentStatus(t *testing.T) {
	status, err := client.CheckDocumentStatus("document_id", "document_key").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(status)
}

func TestDeepl_DownloadDocument(t *testing.T) {
	create, err := os.Create("output.pdf")
	if err != nil {
		log.Fatalln(err)
	}
	file, err := client.DownloadDocument("document_id", "document_key").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = io.Copy(create, bytes.NewBuffer(file))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestDeepl_ListGlossaryPairsWithContext(t *testing.T) {
	pairs, err := client.ListGlossaryPairsWithContext(context.Background()).Sync()
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range pairs {
		log.Println(item)
	}
}

func TestDeepl_CreateGlossaryWithContext(t *testing.T) {
	body := AcquireCreateGlossaryParams()
	defer RecycleParams(body)
	body.Name = "demo"
	body.SourceLang = "en"
	body.TargetLang = "Zh"
	body.Entries = "Hello\tGuten Tag"
	body.EntriesFormat = EntriesFormatTSV
	result, err := client.CreateGlossaryWithContext(context.Background(), body).Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}

func TestDeepl_ListGlossariesWithContext(t *testing.T) {
	result, err := client.ListGlossariesWithContext(context.Background()).Sync()
	if err != nil {
		log.Fatalln(err)
	}
	for _, item := range result {
		log.Println(item)
	}
	log.Println(len(result))
}

func TestDeepl_GlossaryDetailWithContext(t *testing.T) {
	result, err := client.GlossaryDetailWithContext(context.Background(), "glossary_id").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}

func TestDeepl_GlossaryEntriesWithContext(t *testing.T) {
	result, err := client.GlossaryEntriesWithContext(context.Background(), "glossary_id", "text/tab-separated-values").Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(result)
}

func TestDeepl_DeleteGlossaryWithContext(t *testing.T) {
	_, err := client.DeleteGlossaryWithContext(context.Background(), "glossary_id").Sync()
	if err != nil {
		log.Fatalln(err)
	}
}
