package deepl

type TextTranslateHandler func([]*TextResult, error)

type Recyclable interface {
	recycle()
}

// BaseParams is public request params
type BaseParams struct {
	SourceLang string `json:"source_lang,omitempty"`
	TargetLang string `json:"target_lang,omitempty"`
	Formality  string `json:"formality,omitempty"`
	GlossaryId string `json:"glossary_id,omitempty"`
}

type TextTranslateParams struct {
	BaseParams
	Text                 []string `json:"text,omitempty"`
	Context              string   `json:"context,omitempty"`
	ShowBilledCharacters bool     `json:"show_billed_characters"`
	SplitSentences       string   `json:"split_sentences,omitempty"`
	PreserveFormatting   bool     `json:"preserve_formatting"`
	TagHandling          string   `json:"tag_handling,omitempty"`
	OutlineDetection     bool     `json:"outline_detection,omitempty"`
	NonSplittingTags     []string `json:"non_splitting_tags,omitempty"`
	SplittingTags        []string `json:"splitting_tags,omitempty"`
	IgnoreTags           []string `json:"ignore_tags,omitempty"`
}

type TextImprovementParams struct {
	Text         []string `json:"text,omitempty"`
	TargetLang   string   `json:"target_lang,omitempty"`
	WritingStyle string   `json:"writing_style,omitempty"`
	Tone         string   `json:"tone,omitempty"`
}

type DocumentTranslateParams struct {
	BaseParams
	Filename     string `json:"filename"`
	OutputFormat string `json:"output_format"`
}

type CreateGlossaryParams struct {
	Name          string `json:"name"`
	SourceLang    string `json:"source_lang,omitempty"`
	TargetLang    string `json:"target_lang"`
	Entries       string `json:"entries"`
	EntriesFormat string `json:"entries_format"`
}

func (self *TextTranslateParams) recycle() {
	self.Text = nil
	self.SourceLang = ""
	self.TargetLang = ""
	self.Context = ""
	self.ShowBilledCharacters = false
	self.SplitSentences = ""
	self.PreserveFormatting = false
	self.GlossaryId = ""
	self.TagHandling = ""
	self.OutlineDetection = false
	self.NonSplittingTags = nil
	self.SplittingTags = nil
	self.IgnoreTags = nil
}

func (self *DocumentTranslateParams) recycle() {
	self.SourceLang = ""
	self.TargetLang = ""
	self.Formality = ""
	self.GlossaryId = ""
	self.Filename = ""
	self.OutputFormat = ""
}

func (self *TextImprovementParams) recycle() {
	self.Text = nil
	self.TargetLang = ""
	self.WritingStyle = ""
	self.Tone = ""
}

func (self *CreateGlossaryParams) recycle() {
	self.Name = ""
	self.SourceLang = ""
	self.TargetLang = ""
	self.Entries = ""
	self.EntriesFormat = ""
}

type TextTranslateResultOptional struct {
	Translations []*TextResult `json:"translations"`
}

type TextImprovementResultOptional struct {
	Improvements []*TextResult `json:"improvements"`
}

type GlossaryPairsOptional struct {
	SupportedLanguages []PairResult `json:"supported_languages"`
}

type GlossariesOptional struct {
	Glossaries []*GlossaryResult `json:"glossaries"`
}

type TextResult struct {
	DetectedSourceLanguage string `json:"detected_source_language,omitempty"`
	Text                   string `json:"text,omitempty"`
	BilledCharacters       int    `json:"billed_characters,omitempty"`
	ModelTypeUsed          string `json:"model_type_used,omitempty"`
}

type DocumentResult struct {
	DocumentId  string `json:"document_id"`
	DocumentKey string `json:"document_key"`
}

type CheckDocumentResult struct {
	DocumentId       string `json:"document_id"`
	Status           string `json:"status"`
	SecondsRemaining int    `json:"seconds_remaining"`
}

type UsageResult struct {
	CharacterCount int64 `json:"character_count"`
	CharacterLimit int64 `json:"character_limit"`
}

type LanguageResult struct {
	Language          string `json:"language"`
	Name              string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"`
}
type PairResult struct {
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type GlossaryResult struct {
	GlossaryId   string `json:"glossary_id"`
	Ready        bool   `json:"ready"`
	Name         string `json:"name"`
	SourceLang   string `json:"source_lang"`
	TargetLang   string `json:"target_lang"`
	CreationTime string `json:"creation_time"`
	EntryCount   int64  `json:"entry_count"`
}
