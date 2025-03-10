# Intellectual_property_mannagement_system(知识产权管理系统)

## 项目结构:

```text
Intellectual_property_mannagement_system
├─cmd               # 应用程序入口
│  └─myapp          # 具体应用
├─config            # 配置文件
├─internal          # 私有代码库
│  ├─api            # Api定义和协议
│  └─services       # 服务文件
├─static            # md文件素材,日志文件存放位置
├─pkg               # 内部库代码
│  ├─models         # 数据模型
│  └─utils          #工具类函数
└─tests             # 测试代码

```

## 专利数据来源

```text
国家知识产权局网站的专利数据库
https://ipdps.cnipa.gov.cn
```

## 参考网站

以下系统仅仅用于学习

```text
https://app.wepatent.cn    # WePatent 管理平台
```

## 问题

```text
修改环境变量要把编译器重启
```

## 第三方包
```text
github.com/gin-gonic/gin          # gin框架
github.com/gorm.io/gorm           # gorm数据库操作
github.com/gorm.io/driver/mysql   # mysql驱动
github.com/spf13/viper            # 配置管理
go.uber.org/zap                   # 日志包
gopkg.in/natefinch/lumberjack.v2  # 日志写入滚动文件，配合zap使用
```