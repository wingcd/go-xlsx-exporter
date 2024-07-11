### 导出支持

- [x] [](#proto)支持.proto文件导出(type:[proto](../gen/all.proto))

- [x] [](#proto_bytes)支持序列化为protobuf文件导出(type:[proto_bytes](../gen/user.bytes))

- [x] [](#golang)支持golang数据结构代码导出(type:[golang](../gen/DataMode.pb.go))

- [x] [](#csharp)支持csharp数据接口代码导出(type:[csharp](../gen/DataMode.cs))

- [x] [](#js)支持javascript数据接口代码导出(type:[js](../gen/data_mode.js))

- [x] [](#dts)支持typescript接口代码导出, 用于配合javascript代码(type:[dts](../gen/data_mode.d.ts))

- [x] [](#ts)支持typescript数据接口代码导出(type:[ts](../gen/data_mode.ts))

- [x] [](#json)支持json数据导出(type:[json](../gen/User.json))

- [ ] [](#lua)支持lua结构及数据导出

- [x] [](#charset)支持多语言表使用的文字导出为文本文件，用于生成字符集(type:charset)

- [x] [](#i18n)支持多语言数据导出

- [ ] [](#sqlite)支持sqlite表结构及数据导出

- [x] [](#regex)支持列正则检查
  
- [x] [](#message)支持消息导出（无需单独导出类型，按相关语言导出即可）[案例](../gen/message.go)

- [x] [](#channel)支持列级别的渠道数据配置，通过配置或者命令行导出不同渠道数据

- [x] [](#custom)支持自定义导出脚本，通过编写lua生成导出代码、数据