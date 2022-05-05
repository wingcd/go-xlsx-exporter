golang编写的将xlsx表文件数据及结构导出工具

### 功能列表

#### 类型支持

数据类型表达式：`(基础类型|void)[[数组分割符]][[?]自定义类型][<规则ID>]`,其中技基础数据必须包含，实例如下如： 
  string,int,int32,void,string[,],string[,]?UserInfo,string?UserInfo,string[,]?UserInfo<1>等...
  一般情况，如果没有特殊需求，基本类型即可满足要求

- bool

- int

- uint

- int64

- uint64

- float

- double

- string

- bytes 用于配置16进制数据，部分语言支持，不支持语言将按字符串形式表达

- void   此类型不会单独生成字段，且无需配置数值，只会在代码层生成Get函数，比如：size字段，将csharp中将生成public object GetSize(),获取值为用户注册转换函数的返回值，并缓存起来；可用于解决多字段组合问题，或者需要根据此列表生成的对象

- 及以上数据类型的数组类型，如bool[],int[],数组通过自定义分割符（默认配置为‘|’）分割，通过两个分割符（如：‘||’）可转义此分割符

- 通过在可单独设置列的分隔符，如bool[,],即使用逗号‘，’分割此列

- 支持在类型后添加‘?’, 如string?，将在生成除此字段的属性值外，多一个Get函数，类似void生成，用于解决此字段数据转特殊对象的应用

- 如果问号后面带上自定义类型，可在生成Get函数时定义返回类型，减少类型转换，如果有自定义类型，一般会在生成配置中`添加引用`，或者直接修改相关代码生成模板，在`模板中添加引用`

- 如果使用规则，需要在最后添加<规则id>,如<1>标明使用id=1的规则对列数据进行检测，且检测一般发生在数据生成时

#### 数据配置

- [x] 支持枚举类型

- [x] 支持自定义类型

- [x] 支持全局定义

- [x] 支持客户端/服务器导出

- [x] 支持注释

- [x] 支持多语言

- [x] 忽略空行/列

- [x] 分表配置支持
  
- [x] 支持消息协议配置(xml格式， type=[define/message])

#### 导出支持

- [x] 支持.proto文件导出(type:proto)

- [x] 支持序列化为protobuf文件导出(type:proto_bytes)

- [x] 支持golang数据结构代码导出(type:golang)

- [x] 支持csharp数据接口代码导出(type:csharp)

- [x] 支持javascript数据接口代码导出(type:js)

- [x] 支持typescript接口代码导出, 用于配合javascript代码(type:dts)

- [x] 支持typescript数据接口代码导出(type:ts)

- [x] 支持json数据导出(type:json)

- [ ] 支持lua结构及数据导出

- [x] 支持多语言表使用的文字导出为文本文件，用于生成字符集(type:charset)

- [x] 支持多语言数据导出

- [ ] 支持sqlite表结构及数据导出

- [x] 支持列正则检查
  
- [x] 支持消息导出（无需单独导出类型，按相关语言导出即可）

#### 读取支持

- [x] csharp读取支持：

  - [x] 支持设置读取

  - [x] 支持表格数据读取

  - [x] 支持转Lua表
  
  - [x] 多语言读取 

  - [x] 空类型(void)配置与读取

  - [x] 计算类型(?)配置与读取

  - [x] ?配置增加指定类型

- [x] golang读取支持

  - [x] 设置读取

  - [x] 表格数据读取
  
  - [x] 多语言读取 

### 快速开始

#### 使用案例

- 下载发布包

- 复制conf.template.yaml，并改名为conf.yaml

- 命令行运行gxe

- 在gen目录下查看生成数据

#### 项目中使用

- 修改conf.yaml数据，其中：

```YAML
package: "GameData" # 导出代码的包名（命名空间）
pb_bytes_file_ext: ".bytes" # 二进制数据文件后缀名，unity中请使用.bytes
comment_symbol: "#" # 定义表格中的注释符
export_type: 1 # 全局导出类型设置，1-导出前后端代码/数据,2-仅导出客户端代码/数据,3-仅导出服务端代码/数据,4-忽略(表格式无需添加前后端导出配置行)；对消息类型无效 
array_split_char: "|" #默认数组分割符号
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

- 命令参数

```text
  -cfg string
        设置配置文件 (default "./conf.yaml")
  -cmt string
        设置表格注释符号 (default "#")
  -exports string
        设置需要导出的配置项，默认为空，全部导出, 参考：1,2,5-7
  -ext string
        设置二进制数据文件后缀(unity必须为.bytes)
  -h    获取帮助
  -lang
        是否生成语言类型到代码（仅测试用，默认为false）
  -pkg string
        设置导出包名
  -silence
        是否静默执行（默认为false）
  -v    获取工具当前版本号
```

- [表格定义](#表格定义)

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

    此类型用于定义数据表，前四行(忽略前后端导出配置时为三行)用于字段描述，支持使用定义表中定义的枚举类型，导出的二进制数据类型为“类型名_ARRAY”，并在此类型中定义了Items的数组作为此表数据集合，如图：

    ![](./doc/imgs/data-table.png)

    ![](./doc/imgs/data-code.png)

    1. 第一行为字段描述

    2. 第二行为字段类型

    3. 第三行定义此字段为前后端导出类型(可忽略)，c表示支持客户端导出，s表示支持服务端导出

    4. 第四行为字段名(忽略前后端配置时为第三行)，生成的代码结构中的名称，在go语言中，首字母将会被大写

  - 语言表

    此表是数据表的一种特例，专门用来存放多语言数据，结构同数据表。但此表会固定生成language.xx.bytes类似的多个数据文件，xx表示字段名的小写，如：language.cn.bytes

    ![](./doc/imgs/define-table.png)

  - 其他

  > 数据表支持注释项，以及过滤空行等

  > 支持所有类型表的分表存储，只需类型名相同即可



#### 使用数据

  现导出的二进制数据为protobuf3.0协议数据，可通过导出的proto文件生成个种语言的代码文件，然后加载此工具导出的数据。当然工具也准备了通用的数据读取类：

- csharp

  将reader/csharp下的DataAccess.cs及I18N.cs拷贝至项目中,unity环境中使用时，请在DataAccess.cs首行加上#define UNITY_ENGINE, 使用protobuf-net读取protobuf数据  
   
  导出数据模型参考：[csharp demo](./gen/DataMode.cs)

  测试参考：[csharp test](./reader/csharp/test.cs)

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

3. 多语言（可选）

  不需要多语言时，无需拷贝I18N.cs文件

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

5. 读取LuaTable

```C#
var userData = DataContainer<uint, User>.Instance.GetTable(1);
Console.WriteLine($"用户数据：{userData["ID"]}, {userData["Name"]}, {userData["Age"]}, {userData["Head"]}");
```

6. 计算数据类型

注册计算类型转换，通过表类型与字段名，计算并缓存此值，可用来配置复杂结构类型。

```C#
DataAccess.DataConvertHandler = (item, field, data) =>
  {
      if (item is XXX1Cfg)
      {
          var dt = item as XXX1Cfg;
          switch (field)
          {
              case "Merge":
                  return MultipleWeightCalculator.Parse(data?.ToString());
              case "Size":
                  return new Vector2Int(dt.Width, dt.Height);
          }
      }
      else if(item is XXX2Cfg)
      {
          var dt = item as XXX2Cfg;
          switch (field)
          {
              case "ClickProduce":
              case "ClickKindProduct":
              case "ClickLimit":
                  return WeightCalculator.Parse(data?.ToString());
          
          }
      }

      return data;
  };
```

- go

  拷贝reader/golang目录reader.go即可开始使用  
   
  导出数据模型参考：[golang demo](./gen/DataMode.pb.go)

  测试参考：[golang test](./gen/go_proto_test.go)

  测试参考：[golang test](./gen/go_message_test.go)
  
  ``` golang
  // 项目中导入模块
  import (
    gxe "xxx/go_xlsx_exporter"
  )
  ```

1. 初始化

   ``` golang
   // 配置二进制数据目录，以及默认表索引名
   gxe.Initial("./bytes", "ID")

   // 当不使用默认索引名时，请提前注册表
   gxe.RegistDataTable("Index", reflect.TypeOf(User{}))
   ```
2. 表数据
   ``` golang
    // 获取数据表
    var userTable = gxe.GetDataTable(reflect.TypeOf(User{}))
    // 获取所有数据项
    var users = userTable.Items()
    for _, dt := range users {
      var user = dt.(*User)
      fmt.Printf("Name:%s, Age:%v, Sex: %v \n", user.Name, user.Age, user.Sex)
    }
    // 根据索引获取数据，索引都转化为了字符串类型
    var user = userTable.GetItem("1").(*User)
    fmt.Printf("\nByMap ID=1 Name:%s, Age:%v, Sex: %v \n\n", user.Name, user.Age, user.Sex)
   ```
3. 多语言(可选)
  
  不需要多语言时，无需拷贝语言相关代码文件
   ``` golang
    // 设置当前语言为中文
    gxe.SetLanguage(gxe.DefaultIndexKey, "cn") 
    fmt.Printf("中文 cn id=1, text=%v \n", gxe.Translate("1"))
   
    // 设置当前语言为英文
    gxe.SetLanguage(gxe.DefaultIndexKey, "en")
    fmt.Printf("English en id=1, text=%v \n", gxe.Translate("1"))
   ```
4. 配置
   ``` golang
    var dt = gxe.GetDataItem(reflect.TypeOf(Settings{}))
    var settings = dt.Item().(*Settings)
    fmt.Printf("settings version=%v,maxconn=%v\n\n", settings.VERSION,    settings.MAX_CONNECT)
   ```

#### 通信协议导出
  支持网络协议(protobuf3)导出
  1. define文件定义(后缀.xml)，此文件可定义`枚举类型`, 方便后续使用
  ```xml
  <define>
    <enum name="EMsgType">
        <field name="JSON" value="1" desc="用户ID"/>
        <field name="XML"/>
    </enum>
</define>
  ```

  2. 消息文件定义(后缀.xml)，此文件定义用于网络传输的消息，可以使用预定义的枚举类型，以及此文件中的所有类型，当不定义消息id时，将不会对此消息进行注册
 ``` xml
 <proto>
    <message name="MessageWrapper">
        <field name="id" type="int32"/>
        <field name="data" type="bytes"/>
    </message>
    <message id="10001" name="C2S_GetPlayerInfo">
        <field name="id" alias="用户ID" type="int32" desc="用户ID"/>
        <field name="name" type="string"/>
    </message>
    <message id="10002" name="S2C_GetPlayerInfo">
        <field name="id" type="int32" desc="用户ID"/>
        <field name="name" type="string"/>      
        <field name="type" type="EMsgType"/>  
        <field name="items" type="Item[]"/>
    </message>
    <message name="Item">
        <field name="id" type="int32" desc="ID"/>
        <field name="name" type="string"/>
    </message>
  </proto>
 ```
 3. 消息使用，可注册一个消息封装体，用来二次包装传送消息，方便拿取id后进行消息创建，如golang中通过id和二进制数据创建消息
 ``` golang
err, dt := LoadMessage(10001, bytes)
 ```
   