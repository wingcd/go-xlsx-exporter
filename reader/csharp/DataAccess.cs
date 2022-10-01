using ProtoBuf;
using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using UnityEngine;

#if LUA_SUPPORT
using LuaInterface;
#endif

internal class T_LOG
{
    public static void Log(string info)
    {
        Debug.Log(info);
    }
    
    public static void LogError(string info)
    {
        Debug.LogError(info);
    }
}

public class PBDataModel
{
    private Dictionary<string,object> _converted = new Dictionary<string, object>();

    protected object GetConvertData(string fieldName, object value, string alias, bool cachable)
    {
        if(cachable && _converted.ContainsKey(fieldName))
        {
            return _converted[fieldName];
        }

        if(DataAccess.DataConvertHandler == null)
        {
            throw new Exception($"convert field {fieldName} value need a convetor");
        }

        var data = DataAccess.DataConvertHandler(this, fieldName, value, alias, cachable);
        if (cachable)
        {
            _converted[fieldName] = data;
        }
        return data;
    }     
    
    public object Clone()
    {
        var type = GetType();
        var clone = Activator.CreateInstance(type);
        var props = type.GetProperties();
        for(var i=0;i<props.Length;i++)
        {
            var prop = props[i];
            var attris = prop.GetCustomAttributes(false);
            for(var ai = 0; ai < attris.Length; ai++)
            {
                var attri = attris[ai];
                if(attri is ProtoMemberAttribute)
                {
                    prop.SetValue(clone, prop.GetValue(this));

                    break;
                }
            }
        }
        return clone;
    }

    #if LUA_SUPPORT
    public virtual LuaTable ToLuaTable()
    {
        if (DataAccess.luaState == null)
        {
            Debug.LogError("must set luastate first!");
        }
        
        var table = DataAccess.luaState.NewTable();
        var type = GetType();
        var props = type.GetProperties();
        for(var i=0;i<props.Length;i++)
        {
            var prop = props[i];
            var attris = prop.GetCustomAttributes(false);
            for(var ai = 0; ai < attris.Length; ai++)
            {
                var attri = attris[ai];
                if(attri is ProtoMemberAttribute)
                {
                    if(DataAccess.UseProtoMemberTagAsHashtableKey)
                    {
                        table[(attri as ProtoMemberAttribute).Tag] = prop.GetValue(this);
                    }
                    else
                    {
                        table[prop.Name] = prop.GetValue(this);
                    }

                    break;
                }
            }
        }
        return table;
    }
    #endif
    
    public object GetPropValue(string name)
    {
        var tp = GetType();
        return tp.GetProperty(name).GetValue(this);
    }
}

public class PBDataModels
{
    
}

public delegate string FileNameGenerateHandler(string typeName);
public delegate byte[] LoadDataHandler(string datafile);

public partial class DataAccess
{
    #if LUA_SUPPORT
    public static LuaState luaState;
    #endif
    
    /// <summary>
    /// 是否使用ProtoMember的tag作为hashtable的key
    /// </summary>
    public static bool UseProtoMemberTagAsHashtableKey = false;
    public static bool CacheHashValue = true;
    public static string DataExt = ".bytes";
    public const string DefaultIDName = "Id";

    /// <summary>
    /// 数据文件以assets为根目录的路径
    /// </summary>
    public static string DataDir { get; private set; }
    public static LoadDataHandler Loader { get; private set; }
    public static FileNameGenerateHandler Generator { get; private set; }

    public delegate object ConvertHandler(object item, string field, object data, string alias, bool cachable);
    public static ConvertHandler DataConvertHandler;

    public static void Initial(string dataDir, LoadDataHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        DataDir = dataDir;
        Loader = loadHandle;
        Generator = fileNameGenerateHandle;
    }
}

public partial class DataContainer<TItem> : DataContainer
    where TItem : PBDataModel
{
    protected LoadDataHandler localLoader;
    protected FileNameGenerateHandler localGenerator;

    private string OnGenerateFilename(string typeName)
    {
        if (DataAccess.DataDir == null)
        {
            Debug.LogError("you need initial data access first!");
            return typeName.ToLower();
        }
        
        if (localGenerator != null)
        {
            return localGenerator(typeName);
        }

        if (DataAccess.Generator != null)
        {
            return DataAccess.Generator(typeName);
        }

        return Path.Combine(DataAccess.DataDir, typeName.ToLower());
    }

    protected byte[] OnLoadData(string typeName)
    {
        var datafile = OnGenerateFilename(typeName);

        if (localLoader != null)
        {
            return localLoader(datafile);
        }

        if (DataAccess.Loader != null)
        {
            return DataAccess.Loader(datafile);
        }

#if UNITY_ENGINE
        TextAsset text = Resources.Load<TextAsset>(dataPath);
		var bytes = text.bytes;
        GameObject.Destroy(text);
        return bytes;
#else
        var filename = datafile + DataAccess.DataExt;
        if (File.Exists(filename))
        {
            return File.ReadAllBytes(filename);
        }
        else
        {
            throw new Exception($"can not load data:{datafile}");
        }
#endif
    }

    public void Initial(LoadDataHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        localLoader = loadHandle;
        localGenerator = fileNameGenerateHandle;
    }

    static DataContainer<TItem> _instance;
    public static DataContainer<TItem> Instance
    {
        get
        {
            if (_instance == null)
            {
                _instance = new DataContainer<TItem>();
            }
            return _instance;
        }
    }

    public override void Clear()
    {
        _item = null;
    }

    public DataContainer<TItem> Preload(byte[] bytes = null)
    {
        if (_item == null)
        {
            _item = Load(bytes);
        }

        return this;
    }

    private TItem Load(byte[] bytes = null)
    {
        var type = typeof(TItem);
        bytes = bytes ?? OnLoadData(type.Name);

        using (var stream = new MemoryStream(bytes))
        {
            object data = Serializer.NonGeneric.Deserialize(type, stream);
            return data as TItem;
        }
    }

    TItem _item;
    public TItem Data
    {
        get
        {
            if (_item == null)
            {
                try
                {
                    _item = Load();
                }
                catch (Exception e)
                {
                    T_LOG.LogError(typeof(TItem).Name + e);
                    return null;
                }
            }
            return _item;
        }
        set
        {
            _item = value;
        }
    }
}

#if LUA_SUPPORT
public interface ILuaDataContainer<TID> where TID : IComparable
{
    int Count { get; }
    List<TID> IDs { get; }
    LuaTable GetTableByIndex(int index);
    LuaTable GetTable(TID ID);
}
#endif

public class DataContainer
{
    public virtual void Clear()
    {
        
    }
}

public partial class DataContainer<TID, TItem> : DataContainer<TItem>
#if LUA_SUPPORT
    , ILuaDataContainer<TID>
#endif
    where TID : IComparable
    where TItem : PBDataModel
{

    private List<TItem> Load(byte[] bytes = null)
    {
        var type = typeof(TItem);
        var arrTypeName = type.FullName + "_ARRAY";
        var arrayType = type.Assembly.GetType(arrTypeName);
        if (arrayType == null)
        {
            throw new Exception($"can not find data type:{arrTypeName}");
        }

        bytes = bytes ?? OnLoadData(type.Name);

        using (var stream = new MemoryStream(bytes))
        {
            object array = Serializer.NonGeneric.Deserialize(arrayType, stream);
            var pItems = arrayType.GetProperty("Items");
            return pItems.GetValue(array) as List<TItem>;
        }
    }

    static DataContainer<TID, TItem> _instance;
    public new static DataContainer<TID, TItem> Instance
    {
        get
        {
            if (_instance == null)
            {
                _instance = new DataContainer<TID, TItem>();
            }
            return _instance;
        }
    }

    public static DataContainer<TID, TItem> GetInstance(string keyName)
    {
        if (_instance == null)
        {
            _instance = new DataContainer<TID, TItem>(keyName);
        }
        return _instance;
    }

    string keyName;
    Dictionary<TID, TItem> _itemMap;
    public Dictionary<TID, TItem> ItemMap
    {
        get
        {
            if (_itemMap == null)
            {
                _itemMap = InitDataAsDict();
            }
            return _itemMap;
        }
        private set
        {
            _itemMap = value;
        }
    }
    
    public new DataContainer<TID, TItem> Preload(byte[] bytes = null)
    {
        if (_items == null)
        {
            _items = InitDataAsList(bytes);
        }

        return this;
    }

    List<TItem> _items;
    public List<TItem> Items
    {
        get
        {
            if (_items == null)
            {
                _items = InitDataAsList();
            }
            return _items;
        }
        private set
        {
            _items = value;
        }
    }

    public int Count
    {
        get
        {
            if (Items == null)
            {
                return 0;
            }

            return Items.Count;
        }
    }

    List<TID> _IDs;
    public List<TID> IDs
    {
        get
        {
            if (_IDs == null)
            {
                _IDs = new List<TID>(ItemMap.Keys);
            }
            return _IDs;
        }
        private set
        {
            _IDs = value;
        }
    }

    public TItem this[TID ID]
    {
        get
        {
            if (!ItemMap.ContainsKey(ID))
            {
                T_LOG.Log($"can not find id={ID} in {typeof(TItem).Name}");
                return default(TItem);
            }
            return ItemMap[ID];
        }
    }

    #if LUA_SUPPORT
    Dictionary<TID, LuaTable> luaTables = null;

    public LuaTable GetTableByIndex(int index)
    {
        if (index >= this.IDs.Count)
        {
            return null;
        }

        var id = this.IDs[index];
        return GetTable(id);
    }

    public LuaTable GetTable(TID ID)
    {
        if(DataAccess.CacheHashValue && luaTables == null)
        {
            luaTables = new Dictionary<TID, LuaTable>();
        }

        if (luaTables != null && luaTables.ContainsKey(ID))
        {
            return luaTables[ID];
        }
        else if(Contains(ID))
        {
            var ht = this[ID].ToLuaTable();
            if (luaTables != null)
            {
                luaTables[ID] = ht;
            }
            return ht;
        }
        return null;
    }
    #endif

    public DataContainer(string keyName = DataAccess.DefaultIDName)
    {
        this.keyName = keyName;
    }

    protected Dictionary<TID, TItem> InitDataAsDict()
    {
        Dictionary<TID, TItem> itemMap = null;
        try
        {
            var propID = typeof(TItem).GetProperty(keyName);
            if (propID == null)
            {
                throw new Exception($"can not find key:{keyName}");
            }
            itemMap = new Dictionary<TID, TItem>();
            var itr = Items.GetEnumerator();
            while (itr.MoveNext())
            {
                TID ID = (TID)propID.GetValue(itr.Current);
                if (ID == null || ID is string && ID as string == "")
                {
                    Debug.LogWarning("id can not be null or empty value!");
                    continue;
                }
                itemMap[ID] = itr.Current;
            }
            return itemMap;
        }
        catch (Exception e)
        {
            T_LOG.LogError(typeof(TItem).Name + e);
            return itemMap;
        }
    }

    protected List<TItem> InitDataAsList(byte[] bytes = null)
    {
        List<TItem> list = null;
        try
        {
            return Load(bytes);
        }
        catch (Exception e)
        {
            T_LOG.LogError(typeof(TItem).Name + e);
            return list;
        }
    }

    public bool Contains(TID ID)
    {
        if (ItemMap != null)
        {
            return ItemMap.ContainsKey(ID);
        }
        return false;
    }

    public override void Clear()
    {
        base.Clear();

        if (_itemMap != null)
        {
            _itemMap.Clear();
            _itemMap = null;
        }

        if (_items != null)
        {
            _items.Clear();
            _items = null;
        }
        
        #if LUA_SUPPORT
        if (luaTables != null)
        {
            foreach (var table in luaTables)
            {
                table.Value.Dispose();
            }

            luaTables = null;
        }
        #endif
    }
}
