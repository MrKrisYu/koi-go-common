# koi-go-common

> - @author: krisyu291157877@gmail.com
> - @description: Go开发常用工具类, 参考[go-admin-core](https://github.com/MrKrisYu/koi-go-common/)设计
> - @repository:  [GitHub](https://github.com/MrKrisYu/koi-go-common)


## 功能

- [x] controller&service组件
- [x] log组件
- [x] config组件





### 实现目标

1. 依赖注入的形式，引入则使用，不引入则不使用
2. 约定大于配置的原则
3. 容器的概念，一切配置的作用域均在容器之下，最大上下文：application









#### 日志设计

- [x] 提供一个日志使用泛型: 基本的日志等级设置，不同等级的打印；
- [x] 基于SPI的思路，来引入不同的实现： 想用使用不同的日志实现库，就要引入实现了指定接口的日志库
- [x] 日志的属性配置： 提供默认值与选配值，选配置看用户配置的属性



#### 自动配置的设计

- [ ] 初始化配置： 确定文件来源，确定文件来源解析器，确定完成配置加载后的回调函数（用于执行配置对应组件的初始化函数）
- [ ] 解析器除了解析内置的属性外，还要把用户自定义的属性读取到运行时的config结构体中，用户通过键即可获取对应的值
- [ ] 
