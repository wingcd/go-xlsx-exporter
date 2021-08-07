using ProtoBuf;
using System;
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

}

public class PBDataModels : PBDataModel
{

}

public delegate string FileNameGenerateHandler(string typeName);
public delegate byte[] LoadDataHandler(string datafile);

public class DataAccess
{
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

public class DataContainer<TID, TItem>
    where TID : IComparable
    where TItem : PBDataModel
{

    protected LoadDataHandler localLoader;
    protected FileNameGenerateHandler localGenerator;

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

    private byte[] OnLoadData(string typeName)
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

    public void Initial(LoadDataHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        localLoader = loadHandle;
        localGenerator = fileNameGenerateHandle;
    }

    public List<TItem> Load<TItem>()
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
    public static DataContainer<TID, TItem> Instance
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
    ///eg: var errDefines = DataContainer<uint, VCity.Deploy.ErrorDefine>.CreateInstace(new MemoryStream(txt.bytes), "tid").Items;
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

    MemoryStream source;
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
        set
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
        set
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


    public DataContainer(string keyName = "ID")
    {
        this.keyName = keyName;
    }

    void SetSource(MemoryStream source)
    {
        this.source = source;
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
            return Load<TItem>();
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

    public void Clear()
    {
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

        if (source != null)
        {
            source.Close();
            source = null;
        }
    }
}
