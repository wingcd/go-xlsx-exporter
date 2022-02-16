#if ASYNC_SUPPORT
using ProtoBuf;
using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using UnityEngine;
using Cysharp.Threading.Tasks;

public delegate UniTask<byte[]> LoadDataAsyncHandler(string datafile);

public partial class DataAccess
{
    public static LoadDataAsyncHandler LoaderAsync { get; private set; }

    public static void InitialAsync(string dataDir, LoadDataAsyncHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        DataDir = dataDir;
        LoaderAsync = loadHandle;
        Generator = fileNameGenerateHandle;
    }
}

public partial class DataContainer<TItem>
{
    protected LoadDataAsyncHandler localLoaderAsync;

    protected async UniTask<byte[]> OnLoadDataAsync(string typeName)
    {
        var datafile = OnGenerateFilename(typeName);

        if (localLoader != null)
        {
            return await localLoaderAsync(datafile);
        }

        if (DataAccess.LoaderAsync != null)
        {
            return await DataAccess.LoaderAsync(datafile);
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

    public void InitialAsync(LoadDataAsyncHandler loadHandle = null, FileNameGenerateHandler fileNameGenerateHandle = null)
    {
        localLoaderAsync = loadHandle;
        localGenerator = fileNameGenerateHandle;
    }
    
    public async UniTask<DataContainer<TItem>> PreloadAsync()
    {
        if (_item == null)
        {
            _item = Load();
        }

        return this;
    }

    private async UniTask<TItem> LoadAsync()
    {
        var type = typeof(TItem);
        var bytes = await OnLoadDataAsync(type.Name);

        using (var stream = new MemoryStream(bytes))
        {
            object data = Serializer.NonGeneric.Deserialize(type, stream);
            return data as TItem;
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

public partial class DataContainer<TID, TItem> : DataContainer<TItem>
#if LUA_SUPPORT
    , ILuaDataContainer<TID>
#endif
    where TID : IComparable
    where TItem : PBDataModel
{
    private async UniTask<List<TItem>> LoadAsync()
    {
        var type = typeof(TItem);
        var arrTypeName = type.FullName + "_ARRAY";
        var arrayType = type.Assembly.GetType(arrTypeName);
        if (arrayType == null)
        {
            throw new Exception($"can not find data type:{arrTypeName}");
        }

        var bytes = await OnLoadDataAsync(type.Name);

        using (var stream = new MemoryStream(bytes))
        {
            object array = Serializer.NonGeneric.Deserialize(arrayType, stream);
            var pItems = arrayType.GetProperty("Items");
            return pItems.GetValue(array) as List<TItem>;
        }
    }
    
    public new async UniTask<DataContainer<TID, TItem>> PreloadAsync()
    {
        if (_items == null)
        {
            _items = await InitDataAsListAsync();
        }

        return this;
    }

    protected async UniTask<Dictionary<TID, TItem>> InitDataAsDictAsync()
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

    protected async UniTask<List<TItem>> InitDataAsListAsync()
    {
        List<TItem> list = null;
        try
        {
            return await LoadAsync();
        }
        catch (Exception e)
        {
            T_LOG.LogError(typeof(TItem).Name + e);
            return list;
        }
    }
}

#endif
