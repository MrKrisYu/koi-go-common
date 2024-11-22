package example

import (
	"github.com/MrKrisYu/koi-go-common/sdk/api/header"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

const (
	LanguageReqQueryKey = "lang"
)

func GetLang(ctx *gin.Context) language.Tag {
	matchedTag := DefaultLanguage
	acceptLanguage := ctx.GetHeader(header.AcceptLanguageFlag)
	if len(acceptLanguage) == 0 {
		acceptLanguage, _ = ctx.GetQuery(LanguageReqQueryKey)
	}
	if len(acceptLanguage) == 0 {
		//fmt.Printf("[Middleware-AcceptLanguage] got lang=%s, match=%s \n\n", acceptLanguage, matchedTag.String())
		return matchedTag
	}
	matcher := language.NewMatcher(AllowedLanguage)
	_, index, confidence := matcher.Match(language.Make(acceptLanguage))
	if confidence != language.No {
		matchedTag = AllowedLanguage[index]
	}
	//fmt.Printf("[Middleware-AcceptLanguage] got lang=%s, match=%s \n\n", acceptLanguage, matchedTag.String())
	return matchedTag
}

func AcceptLanguage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang := GetLang(ctx)
		ctx.Set(header.AcceptLanguageFlag, lang)
		ctx.Next()
	}
}
