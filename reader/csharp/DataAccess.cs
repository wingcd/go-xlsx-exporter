using ProtoBuf;
using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
#if UNITY_ENGINE
using UnityEngine;
#endif

internal class T_LOG
{
    public static void Log(string info)
    {
#if UNITY_ENGINE
        Debug.Log(info);
#else
        Console.WriteLine(info);
#endif
    }
}

public class PBDataModel
{
    public virtual Hashtable ToHashtable()
    {
        var hashtable = new Hashtable();
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
                        hashtable[(attri as ProtoMemberAttribute).Tag] = prop.GetValue(this);
                    }
                    else
                    {
                        hashtable[prop.Name] = prop.GetValue(this);
                    }

                    break;
                }
            }
        }
        return hashtable;
    }
}

public class PBDataModels : PBDataModel
{
    public override Hashtable ToHashtable()
    {
        return null;
    }
}

public delegate string FileNameGenerateHandler(string typeName);
public delegate byte[] LoadDataHandler(string datafile);

public class DataAccess
{
    /// <summary>
    /// 是否使用ProtoMember的tag作为hashtable的key
    /// </summary>
    public static bool UseProtoMemberTagAsHashtableKey = false;
    public static bool CacheHashValue = true;
    public static string DataExt = ".bytes";

    /// <summary>
    /// 数据文件以assets为根目录的路径
    /// </summary>
    public static string DataDir { get; private set; }
    public static LoadDataHandler Loader { get; private set; }
    public static FileNameGenerateHandler Generator { get; private set; }

    public static void Initial(string dataDir, LoadDataHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        DataDir = dataDir;
        Loader = loadHandle;
        Generator = fileNameGenerateHandle;
    }
}

public class DataContainer<TItem>
    where TItem : PBDataModel
{
    protected LoadDataHandler localLoader;
    protected FileNameGenerateHandler localGenerator;

    protected MemoryStream source;

    private string OnGenerateFilename(string typeName)
    {
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
        if (source != null)
        {
            return source.ToArray();
        }

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


    protected void SetSource(MemoryStream source)
    {
        this.source = source;
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

    /// <summary>
    /// </summary>
    /// <param name="source"></param>
    /// <returns></returns>
    public static DataContainer<TItem> CreateInstace(MemoryStream source)
    {
        var instance = new DataContainer<TItem>();
        instance.SetSource(source);
        return instance;
    }

    public virtual void Clear()
    {
        if (source != null)
        {
            source.Close();
            source = null;
        }
    }

    private TItem Load()
    {
        var type = typeof(TItem);
        var bytes = OnLoadData(type.Name);

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
                    T_LOG.Log(typeof(TItem).Name + e);
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

public class DataContainer<TID, TItem> : DataContainer<TItem>
    where TID : IComparable
    where TItem : PBDataModel
{

    private List<TItem> Load()
    {
        var type = typeof(TItem);
        var arrTypeName = type.FullName + "_ARRAY";
        var arrayType = type.Assembly.GetType(arrTypeName);
        if (arrayType == null)
        {
            throw new Exception($"can not find data type:{arrTypeName}");
        }

        var bytes = OnLoadData(type.Name);

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

    /// <summary>
    /// </summary>
    /// <param name="source"></param>
    /// <param name="keyName"></param>
    /// <returns></returns>
    public static DataContainer<TID, TItem> CreateInstace(MemoryStream source, string keyName = "ID")
    {
        var instance = new DataContainer<TID, TItem>(keyName);
        instance.SetSource(source);
        return instance;
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
                return default(TItem);
            }
            return ItemMap[ID];
        }
    }

    Dictionary<TID, Hashtable> hashtables = null;
    public Hashtable GetHashtable(TID ID)
    {
        if(DataAccess.CacheHashValue && hashtables == null)
        {
            hashtables = new Dictionary<TID, Hashtable>();
        }

        if (hashtables != null && hashtables.ContainsKey(ID))
        {
            return hashtables[ID];
        }
        else if(Contains(ID))
        {
            var ht = this[ID].ToHashtable();
            if (hashtables != null)
            {
                hashtables[ID] = ht;
            }
            return ht;
        }
        return null;
    }

    public DataContainer(string keyName = "ID")
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
                itemMap[ID] = itr.Current;
            }
            return itemMap;
        }
        catch (Exception e)
        {
            T_LOG.Log(typeof(TItem).Name + e);
            return itemMap;
        }
    }

    protected List<TItem> InitDataAsList()
    {
        List<TItem> list = null;
        try
        {
            return Load();
        }
        catch (Exception e)
        {
            T_LOG.Log(typeof(TItem).Name + e);
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
    }
}
