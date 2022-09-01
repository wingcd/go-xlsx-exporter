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

export interface IDataArray {
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

    private static _items:{[key:string]:any} = {};
    private static _tables: {[key:string]:any} = {};

    public static initial(dataDir:string, loadHandle?:LoadDataHandler, fileNameGenerateHandle?:FileNameGenerateHandler){
        this.dataDir = dataDir;
        this.generator = fileNameGenerateHandle;
        this.loader = loadHandle;
    }

    /**
     * 获取配置表
     * @param dataType 配置的数据类型
     * @returns 
     */
    public static getDataItem<T>(dataType: new()=>T): DataItem<T> {
        let typename = dataType["__type_name__"];
        if(this._items[typename]) {
            return this._items[typename];
        }
        return this._items[typename] = new DataItem<T>(dataType);
    }

    /**
     * 获取配置表, 可自定义主键名称
     * @param dataType 配置的数据类型
     * @param keyName 主键名称
     * @returns 
     */
    public static getDataTable<T>(dataType:new()=>T, keyName = "ID") : DataTable<T> {
        let typename = dataType["__type_name__"];
        if(this._tables[typename]) {
            return this._tables[typename];
        }
        return this._tables[typename] = new DataTable<T>(dataType, keyName);
    }
}

export class DataItem<T> {
    protected localGenerator: FileNameGenerateHandler;
    protected localLoader: LoadDataHandler;

    protected dataType: new()=>T;
    protected source: BufferAsset;

    public constructor(dataType: new()=>T) {
        this.dataType = dataType;
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

    public initial(dataType:new()=>T, loadHandle?:LoadDataHandler, fileNameGenerateHandle?:FileNameGenerateHandler) : void {
        this.dataType = dataType;
        this.localLoader = loadHandle;
        this.localGenerator = fileNameGenerateHandle;
    }

    public static create<T>(dataType: new()=>T, source: BufferAsset): DataItem<T> {
        var instance = new DataItem<T>(dataType);
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
        var buffer = this.onLoadData(this.dataType["__array_type_name__"]);
        var msgType:IMessage = this.dataType["__array_type__"] as any;   
        return msgType.decode(buffer);
    }

    private _item: T;
    public get data(): Readonly<T> {
        if (this._item == null)
        {
            this._item = this.load();
        }
        return this._item;
    }
}

export class DataTable<T> extends DataItem<T> {
    private _keyName = "ID";
    public get keyName() {
        return this._keyName;
    }

    constructor(dataType:new()=>T, keyName = "ID") {
        super(dataType);

        this._keyName = keyName;
    }

    protected load() : T[] {
        var arrTypeName = this.dataType["__array_type_name__"];

        var buffer = this.onLoadData(arrTypeName);
        var msgType:IMessage = this.dataType["getArrayType"]() as any;
        var message = msgType.decode(buffer);
        return (message as IDataArray).Items;        
    }

    public static create<T>(dataType: new()=>T, source: BufferAsset, keyName = "ID"): DataTable<T>{
        var instance = new DataTable<T>(dataType, keyName);
        instance.setSource(source);
        return instance;
    }

    private _itemMap: DataKeyMap;
    public get itemMap(): Readonly<DataKeyMap> {
        if (this._itemMap == null)
        {
            this._itemMap = this.initDataAsDict();
        }
        return this._itemMap;
    }

    private _items: T[];
    public get items(): Readonly<Readonly<T>[]> {
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
    public get IDs(): Readonly<KeyType[]>
    {
        if (this._ids == null)
        {
            this._ids = this.items.map((val, idx, arr)=>{
                //@ts-ignore
                arr.push(val[this.keyName]);
            }, []) as any as KeyType[];
        }
        return this._ids;
    }

    public getItem(key: KeyType): Readonly<T>
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

    protected initDataAsList(): T[]{
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
