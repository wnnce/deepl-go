package deepl

import (
	"bytes"
	"context"
	"fmt"
	"go/types"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

var (
	uuidRegex        = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	documentIdRegex  = regexp.MustCompile(`^[0-9A-Z]{32}$`)
	documentKeyRegex = regexp.MustCompile(`^[0-9A-Z]{64}$`)
)

type Deepl struct {
	client *http.Client
	config Config
	host   string
}

func NewDeepl(config Config) (*Deepl, error) {
	if !strings.HasSuffix(config.AuthKey, ":fx") || !uuidRegex.MatchString(config.AuthKey[:len(config.AuthKey)-3]) {
		return nil, fmt.Errorf("Token does not exist or is not formatted correctly, your Token: %s ", config.AuthKey)
	}
	if config.Timeout == 0 {
		config.Timeout = DefaultConfig.Timeout
	}
	if config.AccountType < 1 || config.AccountType > 2 {
		config.AccountType = DefaultConfig.AccountType
	}
	if config.JSONEncode == nil {
		config.JSONEncode = DefaultConfig.JSONEncode
	}
	if config.JSONDecode == nil {
		config.JSONDecode = DefaultConfig.JSONDecode
	}
	client := &http.Client{
		Timeout: config.Timeout,
	}
	host := freeHost
	if config.AccountType == ProAccount {
		host = proHost
	}
	return &Deepl{
		client: client,
		config: config,
		host:   host,
	}, nil
}

// TextTranslate Is single text translate
func (self *Deepl) TextTranslate(text, target string) *CMD[*TextResult] {
	return self.TextTranslateWithContext(context.Background(), text, "", target)
}

// TextsTranslate Is multiple text translate
func (self *Deepl) TextsTranslate(texts []string, target string) *CMD[[]*TextResult] {
	return self.TextsTranslateWithContext(context.Background(), texts, "", target)
}

// TextTranslateWithSource Is specify a single text translate in the source language
func (self *Deepl) TextTranslateWithSource(text, source, target string) *CMD[*TextResult] {
	return self.TextTranslateWithContext(context.Background(), text, source, target)
}

func (self *Deepl) TextsTransLateWithSource(texts []string, source, target string) *CMD[[]*TextResult] {
	return self.TextsTranslateWithContext(context.Background(), texts, source, target)
}

// TextTranslateWithContext All individual text translators end up calling a method
// that creates the request parameters and recycles them after the call is complete
func (self *Deepl) TextTranslateWithContext(ctx context.Context, text, source, target string) *CMD[*TextResult] {
	return NewCMD(ctx, func() (*TextResult, error) {
		body := AcquireTextTranslateParams()
		body.Text = []string{text}
		body.SourceLang = source
		body.TargetLang = target
		defer RecycleParams(body)
		result, err := self.doTextTranslate(ctx, body)
		if err == nil && result != nil && len(result) > 0 {
			return result[0], nil
		}
		return nil, err
	})

}

// TextsTranslateWithContext All multiple text translations will eventually call the method
func (self *Deepl) TextsTranslateWithContext(ctx context.Context, texts []string, source, target string) *CMD[[]*TextResult] {
	return NewCMD(ctx, func() ([]*TextResult, error) {
		body := AcquireTextTranslateParams()
		body.Text = texts
		body.SourceLang = source
		body.TargetLang = target
		defer RecycleParams(body)
		return self.doTextTranslate(ctx, body)
	})
}

func (self *Deepl) TextTranslateWithParams(ctx context.Context, body *TextTranslateParams) *CMD[[]*TextResult] {
	return NewCMD(ctx, func() ([]*TextResult, error) {
		return self.doTextTranslate(ctx, body)
	})
}

// All text translations end up calling the method
func (self *Deepl) doTextTranslate(ctx context.Context, body *TextTranslateParams) ([]*TextResult, error) {
	request, err := self.createRequestWithJSON(ctx, textTranslateUri, http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	result := &TextTranslateResultOptional{}
	if err = self.doRequest(request, result); err != nil {
		return nil, err
	}
	return result.Translations, nil
}

// Usage Is Check Usage and Limits
func (self *Deepl) Usage() *CMD[UsageResult] {
	return self.UsageWithContext(context.Background())
}

// UsageWithContext Usage of the transitive context
func (self *Deepl) UsageWithContext(ctx context.Context) *CMD[UsageResult] {
	return NewCMD(ctx, func() (UsageResult, error) {
		request, err := self.createRequestWithJSON(ctx, usageUri, http.MethodGet, nil)
		var result UsageResult
		if err != nil {
			return result, err
		}
		err = self.doRequest(request, &result)
		return result, err
	})
}

// Languages Is retrieve supported languages
func (self *Deepl) Languages() *CMD[[]LanguageResult] {
	return self.LanguagesWithContext(context.Background(), LanguagesTypeSource)
}

// LanguagesWithType Is retrieve supported languages for type
func (self *Deepl) LanguagesWithType(t string) *CMD[[]LanguageResult] {
	return self.LanguagesWithContext(context.Background(), t)
}

// LanguagesWithContext LanguagesWithType of the transitive context
func (self *Deepl) LanguagesWithContext(ctx context.Context, t string) *CMD[[]LanguageResult] {
	return NewCMD(ctx, func() ([]LanguageResult, error) {
		uri := languagesUri + "?type=" + t
		request, err := self.createRequestWithJSON(ctx, uri, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}
		result := make([]LanguageResult, 0)
		err = self.doRequest(request, &result)
		return result, err
	})
}

// TextImprovement Is single text improvement
func (self *Deepl) TextImprovement(text string) *CMD[*TextResult] {
	return self.TextImprovementWithContext(context.Background(), text)
}

// TextsImprovement Is multiple text improvement
func (self *Deepl) TextsImprovement(texts []string) *CMD[[]*TextResult] {
	return self.TextsImprovementWithContext(context.Background(), texts)
}

// TextImprovementWithContext Methods that are called by all single text improvement can pass the context
func (self *Deepl) TextImprovementWithContext(ctx context.Context, text string) *CMD[*TextResult] {
	return NewCMD(ctx, func() (*TextResult, error) {
		body := AcquireTextImprovementParams()
		body.Text = []string{text}
		defer RecycleParams(body)
		result, err := self.doTextImprovement(ctx, body)
		if err == nil && result != nil && len(result) > 0 {
			return result[0], nil
		}
		return nil, err
	})
}

// TextsImprovementWithContext Methods that are called by all multiple text improvement can pass the context
func (self *Deepl) TextsImprovementWithContext(ctx context.Context, texts []string) *CMD[[]*TextResult] {
	return NewCMD(ctx, func() ([]*TextResult, error) {
		body := AcquireTextImprovementParams()
		body.Text = texts
		defer RecycleParams(body)
		return self.doTextImprovement(ctx, body)
	})
}

func (self *Deepl) TextImprovementWithParams(ctx context.Context, body *TextImprovementParams) *CMD[[]*TextResult] {
	return NewCMD(ctx, func() ([]*TextResult, error) {
		return self.doTextImprovement(ctx, body)
	})
}

// All text improvement methods are finally called
func (self *Deepl) doTextImprovement(ctx context.Context, body *TextImprovementParams) ([]*TextResult, error) {
	request, err := self.createRequestWithJSON(ctx, textImprovementUri, http.MethodPost, body)
	if err != nil {
		return nil, err
	}
	result := &TextImprovementResultOptional{}
	if err = self.doRequest(request, result); err != nil {
		return nil, err
	}
	return result.Improvements, nil
}

// DocumentTranslate Is simple document translate
// document Documents that need to be translated
// filename The name of the file to be translated
func (self *Deepl) DocumentTranslate(document io.Reader, filename, target string) *CMD[DocumentResult] {
	return self.DocumentTranslateWithContext(context.Background(), document, filename, "", target)
}

// DocumentTranslateWithSource Is simple document translate for sourceLang
func (self *Deepl) DocumentTranslateWithSource(document io.Reader, filename, source, target string) *CMD[DocumentResult] {
	return self.DocumentTranslateWithContext(context.Background(), document, filename, source, target)
}

// DocumentTranslateWithContext DocumentTranslateWithSource of the transitive context
func (self *Deepl) DocumentTranslateWithContext(ctx context.Context, document io.Reader, filename, source, target string) *CMD[DocumentResult] {
	return NewCMD(ctx, func() (DocumentResult, error) {
		body := AcquireDocumentTranslateParams()
		body.SourceLang = source
		body.TargetLang = target
		defer RecycleParams(body)
		return self.doDocumentTranslate(ctx, document, filename, body)
	})
}

func (self *Deepl) DocumentTransWithParams(ctx context.Context, document io.Reader, filename string, body *DocumentTranslateParams) *CMD[DocumentResult] {
	return NewCMD(ctx, func() (DocumentResult, error) {
		return self.doDocumentTranslate(ctx, document, filename, body)
	})
}

// All document translate methods that are finally called
// filename and body.Filename is used for the filename of the file field
// in the form and the separate filename field in the form
func (self *Deepl) doDocumentTranslate(ctx context.Context, document io.Reader, filename string, body *DocumentTranslateParams) (DocumentResult, error) {
	var result DocumentResult
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer recycleBuffer(buffer)
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return result, err
	}
	if _, err = io.Copy(part, document); err != nil {
		return result, err
	}
	params := map[string]string{
		"filename":      body.Filename,
		"source_lang":   body.SourceLang,
		"target_lang":   body.TargetLang,
		"output_format": body.OutputFormat,
		"formality":     body.Formality,
		"glossary_id":   body.GlossaryId,
	}
	for k, v := range params {
		if v == "" || strings.TrimSpace(v) == "" {
			continue
		}
		if err = writer.WriteField(k, v); err != nil {
			return result, err
		}
	}
	writer.Close()
	request, err := self.createRequest(ctx, documentTranslateUri, http.MethodPost, writer.FormDataContentType(), buffer)
	if err != nil {
		return result, err
	}
	err = self.doRequest(request, &result)
	return result, err
}

func (self *Deepl) CheckDocumentStatus(documentId, documentKey string) *CMD[CheckDocumentResult] {
	return self.CheckDocumentStatusWithContext(context.Background(), documentId, documentKey)
}

func (self *Deepl) CheckDocumentStatusWithContext(ctx context.Context, documentId, documentKey string) *CMD[CheckDocumentResult] {
	return NewCMD(ctx, func() (CheckDocumentResult, error) {
		var result CheckDocumentResult
		if err := self.validateDocumentIdAndKey(documentId, documentKey); err != nil {
			return result, err
		}
		requestUri := fmt.Sprintf(checkDocumentStatusUri, documentId)
		buffer := bufferPool.Get().(*bytes.Buffer)
		defer recycleBuffer(buffer)
		buffer.WriteString("{\"document_key\":\"" + documentKey + "\"}")
		request, err := self.createRequestWithJSON(ctx, requestUri, http.MethodPost, buffer)
		if err != nil {
			return result, err
		}
		err = self.doRequest(request, &result)
		return result, err
	})
}

func (self *Deepl) DownloadDocument(documentId, documentKey string) *CMD[[]byte] {
	return self.DownloadDocumentWithContext(context.Background(), documentId, documentKey)
}

func (self *Deepl) DownloadDocumentWithContext(ctx context.Context, documentId, documentKey string) *CMD[[]byte] {
	return NewCMD(ctx, func() ([]byte, error) {
		if err := self.validateDocumentIdAndKey(documentId, documentKey); err != nil {
			return nil, err
		}
		requestUri := fmt.Sprintf(downloadDocumentsUri, documentId)
		buffer := bufferPool.Get().(*bytes.Buffer)
		defer recycleBuffer(buffer)
		buffer.WriteString("{\"document_key\": \"" + documentKey + "\"}")
		request, err := self.createRequestWithJSON(ctx, requestUri, http.MethodPost, buffer)
		if err != nil {
			return nil, err
		}
		result := make([]byte, 0)
		err = self.doRequest(request, &result)
		return result, err
	})
}

func (self *Deepl) validateDocumentIdAndKey(id, key string) error {
	if !documentIdRegex.MatchString(id) {
		return fmt.Errorf("the document-id format is incorrect, your document-id: %s", id)
	}
	if !documentKeyRegex.MatchString(key) {
		return fmt.Errorf("the document-key format is incorrect, your document-key: %s", key)
	}
	return nil
}

// ListGlossaryPairs Is list language pairs supported by glossaries
func (self *Deepl) ListGlossaryPairs() *CMD[[]PairResult] {
	return self.ListGlossaryPairsWithContext(context.Background())
}

func (self *Deepl) ListGlossaryPairsWithContext(ctx context.Context) *CMD[[]PairResult] {
	return NewCMD(ctx, func() ([]PairResult, error) {
		request, err := self.createRequestWithJSON(ctx, listGlossaryPairsUri, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}
		result := &GlossaryPairsOptional{}
		if err = self.doRequest(request, result); err != nil {
			return nil, err
		}
		return result.SupportedLanguages, nil
	})
}

func (self *Deepl) CreateGlossary(body *CreateGlossaryParams) *CMD[*GlossaryResult] {
	return self.CreateGlossaryWithContext(context.Background(), body)
}

func (self *Deepl) CreateGlossaryWithContext(ctx context.Context, body *CreateGlossaryParams) *CMD[*GlossaryResult] {
	return NewCMD(ctx, func() (*GlossaryResult, error) {
		request, err := self.createRequestWithJSON(ctx, createGlossaryUri, http.MethodPost, body)
		if err != nil {
			return nil, err
		}
		result := &GlossaryResult{}
		err = self.doRequest(request, result)
		return result, err
	})
}

// ListGlossaries Is list all glossaries
func (self *Deepl) ListGlossaries() *CMD[[]*GlossaryResult] {
	return self.ListGlossariesWithContext(context.Background())
}

func (self *Deepl) ListGlossariesWithContext(ctx context.Context) *CMD[[]*GlossaryResult] {
	return NewCMD(ctx, func() ([]*GlossaryResult, error) {
		request, err := self.createRequestWithJSON(ctx, listGlossariesUri, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}
		result := &GlossariesOptional{}
		if err = self.doRequest(request, result); err != nil {
			return nil, err
		}
		return result.Glossaries, nil
	})
}

// GlossaryDetail Is retrieve glossary details
func (self *Deepl) GlossaryDetail(glossaryId string) *CMD[*GlossaryResult] {
	return self.GlossaryDetailWithContext(context.Background(), glossaryId)
}

func (self *Deepl) GlossaryDetailWithContext(ctx context.Context, glossaryId string) *CMD[*GlossaryResult] {
	return NewCMD(ctx, func() (*GlossaryResult, error) {
		if !uuidRegex.MatchString(glossaryId) {
			return nil, fmt.Errorf("GlossaryId does not exist or is not formatted correctly, your glossaryId: %s ", glossaryId)
		}
		requestUri := fmt.Sprintf(glossaryDetailsUri, glossaryId)
		request, err := self.createRequestWithJSON(ctx, requestUri, http.MethodGet, nil)
		if err != nil {
			return nil, err
		}
		result := &GlossaryResult{}
		err = self.doRequest(request, result)
		return result, err
	})
}

// GlossaryEntries Is retrieve glossary entries
func (self *Deepl) GlossaryEntries(glossaryId, accept string) *CMD[string] {
	return self.GlossaryEntriesWithContext(context.Background(), glossaryId, accept)
}

func (self *Deepl) GlossaryEntriesWithContext(ctx context.Context, glossaryId, accept string) *CMD[string] {
	return NewCMD(ctx, func() (string, error) {
		if !uuidRegex.MatchString(glossaryId) {
			return "", fmt.Errorf("GlossaryId does not exist or is not formatted correctly, your glossaryId: %s ", glossaryId)
		}
		requestUri := fmt.Sprintf(glossaryEntriesUri, glossaryId)
		request, err := self.createRequestWithJSON(ctx, requestUri, http.MethodGet, nil)
		if err != nil {
			return "", err
		}
		request.Header.Set("Accept", accept)
		result := make([]byte, 0)
		err = self.doRequest(request, &result)
		return string(result), err
	})
}

func (self *Deepl) DeleteGlossary(glossaryId string) *CMD[struct{}] {
	return self.DeleteGlossaryWithContext(context.Background(), glossaryId)
}

func (self *Deepl) DeleteGlossaryWithContext(ctx context.Context, glossaryId string) *CMD[struct{}] {
	return NewCMD(ctx, func() (struct{}, error) {
		if !uuidRegex.MatchString(glossaryId) {
			return struct{}{}, fmt.Errorf("GlossaryId does not exist or is not formatted correctly, your glossaryId: %s ", glossaryId)
		}
		requestUri := fmt.Sprintf(deleteGlossaryUri, glossaryId)
		request, err := self.createRequestWithJSON(ctx, requestUri, http.MethodDelete, nil)
		if err != nil {
			return struct{}{}, err
		}
		return struct{}{}, self.doRequest(request, nil)
	})
}

// Send a request and deserialize the response Body through generics
func (self *Deepl) doRequest(req *http.Request, result any) error {
	response, err := self.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err = self.handlerResponse(response, result); err != nil {
		return err
	}
	return nil
}

// Process the response. If the response code is not 200, return the corresponding error.
// Otherwise, determine the type of result. If it is []byte, return the body directly.
// Otherwise, perform json deserialization.
func (self *Deepl) handlerResponse(resp *http.Response, result any) error {
	switch resp.StatusCode {
	case 200, 201, 204:
		if result == nil {
			return nil
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		switch result.(type) {
		case *[]byte:
			*(result.(*[]byte)) = body
		default:
			return self.config.JSONDecode(body, result)
		}
	case 400:
		body, err := io.ReadAll(resp.Body)
		if err != nil || len(body) == 0 {
			return ErrBadRequest
		}
		return NewError(400, string(body))
	case 401:
		return ErrAuthorization
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFount
	case 413:
		return ErrLimit
	case 414:
		return ErrLongURL
	case 415:
		return ErrNotAccept
	case 429:
		return ErrManyRequests
	case 456:
		return ErrQuotaExceeded
	case 500:
		return ErrInternal
	case 503, 504:
		return ErrResourceUnavailable
	case 529:
		return ErrManyRequests2
	default:
		return NewError(resp.StatusCode, "Unknown error")
	}
	return nil
}

func (self *Deepl) createRequestWithJSON(ctx context.Context, uri, method string, body any) (*http.Request, error) {
	switch v := body.(type) {
	case types.Nil:
		return self.createRequest(ctx, uri, method, "application/json", nil)
	case io.Reader:
		return self.createRequest(ctx, uri, method, "application/json", v)
	default:
		encode, err := self.config.JSONEncode(body)
		if err != nil {
			return nil, err
		}
		buffer := bufferPool.Get().(*bytes.Buffer)
		defer recycleBuffer(buffer)
		buffer.Write(encode)
		return self.createRequest(ctx, uri, method, "application/json", buffer)
	}
}

func (self *Deepl) createRequest(ctx context.Context, uri, method, contentType string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, self.host+uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "DeepL-Auth-Key "+self.config.AuthKey)
	request.Header.Set("Content-Type", contentType)
	return request, nil
}
