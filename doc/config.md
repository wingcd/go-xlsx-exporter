#### 配置文件

表格配置规则，以及数据导出规则需要修改配置表，但一般只需要在配置一次，除非一直有表格增减

- 修改conf.yaml数据，其中相关参数解释如下：

```YAML
package: "GameData" # 导出代码的包名（命名空间）
pb_bytes_file_ext: ".bytes" # 二进制数据文件后缀名，unity中请使用.bytes
comment_symbol: "#" # 定义表格中的注释符
export_type: 1 # 全局导出类型设置，1-导出前后端代码/数据,2-仅导出客户端代码/数据,3-仅导出服务端代码/数据,4-忽略(表格式无需添加前后端导出配置行)；对消息类型无效 
array_split_char: "," #默认数组分割符号
pause_on_end: false # 运行完毕后是否暂停
strict_mode: true # 是否严格模式,如：int配置为空时，严格模式将会报错，非严格模式默认为0
rules: #数据规则，可在字段类型后加入规则id,在数据导出时进行规则检测，严格模式会中断输出，否则只进行日志提示
-
  id: 1              # 规则id，在类型表达式中使用
  rule: '\w+?\.png'  # 规则正则表达式
  desc: '图片'        # 描述
  disable: false     # 是否关闭此规则
-
  id: 2
  rule: '\d+'
  desc: '数字'
  disable: false
includes: # 所有表格数据
 -
  id: 1 # 表格编号，用于生成时过滤
  type: "define" # 表格数据类型,包含： define/table/language
  file: 'data/define.xlsx' # xlsx文件路径
  sheet: 'define' # xlsx中的表格名
  transpose: true  # 是否转置表格(行列翻转)
 - 
  id: 2
  type: "table"
  file: 'data/model.xlsx'
  sheet: 'user'
  type_name: 'User'  
 - 
  id: 2
  type: "table"
  file: 'data/i18n.xlsx'
  sheet: 'location1'
  is_lang: true
 - 
  id: 3 # 消息类型预定义，文件为xml
  type: define
  file: 'data/define.xml'
 - 
  id: 4 # 消息文件，不存在前后端差异
  type: message
  file: 'data/message.xml'
  
exports: # 导出任务集合
-
 id: 1 # 编号，用于任务过滤
 type: "proto"  # 任务类型，包含：proto/proto_bytes/golang/csharp或后续支持类型
 path: "./gen/code/all_proto.proto" # 生成文件名，或者导出路径（针对多文件输出）
 includes: "" # 导出表id集合，默认为导出所有表，逗号分割表示数组，如：1,2,3，横杠分割表示范围,如：1-7，可混用如：1,2-4,6-7
 excludes: "" # 排除id结合，同includes
 export_type: 2  # 可单独设置导出类型，覆盖全局设置（忽略导出除外，只认全局配置） 
 template: "template/proto.gtpl" # 用于自定义模板，默认为空，使用系统指定模板，否则使用指定模板输出
-
 id: 2
 type: "golang"
 path: "./gen/code/data_model.pb.go"
 includes: "1,3-7"
 package: "game_data" # 可单独设置包名，覆盖全局设置，但是会被参数设置的包名覆盖
-
 id: 3
 type: "csharp"
 path: "./gen/code/DataModel.cs"
 includes: ""
 package: "Cfg"
-
 id: 4
 type: "proto_bytes"
 path: "./gen/data/"
-
 id: 5
 type: "charset"
 path: "./gen/data/lang.txt"
-
 id: 6
 type: "js,dts" #部分任务可同时进行，避免多次解析表格，多类型导出按","分割，且输出路径必须与类型数量相同
 path: "./gen/code/data_mode.js,./gen/code/data_mode.d.ts"
 imports:  #当使用了自定义类型转换，可在生成代码中加入引用，防止生成代码错误
  - "import UserData from './userdata'"
  - "import XXX from './xxx'"
-
 id: 7
 type: "ts"
 path: "./gen/code/data_mode.ts"

```