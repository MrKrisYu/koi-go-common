package example

import (
	"encoding/json"
	"fmt"
	i18n2 "github.com/MrKrisYu/koi-go-common/sdk/i18n"
	"github.com/MrKrisYu/koi-go-common/sdk/i18n/example/resource"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
	"slices"
	"strings"
)

var (
	DefaultLanguage = language.AmericanEnglish

	AllowedLanguage = []language.Tag{
		DefaultLanguage,
		language.SimplifiedChinese,
		language.German,
		language.Spanish,
	}
)

func initialI18nBundle(bundle *i18n.Bundle, matcher language.Matcher) (map[language.Tag]*i18n.Localizer, []language.Tag, error) {
	var localizers = make(map[language.Tag]*i18n.Localizer)
	var loadTags []language.Tag
	entries, err := resource.I18n.ReadDir(".")
	if err != nil {
		return nil, loadTags, err
	}

	for _, entry := range entries {
		fileName := entry.Name()
		langStr := strings.TrimRight(fileName, ".json")
		// 检查是否为允许的语言
		_, i, c := matcher.Match(language.Make(langStr))
		if c == language.No {
			//fmt.Printf("[initialI18nBundle] load message file failed, langStr = %s does not match any of the allowed tags(%+v) \n\n", langStr, AllowedLanguage)
			continue
		}
		bindTag := AllowedLanguage[i]
		if _, ok := localizers[bindTag]; ok {
			//fmt.Printf("[initialI18nBundle] localizer of tag=%s already exists, langStr = %s\n\n", bindTag.String(), langStr)
			continue
		}
		// 加载语言资源
		_, err = bundle.LoadMessageFileFS(resource.I18n, fileName)
		if err != nil {
			return nil, loadTags, err
		}
		// 初始化Localizer
		localizers[bindTag] = i18n.NewLocalizer(bundle, bindTag.String())
		loadTags = append(loadTags, bindTag)
	}
	return localizers, loadTags, nil
}

type DefaultTranslator struct {
	bundle         *i18n.Bundle
	allowedMatcher language.Matcher // 允许注册语言的匹配器
	loadedTags     []language.Tag   // 已加载的语言列表
	localizers     map[language.Tag]*i18n.Localizer
}

func NewDefaultTranslator(defaultLang language.Tag, allowedLangs []language.Tag) *DefaultTranslator {
	internalMatcher := language.NewMatcher(allowedLangs)
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	localizers, loadedTags, err := initialI18nBundle(bundle, internalMatcher)
	if err != nil {
		panic(err)
	}
	return &DefaultTranslator{
		localizers:     localizers,
		loadedTags:     loadedTags,
		allowedMatcher: internalMatcher,
		bundle:         bundle,
	}
}

func (d *DefaultTranslator) HasLanguageTag(tag language.Tag) bool {
	return slices.Contains(d.loadedTags, tag)
}
func (d *DefaultTranslator) RegisterLocalizer(fsys fs.FS, path string, lang language.Tag) error {
	// 检查是否为允许的语言
	_, i, c := d.allowedMatcher.Match(lang)
	if c == language.No {
		return fmt.Errorf("[RegisterLocalizerFS] load message file failed, langStr = %s does not match any of the allowed tags(%+v) \n\n",
			lang.String(), AllowedLanguage)
	}
	bindTag := AllowedLanguage[i]
	if _, ok := d.localizers[bindTag]; ok {
		return fmt.Errorf("[RegisterLocalizerFS] localizer of tag=%s already exists, langStr = %s\n\n",
			bindTag.String(), lang.String())
	}
	// 加载语言资源
	_, err := d.bundle.LoadMessageFileFS(fsys, path)
	if err != nil {
		return err
	}
	// 初始化Localizer
	d.localizers[bindTag] = i18n.NewLocalizer(d.bundle, bindTag.String())
	d.loadedTags = append(d.loadedTags, bindTag)
	return nil
}
func (d *DefaultTranslator) GetLocalizer(lang language.Tag) *i18n.Localizer {
	matchedTag, exist := d.getBestMatchedLang(lang)
	if !exist {
		return nil
	}
	return d.localizers[matchedTag]
}
func (d *DefaultTranslator) Tr(lang language.Tag, message i18n2.Message) string {
	localizer := d.GetLocalizer(lang)
	if localizer == nil {
		//fmt.Printf("[Tr] %s does not match any locaizer, using default lang:%s \n", lang.String(), DefaultLanguage.String())
		localizer = d.GetLocalizer(DefaultLanguage)
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: message.ID,
	})
	if err != nil {
		fmt.Printf("[Tr] translator error: %v \n", err)
		return message.DefaultMessage
	}
	return msg
}
func (d *DefaultTranslator) TrWithData(lang language.Tag, message i18n2.Message) string {
	localizer := d.GetLocalizer(lang)
	if localizer == nil {
		//fmt.Printf("%s does not match any locaizer, using default lang:%s \n", lang.String(), DefaultLanguage.String())
		localizer = d.GetLocalizer(DefaultLanguage)
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    message.ID,
		TemplateData: message.Args,
	})
	if err != nil {
		fmt.Printf("[Tr] translator error: %v \n", err)
		return message.DefaultMessage
	}
	return msg
}

// getBestMatchedLang 从已加载的语言资源中获取最佳匹配的语言
// 返回值
// tag 匹配的语言
// exist 是否存在匹配的语言， true-存在；false-不存在
func (d *DefaultTranslator) getBestMatchedLang(lang language.Tag) (tag language.Tag, exist bool) {
	// 仅匹配已加载的语言
	matcher := language.NewMatcher(d.loadedTags)
	_, index, confidence := matcher.Match(lang)
	if confidence == language.No {
		return tag, false
	}
	return d.loadedTags[index], true
}
