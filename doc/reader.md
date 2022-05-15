#### 读取支持

| 功能\语言                                   | csharp |  ts   |  js   | golang | json  |
| ------------------------------------------- | :----: | :---: | :---: | :----: | :---: |
| 枚举                                        |   √    |   √   |   √   |   √    |   ×   |
| 设置读取                                    |   √    |   √   |   √   |   √    |   √   |
| 数据读取                                    |   √    |   √   |   √   |   √    |   √   |
| 转Lua表                                     |   √    |   ×   |   ×   |   ×    |   ×   |
| 多语言                                      |   √    |   √   |   √   |   √    |   √   |
| 空类型([void](./field_types.md#void))       |   √    |   √   |   √   |   ×    |   ×   |
| 计算类型([?](./field_types.md#计算))        |   √    |   √   |   √   |   ×    |   ×   |
| 自定义别名([别名](./field_types.md#自定义)) |   √    |   √   |   √   |   ×    |   ×   |
| bytes支持                                   |   √    |   √   |   √   |  √    |   ×   |  


### 辅助类
____
> 现导出的二进制数据为protobuf3.0协议数据，可通过导出的proto文件生成个种语言的代码文件，然后加载此工具导出的数据。当然工具也准备了通用的数据读取类：

**[`CSHARP`](#CSHARP)**

  将reader/csharp下的DataAccess.cs及I18N.cs拷贝至项目中,unity环境中使用时，请在DataAccess.cs首行加上#define UNITY_ENGINE, 使用protobuf-net读取protobuf数据  
   
  导出数据模型参考：[csharp demo](../gen/DataMode.cs)

  测试参考：[csharp test](../reader/csharp/test.cs)

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
> JS/TS 使用方式类似

**[`GOLANG`](#GOLANG)**

  拷贝reader/golang目录reader.go即可开始使用  
   
  导出数据模型参考：[golang demo](../gen/DataMode.pb.go)

  测试参考：[golang test](../gen/go_proto_test.go)

  测试参考：[golang test](../gen/go_message_test.go)
  
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