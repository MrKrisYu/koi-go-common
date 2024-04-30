package source

import (
	"fmt"
	"os"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	str := "${DB_HOST:}"

	// 定义正则表达式模式
	pattern := `\$\{([a-zA-Z0-9_]+):?([a-zA-Z0-9_]*)\}`

	// 编译正则表达式
	regex := regexp.MustCompile(pattern)

	// 查找匹配项
	matches := regex.FindStringSubmatch(str)
	fmt.Printf("matches = %q", matches)
	if len(matches) == 3 {
		envVar := matches[1]
		defaultValue := matches[2]

		// 获取环境变量的值
		value := os.Getenv(envVar)

		// 如果环境变量不存在或为空值，则使用默认值
		if value == "" {
			value = defaultValue
		}

		fmt.Printf("环境变量: %s, 值: %s\n", envVar, value)
	} else {
		fmt.Println("未找到匹配项")
	}
}
