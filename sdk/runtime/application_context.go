package runtime

import (
	"context"
	"github.com/MrKrisYu/koi-go-common/config"
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/storage"
	"github.com/MrKrisYu/koi-go-common/storage/queue"
	"github.com/casbin/casbin/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

type ApplicationContext struct {
	Context     context.Context
	dbs         map[string]*gorm.DB
	casbins     map[string]*casbin.SyncedEnforcer
	engine      http.Handler
	crontab     map[string]*cron.Cron
	mux         sync.RWMutex
	middlewares map[string]interface{}
	cache       storage.AdapterCache
	queue       storage.AdapterQueue
	memoryQueue storage.AdapterQueue
	locker      storage.AdapterLocker
	//configMap     map[string]interface{} // 系统参数
	defaultConfig config.Config // 系统参数
	logger        *logger.Helper
}

// NewConfig 默认值
func NewConfig() *ApplicationContext {
	return &ApplicationContext{
		Context:     context.Background(),
		dbs:         make(map[string]*gorm.DB),
		casbins:     make(map[string]*casbin.SyncedEnforcer),
		crontab:     make(map[string]*cron.Cron),
		middlewares: make(map[string]interface{}),
		memoryQueue: queue.NewMemory(10000),
		//configMap:   make(map[string]interface{}),
	}
}

// SetContext 设置全局上下文
func (e *ApplicationContext) SetContext(ctx context.Context) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if ctx != nil {
		e.Context = ctx
	}
}

// GetContext 获取全局上下文
func (e *ApplicationContext) GetContext() context.Context {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.Context == nil {
		e.Context = context.Background()
	}
	return e.Context
}

// SetDb 设置对应key的db
func (e *ApplicationContext) SetDb(key string, db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.dbs[key] = db
}

// GetDb 获取所有map里的db数据
func (e *ApplicationContext) GetDb() map[string]*gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.dbs
}

// GetDbByKey 根据key获取db
func (e *ApplicationContext) GetDbByKey(key string) *gorm.DB {
	e.mux.Lock()
	defer e.mux.Unlock()
	if db, ok := e.dbs["*"]; ok {
		return db
	}
	return e.dbs[key]
}

func (e *ApplicationContext) SetCasbin(key string, enforcer *casbin.SyncedEnforcer) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.casbins[key] = enforcer
}

func (e *ApplicationContext) GetCasbin() map[string]*casbin.SyncedEnforcer {
	return e.casbins
}

// GetCasbinKey 根据key获取casbin
func (e *ApplicationContext) GetCasbinKey(key string) *casbin.SyncedEnforcer {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e, ok := e.casbins["*"]; ok {
		return e
	}
	return e.casbins[key]
}

// SetEngine 设置路由引擎
func (e *ApplicationContext) SetEngine(engine http.Handler) {
	e.engine = engine
}

// GetEngine 获取路由引擎
func (e *ApplicationContext) GetEngine() http.Handler {
	return e.engine
}

// SetLogger 设置日志组件
func (e *ApplicationContext) SetLogger(l *logger.Helper) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.logger = l
}

// GetLogger 获取日志组件
func (e *ApplicationContext) GetLogger() *logger.Helper {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.logger
}

// SetCrontab 设置对应key的crontab
func (e *ApplicationContext) SetCrontab(key string, crontab *cron.Cron) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.crontab[key] = crontab
}

// GetCrontab 获取所有map里的crontab数据
func (e *ApplicationContext) GetCrontab() map[string]*cron.Cron {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.crontab
}

// GetCrontabKey 根据key获取crontab
func (e *ApplicationContext) GetCrontabKey(key string) *cron.Cron {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e, ok := e.crontab["*"]; ok {
		return e
	}
	return e.crontab[key]
}

// SetMiddleware 设置中间件
func (e *ApplicationContext) SetMiddleware(key string, middleware interface{}) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.middlewares[key] = middleware
}

// GetMiddleware 获取所有中间件
func (e *ApplicationContext) GetMiddleware() map[string]interface{} {
	return e.middlewares
}

// GetMiddlewareKey 获取对应key的中间件
func (e *ApplicationContext) GetMiddlewareKey(key string) interface{} {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.middlewares[key]
}

// SetCacheAdapter 设置缓存
func (e *ApplicationContext) SetCacheAdapter(c storage.AdapterCache) {
	e.cache = c
}

// GetCacheAdapter 获取缓存
func (e *ApplicationContext) GetCacheAdapter() storage.AdapterCache {
	return NewCache("", e.cache, "")
}

// GetCachePrefix 获取带租户标记的cache
func (e *ApplicationContext) GetCachePrefix(key string) storage.AdapterCache {
	return NewCache(key, e.cache, "")
}

// GetMemoryQueue 获取新的消息队列
func (e *ApplicationContext) GetMemoryQueue(prefix string) storage.AdapterQueue {
	return NewQueue(prefix, e.memoryQueue)
}

// SetQueueAdapter 设置队列适配器
func (e *ApplicationContext) SetQueueAdapter(c storage.AdapterQueue) {
	e.queue = c
}

// GetQueueAdapter 获取队列适配器
func (e *ApplicationContext) GetQueueAdapter() storage.AdapterQueue {
	return NewQueue("", e.queue)
}

// GetQueuePrefix 获取带租户标记的queue
func (e *ApplicationContext) GetQueuePrefix(key string) storage.AdapterQueue {
	return NewQueue(key, e.queue)
}

// SetLockerAdapter 设置分布式锁
func (e *ApplicationContext) SetLockerAdapter(c storage.AdapterLocker) {
	e.locker = c
}

// GetLockerAdapter 获取分布式锁
func (e *ApplicationContext) GetLockerAdapter() storage.AdapterLocker {
	return NewLocker("", e.locker)
}

func (e *ApplicationContext) GetLockerPrefix(key string) storage.AdapterLocker {
	return NewLocker(key, e.locker)
}

// GetStreamMessage 获取队列需要用的message
func (e *ApplicationContext) GetStreamMessage(id, stream string, value map[string]interface{}) (storage.Messager, error) {
	message := &queue.Message{}
	message.SetID(id)
	message.SetStream(stream)
	message.SetValues(value)
	return message, nil
}

// SetDefaultConfig 设置默认的config
func (e *ApplicationContext) SetDefaultConfig(config config.Config) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.defaultConfig = config
	//e.configMap = e.defaultConfig.Map()
}

// GetDefaultConfig 获取默认的config
func (e *ApplicationContext) GetDefaultConfig() config.Config {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.defaultConfig
}

//// SetConfigValue 设置对应key的config
//func (e *ApplicationContext) SetConfigValue(key string, value interface{}) {
//	e.mux.Lock()
//	defer e.mux.Unlock()
//	e.configMap[key] = value
//}
//
//// GetConfigValue 获取对应key的config
//func (e *ApplicationContext) GetConfigValue(key string) interface{} {
//	e.mux.Lock()
//	defer e.mux.Unlock()
//	return e.configMap[key]
//}
