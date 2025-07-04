<div align="center" style="text-align: center">
  <img src="logo.png" width="33%">
</div>

<div align="center">

[![goreportcard](https://goreportcard.com/badge/github.com/falcolee/xutils)](https://goreportcard.com/report/github.com/falcolee/xutils)

</div>

> Go语言工具库

> 集成常用的开发工具包

---

## 📦 目录结构

```mermaid
<!--mermaid-diagram-begin-->
graph TD;
  A[xutils 项目根目录]
  A --> B[xauth]
  A --> C[xcache]
  A --> D[xcli]
  A --> E[xconvert]
  A --> F[xcrypto]
  A --> G[xencoding]
  A --> H[xerror]
  A --> I[xfile]
  A --> J[xfilter]
  A --> K[xgen]
  A --> L[xgit]
  A --> M[xgorm]
  A --> N[xhashring]
  A --> O[xhttp]
  A --> P[ximg]
  A --> Q[xip]
  A --> R[xjson]
  A --> S[xlodash]
  A --> T[xpool]
  A --> U[xregex]
  A --> V[xstruct]
  A --> W[xticker]
  A --> X[xtime]
  A --> Y[xtype]
  A --> Z[xutil]
  A --> AA[xversion]
  B["xauth: JWT认证相关"]
  C["xcache: 缓存接口与实现"]
  D["xcli: 命令行工具与终端交互"]
  E["xconvert: 类型与颜色转换"]
  F["xcrypto: 加解密工具"]
  G["xencoding: 编解码工具"]
  H["xerror: 错误定义与处理"]
  I["xfile: 文件操作"]
  J["xfilter: 过滤器与变量上下文"]
  K["xgen: ID生成等"]
  L["xgit: Git相关操作"]
  M["xgorm: GORM扩展"]
  N["xhashring: 一致性哈希"]
  O["xhttp: HTTP请求工具"]
  P["ximg: 图片处理"]
  Q["xip: IP相关工具"]
  R["xjson: JSON处理"]
  S["xlodash: 集合与切片操作"]
  T["xpool: 协程池"]
  U["xregex: 正则工具"]
  V["xstruct: 结构体工具"]
  W["xticker: 定时器"]
  X["xtime: 时间工具"]
  Y["xtype: 类型转换"]
  Z["xutil: 常用工具函数"]
  AA["xversion: 版本号处理"]
<!--mermaid-diagram-end-->
```

## ✨ 功能模块概述

- **xauth**：JWT签发与解析，简化Token认证流程。
- **xcache**：统一缓存接口，支持多种存储后端。
- **xcli**：命令行交互、终端UI、进度条、表格渲染等。
- **xconvert**：类型转换、颜色转换等常用转换工具。
- **xcrypto**：AES、RSA等加解密算法封装。
- **xencoding**：Base64、URL、JSON等常用编码解码。
- **xerror**：错误码定义、错误处理与调用栈追踪。
- **xfile**：文件读写、类型判断、MD5、JSON序列化等。
- **xfilter**：过滤器、变量上下文、条件表达式。
- **xgen**：ID生成、雪花算法等。
- **xgit**：Git仓库相关操作。
- **xgorm**：GORM数据库扩展，集成缓存。
- **xhashring**：一致性哈希算法。
- **xhttp**：HTTP请求、Header、重试、超时等。
- **ximg**：图片处理、EXIF、合成、圆角等。
- **xip**：IP、MAC、域名等网络工具。
- **xjson**：JSON序列化、反序列化、工具函数。
- **xlodash**：集合、切片、Map等常用操作。
- **xpool**：协程池、任务池。
- **xregex**：正则表达式工具。
- **xstruct**：结构体相关工具。
- **xticker**：定时器。
- **xtime**：时间格式化、定时、超时等。
- **xtype**：类型转换。
- **xutil**：常用工具函数，如三目运算、字符串处理等。
- **xversion**：版本号格式化与比较。

## 🚀 快速上手

以 xauth JWT 为例：

```go
import "github.com/falcolee/xutils/xauth"

// 生成Token
claims := map[string]interface{}{"uid": 123}
token, err := xauth.GetJWT("your-secret", claims, 3600)

// 解析Token
parsed, err := xauth.ParseJWT("your-secret", token)
```

以 xfile 文件操作为例：

```go
import "github.com/falcolee/xutils/xfile"

// 读取JSON文件
var data map[string]interface{}
err := xfile.ReadJSON("data.json", &data)

// 写入JSON文件
err = xfile.WriteJSON("data.json", data)
```

更多模块用法请参考各子目录源码与注释。

---

## 📄 License

MIT License
