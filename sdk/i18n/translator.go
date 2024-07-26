package i18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/fs"
)

type Translator interface {
	// RegisterLocalizer 注册一个特定语言的Localizer
	RegisterLocalizer(fsys fs.FS, path string, lang language.Tag) error

	// GetLocalizer 获取一个特定语言的Localizer
	GetLocalizer(lang language.Tag) *i18n.Localizer

	// HasLanguageTag 判断是否注册了该语言
	HasLanguageTag(tag language.Tag) bool

	// Tr 翻译纯文本消息
	Tr(lang language.Tag, message Message) string

	// TrWithData 翻译带数据消息
	TrWithData(lang language.Tag, message Message) string
}

// Message 国际化消息
type Message struct {
	ID             string `json:"ID"`             // 消息ID
	DefaultMessage string `json:"defaultMessage"` // 默认消息，翻译失败时，使用
	Args           any    `json:"args"`           // 消息参数
}
