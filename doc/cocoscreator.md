### 如何在cocoscreator中使用

[CCC Demo](./demo/GoXlsxExportDemo-ccc.rar)

1. 下载并解压文件gxe.zip
2. 修改conf.yaml文件，定义项目中使用的表格

``` yaml 
package: "GameConfig"
pb_bytes_file_ext: ".bin" #unity must use .bytes extension, cocoscreator use .bin
comment_symbol: "#"
export_type: 4 #全局导出类型设置，1-all,2-client,3-server,4-ignore-不支持此行
array_split_char: "," #默认数组分割符号
pause_on_end: false # 运行完毕后是否暂停
strict_mode: false # 是否严格模式,如：int配置为空时，严格模式将会报错，非严格模式默认为0
sheets:
  - id: 9
    type: table  # 表格类型为数据表格
    file: "data/settings.xlsx"  # 表格文件路径
    sheet: "settings" # sheet表格名
    type_name: "SettingsCfg" # 导出类型名
    transpose: true  # 是否转置表格
  - id: 10
    type: table
    file: "data/item.xlsx"
    sheet: "item"
    type_name: "ItemCfg"
```

3. 修改二进制文件以及ts文件的导出路径
``` yaml
exports:
  - id: 1
    type: "proto_bytes"
    path: "../../mm-client/assets/packages/config/data"
  - id: 2
    type: "ts"
    path: "../../mm-client/assets/scripts/game/config/DataModel.ts"
    imports:
      - "import { DropGroup,ItemGroup,RewardData } from './ConfigExtension';"  # 在导出文件的目录中，可以添加自定义类型的引用代码
```

``` typescript
// ConfigExtension.ts 定义类型
export class DropGroup { //类型|ID|数量|权重|解锁等级
    type: number;
    id: number;
    num: number;
    weight: number;
    limit: number;
}
```

4. 修改项目中的tsconfig.json文件，添加: `"allowSyntheticDefaultImports":true`, 解决protobuf引用失败错误
``` json
{
  /* Base configuration. Do not edit this field. */
  "extends": "./temp/tsconfig.cocos.json",

  /* Add your custom configuration here. */
  "compilerOptions": {
    "strict": false,
    "allowSyntheticDefaultImports":true
  }
}
```

5. 下载源码中的reader/typescript/DataAccess.ts至项目中
6. 安装protobufjs: `npm i protobufjs -S`(可能需要重启编辑器)
7. 添加自定义脚本用于管理配置表：
``` typescript
export class ConfigManager {
    // 配置表实例，外部直接使用
    static ItemTable: DataTable<GameConfig.ItemCfg>;
    static DropTable: DataTable<GameConfig.DropCfg>;
    static SettingsTable: DataTable<GameConfig.SettingsCfg>;

    private static _modelConverter: {[key:string]: (data: DataModel, fieldName:string, value: string)=>void} = {};
    // 加载配置表
    static loadConfig() {
        DataConverter.convertHandler = this.modelConvert.bind(this);

        // 设置bin文件加载方式，返回二进制数据
        DataAccess.initial("data/", (datafile)=>{
            let asset = ResManager.get<BufferAsset>(ResConst.AB_CONFIG, datafile);
            return new Uint8Array(asset.buffer());
        });

        // 配置表实例化
        this.ItemTable = new DataTable(GameConfig.ItemCfg);
        this.DropTable = new DataTable(GameConfig.DropCfg);
        this.SettingsTable = new DataTable(GameConfig.SettingsCfg);

        // 设置自定义类型转换函数
        this._modelConverter["DropGroup"] = this.parseDropGroup.bind(this);
        this._modelConverter["DropGroup[]"] = this.parseDropGroups.bind(this);
    }    

    // 通过注册的自定类型转换相关类型
    private static modelConvert(data: DataModel, fieldName:string, value: string, alias?: string): any {        
        let convert = this._modelConverter[alias];
        if(!convert) {
            console.error(`can not find alias named ${alias} in field=${fieldName}`);
            return;
        }

        return convert(data, fieldName, value);
    }

    // 直接返回第一行数据为配置数据
    static get settings(): GameConfig.SettingsCfg {
        return this.SettingsTable.getItem(1);
    }

    private static parseDropGroup(data: DataModel, fieldName:string, value: string) {
        if(!value) {
            return null;
        }

        let temp = new DropGroup();
        let t = value.split(',');
        if(t.length != 5) {
            //@ts-ignore
            console.error(`DropGroup Config Error: ID=${data.ID}, Filed=${fieldName}`);
            return null;
        }

        temp.type = parseInt(t[0]);
        temp.id = parseInt(t[1]);
        temp.num = parseInt(t[2]);
        temp.weight = parseInt(t[3]);
        temp.limit = parseInt(t[4]);

        return temp;
    }

    private static parseDropGroups(data: DataModel, fieldName, value: string) {
        let d = value.split('|');
        let table: DropGroup[] = [];
        for (let l = 0; l < d.length; ++l) {
            let temp = this.parseDropGroup(data, fieldName, d[l]);
            if(!temp) {
                continue;
            }
            table.push(temp);
        }
        return table;
    }
}
``` 

7. 愉快的使用：
``` typescript
let itemData = ConfigManager.DropTable.getItem(this._item.id);
// 表格定义字段
console.log(itemData.ID);
console.log(itemData.Item1);
// (DropGroup[])获取自定义类型，获取时自动转换，并会缓存结果
console.log(itemData.getItem1());
```
