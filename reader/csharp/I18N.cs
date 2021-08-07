using ProtoBuf;
using System;
using System.Collections.Generic;
using System.IO;

[Serializable]
[ProtoContract]
public class Language : PBDataModel
{
    [ProtoMember(1)]
    public string ID { get; set; }

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

    public static string Translate(string key)
    {
        Language lan;
        if (Instance.ItemMap.TryGetValue(key, out lan))
        {
            return lan.Text;
        }

        return "";
    }
}
