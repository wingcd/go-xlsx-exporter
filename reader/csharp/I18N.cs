using ProtoBuf;
using System;
using System.Collections.Generic;
using System.IO;

[Serializable]
[ProtoContract]
public class Language : PBDataModel
{
    [ProtoMember(1)]
    public string Id { get; set; }

    [ProtoMember(2)]
    public string Text { get; set; }
}

[Serializable]
[ProtoContract]
public class Language_ARRAY : PBDataModels
{
    [ProtoMember(1)]
    public List<Language> Items { get; set; }
}

public class I18N : DataContainer<string, Language> 
{
    public static string CurrentLanguage { get; private set; }
    public static void SetLanguage(string lan = "cn")
    {
        if (CurrentLanguage != "")
        {
            Instance.Clear();
        }

        CurrentLanguage = lan;

        Instance.Initial(null, GetFilename);
    }

    protected static string GetFilename(string typeName)
    {
        var datafile = Path.Combine(DataAccess.DataDir, typeName.ToLower());
        return $"{datafile}.{CurrentLanguage}";
    }
    
    public static string Translate(int key, params object[] args)
    {
        return Translate(key.ToString(), args);
    }

    public static string Translate(string key, params object[] args)
    {
        Language lan;
        if (Instance.ItemMap.TryGetValue(key, out lan))
        {
            if (args == null || args.Length == 0)
            {
                return lan.Text;
            }
            else
            {
                return string.Format(lan.Text, args);
            }
        }

        return null;
    }
}
