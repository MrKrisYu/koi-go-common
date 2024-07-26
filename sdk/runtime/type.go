package runtime

import (
	"context"
	"github.com/MrKrisYu/koi-go-common/config"
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/sdk/i18n"
	"github.com/MrKrisYu/koi-go-common/storage"
	"github.com/casbin/casbin/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"net/http"
)

// Runtime 运行时容器的接口能力， 参看 go-admin-core/sdk/runtime包
type Runtime interface {
	// SetTranslator 设置翻译器
	SetTranslator(t i18n.Translator)
	GetTranslator() i18n.Translator

	// SetContext 设置全局上下文
	SetContext(ctx context.Context)
	GetContext() context.Context

	// SetDb 多db设置，⚠️SetDbs不允许并发,可以根据自己的业务，例如app分库、host分库
	SetDb(key string, db *gorm.DB)
	GetDb() map[string]*gorm.DB
	GetDbByKey(key string) *gorm.DB

	// SetCasbin 运行时的权限认证架构casbin的配置
	SetCasbin(key string, enforcer *casbin.SyncedEnforcer)
	GetCasbin() map[string]*casbin.SyncedEnforcer
	GetCasbinKey(key string) *casbin.SyncedEnforcer

	// SetEngine 运行时的web引擎适配
	SetEngine(engine http.Handler)
	GetEngine() http.Handler

	// SetLogger 使用自定义的logger适配器，参考来源go-admin
	SetLogger(logger *logger.Helper)
	GetLogger() *logger.Helper

	// SetCrontab 运行时的定时任务配置
	SetCrontab(key string, crontab *cron.Cron)
	GetCrontab() map[string]*cron.Cron
	GetCrontabKey(key string) *cron.Cron

	// SetMiddleware 运行时的中间件配置
	SetMiddleware(string, interface{})
	GetMiddleware() map[string]interface{}
	GetMiddlewareKey(key string) interface{}

	// SetCacheAdapter 运行时的缓存配置
	SetCacheAdapter(storage.AdapterCache)
	GetCacheAdapter() storage.AdapterCache
	GetCachePrefix(string) storage.AdapterCache

	// SetQueueAdapter 运行时的任务队列配置
	SetQueueAdapter(storage.AdapterQueue)
	GetMemoryQueue(string) storage.AdapterQueue
	GetQueueAdapter() storage.AdapterQueue
	GetQueuePrefix(string) storage.AdapterQueue

	// SetLockerAdapter 运行时的使用的锁配置
	SetLockerAdapter(storage.AdapterLocker)
	GetLockerAdapter() storage.AdapterLocker
	GetLockerPrefix(string) storage.AdapterLocker

	// SetDefaultConfig 运行时的配置信息配置
	SetDefaultConfig(config config.Config)
	GetDefaultConfig() config.Config
	//SetConfigValue(key string, value interface{})
	//GetConfigValue(key string) interface{}
}
