package reader

import (
	"os"
	"regexp"
)

func ReplaceEnvVars(raw []byte) ([]byte, error) {
	// 定义正则表达式模式
	pattern := `\$\{([a-zA-Z0-9_]+):?([a-zA-Z0-9_]*)\}`
	re := regexp.MustCompile(pattern)
	if re.Match(raw) {
		dataS := string(raw)
		res := re.ReplaceAllStringFunc(dataS, replaceEnvVars)
		return []byte(res), nil
	} else {
		return raw, nil
	}
}

func replaceEnvVars(element string) string {
	// 最后替换的接爱国
	value := ""
	// 定义正则表达式模式
	pattern := `\$\{([a-zA-Z0-9_]+):?([a-zA-Z0-9_]*)\}`
	// 编译正则表达式
	regex := regexp.MustCompile(pattern)
	// 查找匹配项
	matches := regex.FindStringSubmatch(element)
	if len(matches) == 3 {
		envVar := matches[1]
		defaultValue := matches[2]
		// 获取环境变量的值
		value = os.Getenv(envVar)
		// 如果环境变量不存在或为空值，则使用默认值
		if value == "" {
			value = defaultValue
		}
	}
	return value
}
