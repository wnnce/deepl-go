package deepl

const (
	FreeAccount = iota + 1
	ProAccount
)

const (
	freeHost = "https://api-free.deepl.com/v2"
	proHost  = "https://api.deepl.com/v2"
)

const (
	textTranslateUri       = "/translate"
	documentTranslateUri   = "/document"
	checkDocumentStatusUri = "/document/%s"
	downloadDocumentsUri   = "/document/%s/result"
	usageUri               = "/usage"
	languagesUri           = "/languages"
	textImprovementUri     = "/write/rephrase"
	listGlossaryPairsUri   = "/glossary-language-pairs"
	createGlossaryUri      = "/glossaries"
	listGlossariesUri      = "/glossaries"
	glossaryDetailsUri     = "/glossaries/%s"
	glossaryEntriesUri     = "/glossaries/%s/entries"
	deleteGlossaryUri      = "/glossaries/%s"
)

var (
	ErrBadRequest          = NewError(400, "Bad request. Please check error message and your parameters.")
	ErrAuthorization       = NewError(401, "Authorization failed. Please supply a valid DeepL-Auth-Key via the Authorization header.")
	ErrForbidden           = NewError(403, "Forbidden. The access to the requested resource is denied, because of insufficient access rights.")
	ErrNotFount            = NewError(404, "The requested resource could not be found.")
	ErrLimit               = NewError(413, "The request size exceeds the limit.")
	ErrLongURL             = NewError(414, "The request URL is too long. You can avoid this error by using a POST request instead of a GET request, and sending the parameters in the HTTP body.")
	ErrNotAccept           = NewError(415, "The requested entries format specified in the Accept header is not supported.")
	ErrManyRequests        = NewError(429, "Too many requests. Please wait and resend your request.")
	ErrQuotaExceeded       = NewError(456, "Quota exceeded. The character limit has been reached.")
	ErrInternal            = NewError(500, "Internal error. ")
	ErrResourceUnavailable = NewError(504, "Resource currently unavailable. Try again later.")
	ErrManyRequests2       = NewError(529, "Too many requests. Please wait and resend your request.")
)

const (
	TagHandlingXML  = "xml"
	TagHandlingHTML = "html"

	FormalityDefault    = "default"
	FormalityMore       = "more"
	FormalityLess       = "less"
	FormalityPreferMore = "prefer_more"
	FormalityPreferLess = "prefer_less"

	SplitSentencesNoSplit                = "0"
	SplitSentencesPunctuationAndNewLines = "1"
	SplitSentencesNoNewLines             = "nonewlines"

	LanguagesTypeSource = "source"
	LanguagesTypeTarget = "target"

	WritingStyleAcademic       = "academic"
	WritingStyleBusiness       = "business"
	WritingStyleCasual         = "casual"
	WritingStyleDefault        = "default"
	WritingStyleSimple         = "simple"
	WritingStylePreferAcademic = "prefer_academic"
	WritingStylePreferBusiness = "prefer_business"
	WritingStylePreferCasual   = "prefer_casual"
	WritingStylePreferSimple   = "prefer_simple"

	ToneDefault            = "default"
	ToneConfident          = "confident"
	ToneDiplomatic         = "diplomatic"
	ToneEnthusiastic       = "enthusiastic"
	ToneFriendly           = "friendly"
	TonePreferConfident    = "prefer_confident"
	TonePreferDiplomatic   = "prefer_diplomatic"
	TonePreferEnthusiastic = "prefer_enthusiastic"
	TonePreferFriendly     = "prefer_friendly"

	DocumentStatusQueued      = "queued"
	DocumentStatusTranslating = "translating"
	DocumentStatusDone        = "done"
	DocumentStatusError       = "error"

	EntriesFormatTSV = "tsv"
	EntriesFormatCSV = "csv"
)
