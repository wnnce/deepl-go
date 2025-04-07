package deepl

import (
	"bytes"
	"sync"
)

var (
	bufferPool            *sync.Pool
	textParamsPool        *sync.Pool
	documentParamsPool    *sync.Pool
	improvementParamsPool *sync.Pool
	glossaryParamsPool    *sync.Pool
)

func init() {
	bufferPool = &sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}
	textParamsPool = &sync.Pool{
		New: func() any {
			return &TextTranslateParams{}
		},
	}
	documentParamsPool = &sync.Pool{
		New: func() any {
			return &DocumentTranslateParams{}
		},
	}
	improvementParamsPool = &sync.Pool{
		New: func() any {
			return &TextImprovementParams{}
		},
	}
	glossaryParamsPool = &sync.Pool{
		New: func() any {
			return &CreateGlossaryParams{}
		},
	}
}

func AcquireTextTranslateParams() *TextTranslateParams {
	value := textParamsPool.Get()
	params, ok := value.(*TextTranslateParams)
	if !ok {
		params = &TextTranslateParams{}
	}
	return params
}

func AcquireDocumentTranslateParams() *DocumentTranslateParams {
	value := documentParamsPool.Get()
	params, ok := value.(*DocumentTranslateParams)
	if !ok {
		params = &DocumentTranslateParams{}
	}
	return params
}

func AcquireTextImprovementParams() *TextImprovementParams {
	value := improvementParamsPool.Get()
	params, ok := value.(*TextImprovementParams)
	if !ok {
		return &TextImprovementParams{}
	}
	return params
}

func AcquireCreateGlossaryParams() *CreateGlossaryParams {
	value := glossaryParamsPool.Get()
	params, ok := value.(*CreateGlossaryParams)
	if !ok {
		return &CreateGlossaryParams{}
	}
	return params
}

func RecycleParams(params Recyclable) {
	params.recycle()
	switch params.(type) {
	case *TextTranslateParams:
		textParamsPool.Put(params)
	case *TextImprovementParams:
		improvementParamsPool.Put(params)
	case *DocumentTranslateParams:
		documentParamsPool.Put(params)
	case *CreateGlossaryParams:
		glossaryParamsPool.Put(params)
	}
}

func recycleBuffer(buffer *bytes.Buffer) {
	buffer.Reset()
	bufferPool.Put(buffer)
}
