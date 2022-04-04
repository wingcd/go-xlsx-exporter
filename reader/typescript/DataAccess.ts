import { Asset, BufferAsset, resources } from "cc";
import * as $protobuf from "protobufjs";

export interface IMessage {
    create(props?:any);
    encode(message: any, writer?: $protobuf.Writer): $protobuf.Writer;
    encodeDelimited(message?: any, writer?: $protobuf.Writer): $protobuf.Writer;
    decode(reader: ($protobuf.Reader|Uint8Array), length?: number): any;
    decodeDelimited(reader: ($protobuf.Reader|Uint8Array)): any;
    verify(message: { [k: string]: any }): (string|null);
    fromObject(object: { [k: string]: any }): any;
    toObject(message: any, options?: $protobuf.IConversionOptions): { [k: string]: any };
    toJSON(): { [k: string]: any };
}

interface IDataArray {
    Items: any[];
}

type FileNameGenerateHandler = (typeName:string)=>string;
type LoadDataHandler = (datafile: string) => Uint8Array;
type KeyType = number|string;
type DataKeyMap = {[key:KeyType]: any};

export class DataAccess
{
    /// <summary>
    /// 是否使用ProtoMember的tag作为hashtable的key
    /// </summary>
    public static useProtoMemberTagAsHashtableKey = false;
    public static cacheHashValue = true;
    public static dataExt = ".bin";

    /// <summary>
    /// 数据文件以assets为根目录的路径
    /// </summary>
    public static dataDir: string;
    public static generator: FileNameGenerateHandler;
    public static loader: LoadDataHandler;

    private static _items:{[key:string]:DataItem} = {};
    private static _tables: {[key:string]:DataTable} = {};

    public static initial(dataDir:string, loadHandle?:LoadDataHandler, fileNameGenerateHandle?:FileNameGenerateHandler){
        this.dataDir = dataDir;
        this.generator = fileNameGenerateHandle;
        this.loader = loadHandle;

        $protobuf.Method
    }

    public static getDataItem(dataType: Function): DataItem {
        if(this._items[dataType.name]) {
            return this._items[dataType.name];
        }
        return this._items[dataType.name] = new DataItem(dataType);
    }

    public static getDataTable(dataType:Function, keyName = "ID") : DataTable {
        if(this._tables[dataType.name]) {
            return this._tables[dataType.name];
        }
        return this._tables[dataType.name] = new DataTable(dataType, keyName);
    }
}

export class DataItem {
    protected localGenerator: FileNameGenerateHandler;
    protected localLoader: LoadDataHandler;

    protected dataType: ObjectConstructor;
    protected source: BufferAsset;

    public constructor(dataType: Function) {
        this.dataType = dataType as ObjectConstructor;
    }

    private onGenerateFilename(typeName: string): string
    {
        if (this.localGenerator != null)
        {
            return this.localGenerator(typeName);
        }

        if (DataAccess.generator != null)
        {
            return DataAccess.generator(typeName);
        }

        typeName = typeName.replace("_ARRAY", "");
        return DataAccess.dataDir + typeName.toLocaleLowerCase();
    }

    protected onLoadData(typeName: string): Uint8Array
    {
        if (this.source != null)
        {
            return new Uint8Array(this.source.buffer());
        }

        var datafile = this.onGenerateFilename(typeName);

        if (this.localLoader != null)
        {
            return this.localLoader(datafile);
        }

        if (DataAccess.loader != null)
        {
            return DataAccess.loader(datafile);
        }

        this.source = resources.get<BufferAsset>(datafile, BufferAsset);
        if(this.source) {
            return new Uint8Array(this.source.buffer());
        }

        return null;
    }

    protected setSource(source: BufferAsset){
        this.source = source;
    }

    public initial(dataType:ObjectConstructor, loadHandle?:LoadDataHandler, fileNameGenerateHandle?:FileNameGenerateHandler) : void {
        this.dataType = dataType;
        this.localLoader = loadHandle;
        this.localGenerator = fileNameGenerateHandle;
    }

    public static create(dataType: ObjectConstructor, source: BufferAsset): DataItem {
        var instance = new DataItem(dataType);
        instance.setSource(source);
        return instance;
    }

    public clear() {
        if (this.source != null)
        {
            this.source.destroy();
            this.source = null;
        }
    }

    protected load(): any {
        var buffer = this.onLoadData(this.dataType["__type_name__"]);
        var msgType:IMessage = this.dataType as any;        
        return msgType.decode(buffer);
    }

    private _item: any[];
    public get data(): any[] {
        if (this._item == null)
        {
            this._item = this.load();
        }
        return this._item;
    }

    public set data(value: any) {
        this._item = value;
    }
}

export class DataTable extends DataItem {
    private _keyName = "ID";
    public get keyName() {
        return this._keyName;
    }

    constructor(dataType:Function, keyName = "ID") {
        super(dataType);

        this._keyName = keyName;
    }

    protected load() : any[] {
        // var arrTypeName = this.dataType["__type_name__"] + "_ARRAY"; 

        var buffer = this.onLoadData(this.dataType["__type_name__"]);
        var msgType:IMessage = this.dataType as any;
        var message = msgType.decode(buffer);
        return (message as IDataArray).Items;        
    }

    public static create(dataType: ObjectConstructor, source: BufferAsset, keyName = "ID"): DataTable{
        var instance = new DataTable(dataType, keyName);
        instance.setSource(source);
        return instance;
    }

    private _itemMap: DataKeyMap;
    public get itemMap(): DataKeyMap {
        if (this._itemMap == null)
        {
            this._itemMap = this.initDataAsDict();
        }
        return this._itemMap;
    }

    private _items: any[];
    public get items(): any[] {
        try{
            if (this._items == null){
                this._items = this.initDataAsList();
            }
        }catch(e){
            console.log(`config data load error: ${e}`)
        }
        return this._items;
    }

    private _ids: KeyType[];
    public get IDs(): KeyType[]
    {
        if (this._ids == null)
        {
            this._ids = this.items.map((val, idx, arr)=>{
                arr.push(val[this.keyName]);
            }, []) as any as KeyType[];
        }
        return this._ids;
    }

    public getItem(key: KeyType): any
    {
        return this.itemMap[key];
    }

    protected initDataAsDict(): DataKeyMap {
        let itemMap: DataKeyMap = {};
        try{
            this.items.forEach((val)=>{
                itemMap[val[this.keyName]] = val;
            });
            return itemMap;
        }
        catch(e) {
            console.error(`can not get data map by key ${this.keyName}`);
            return itemMap;
        }
    }

    protected initDataAsList(): any[]{
        return this.load();
    }

    public contains(id: KeyType): boolean{
        return this.itemMap[id] != null;
    }

    public clear(){
        super.clear();

        this._itemMap = null;
        this._items = null;
    }
}
