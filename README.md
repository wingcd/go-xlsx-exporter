golang编写的将xlsx表文件数据及结构导出工具

### 功能列表

#### 类型支持

- bool

- int

- uint

- int64

- uint64

- float

- double

- string

- 及以上数据类型的数组类型，如bool[],int[],数组通过符号“|”分割，通过两个“|”即“||”可转义此分割符

#### 数据配置

- [x] 支持枚举类型

- [x] 支持自定义类型

- [x] 支持全局定义

- [x] 支持客户端/服务器导出

- [x] 支持注释

- [x] 支持多语言

- [x] 忽略空行/列

- [x] 分表配置支持

#### 导出支持

- [x] 支持.proto文件导出

- [x] 支持序列化为protobuf文件导出

- [x] 支持golang数据结构代码导出

- [x] 支持csharp数据接口代码导出

- [x] 支持多语言导出

- [ ] 支持json结构及数据导出

- [ ] 支持lua结构及数据导出

- [ ] 支持sqlite表结构及数据导出

#### 读取支持

- [x] csharp读取支持：

  - [x] 支持自定义数据读取

  - [x] 支持设置读取

  - [x] 支持表格数据读取

  - [x] 支持转哈希表

- [ ] golang读取支持

  - [ ] 自定义数据读取

  - [ ] 设置读取

  - [ ] 表格数据读取

### 快速开始

#### 使用案例

- 下载发布包

- 复制conf.template.yaml，并改名为conf.yaml

- 命令行运行gxe.exe

- 在gen目录下查看生成数据

#### 项目中使用

- 修改conf.yaml数据，其中：

```YAML
package: "GameData" # 导出代码的包名（命名空间）
pb_bytes_file_ext: ".bytes" # 二进制数据文件后缀名，unity中请使用.bytes
comment_symbol: "#" # 定义表格中的注释符
export_type: 1 # 全局导出类型设置，1-忽略前后端配置,2-仅导出客户端代码/数据,3-仅导出服务端代码/数据
sheets: # 所有表格数据
 -
  id: 1 # 表格编号，用于生成时过滤
  type: "define" # 表格数据类型,包含： define/table/language
  xls_file: 'data/define.xlsx' # xlsx文件路径
  sheet: 'define' # xlsx中的表格名
 - 
  id: 2
  type: "table"
  xls_file: 'data/model.xlsx'
  sheet: 'user'
  type_name: 'User'  
 - 
  id: 2
  type: "table"
  xls_file: 'data/i18n.xlsx'
  sheet: 'location1'
  is_lang: true
  
exports: # 导出任务集合
-
 id: 1 # 编号，用于任务过滤
 type: "proto"  # 任务类型，包含：proto/proto_bytes/golang/csharp或后续支持类型
 path: "./gen/code/all_proto.proto" # 生成文件名，或者导出路径（针对多文件输出）
 sheets: "" # 导出表id集合，默认为导出所有表，逗号分割表示数组，如：1,2,3，横杠分割表示范围,如：1-7，可混用如：1,2-4,6-7
 export_type: 2  # 可单独设置导出类型，覆盖全局设置
-
 id: 2
 type: "golang"
 path: "./gen/code/data_model.pb.go"
 sheets: "1,3-7"
 package: "game_data" # 可单独设置包名，覆盖全局设置，但是会被参数设置的包名覆盖
-
 id: 3
 type: "proto_bytes"
 path: "./gen/data/"

```

- 命令参数

```text
 -cfg string
        设置配置文件 (default "./conf.yaml")
 -cmt string
        设置表格注释符号 (default "#")
 -exports string
        设置需要导出的配置项，默认为空，全部导出, 参考：1,2,5-7
 -ext string
        设置二进制数据文件后缀(unity必须为.bytes) (default ".bytes")
 -h    获取帮助
 -lang
        是否生成语言类型到代码（仅测试用，默认为false）
 -pkg string
        设置导出包名
 -v    获取工具当前版本号
```

- 表格定义

  - 定义表

    此表用于定义数据表可使用的枚举类型，以及结构体及全局配置（可不使用后两种类型）,表结构如下：

    ![](./doc/imgs/define-table.png)

    - 枚举类型

    > enum, 同一枚举类型名一样，类型为空，值为枚举索引值

    ![](./doc/imgs/define-enum.png)

    - 结构体类型

    > struct, 此项预留

    ![](./doc/imgs/define-struct.png)

    - 常量类型

    > const, 用于定义游戏中的一些常量值，常量数据也会导出为二进制文件，与数据文件不同的是，导出的数据不是列表数据，仅为一个消息体（类）的数据

    ![](./doc/imgs/define-const.png)

  - 数据表

    此类型用于定义数据表，前四行用于字段描述，支持使用定义表中定义的枚举类型，导出的二进制数据类型为“类型名_ARRAY”，并在此类型中定义了Items的数组作为此表数据集合，如图：

    ![](./doc/imgs/data-table.png)

    ![](./doc/imgs/data-code.png)

    1. 第一行为字段描述

    2. 第二行为字段类型

    3. 第三行定义此字段导出类型，c表示支持客户端导出，s表示支持服务端导出

    4. 第三行为字段名，生成的代码结构中的名称，在go语言中，首字母将会被大写

  - 语言表

    此表是数据表的一种特例，专门用来存放多语言数据，结构同数据表。但此表会固定生成language.xx.bytes类似的多个数据文件，xx表示字段名的小写，如：language.cn.bytes

    ![](./doc/imgs/define-table.png)

  - 其他

  > 数据表支持注释项，以及过滤空行等

  > 支持所有类型表的分表存储，只需类型名相同即可



#### 使用数据

  现导出的二进制数据为protobuf3.0协议数据，可通过导出的proto文件生成个种语言的代码文件，然后加载此工具导出的数据。当然工具也准备了通用的数据读取类：

- csharp

  将reader/csharp下的DataAccess.cs及I18N.cs拷贝至项目中,unity环境中使用是，请在DataAccess.cs首行加上#define UNITY_ENGINE

1. 初始化

```C#
// 配置二进制数据目录
DataAccess.Initial("./data/");

// 也可使用自定义的数据加载接口进行数据加载, 或自定义的文件名生成器获取数据文件
DataAccess.Initial("./data/", LoadDataHandler, FileNameGenerateHandler);

```

2. 表数据

```C#
var userdata = DataContainer<uint, User>.Instance.Items;
foreach (var item in userdata)
{
    Console.WriteLine($"{item.ID}, {item.Name}, {item.Age}, {item.Head}");
}
```

3. 多语言

  如果使用此工具配置的多语言，则需要I18N.cs文件，否则可不拷贝此文件

```C#
var lanKey = "1";
I18N.SetLanguage("cn");
Console.WriteLine($"\n中文：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");
I18N.SetLanguage("en");
Console.WriteLine($"english：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");
```

4. 配置

```C#
var settings = DataContainer<Settings>.Instance.Data;
Console.WriteLine($"\n配置：maxconn={settings.MAX_CONNECT}, version={settings.VERSION}");
```

5. 读取哈希表

```C#
var userHT = DataContainer<uint, User>.Instance.GetHashtable(1);
Console.WriteLine($"用户哈希值：{userHT["ID"]}, {userHT["Name"]}, {userHT["Age"]}, {userHT["Head"]}");
```

- go

  暂未实现