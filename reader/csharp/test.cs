using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

using ProtoBuf;
using PBGen;
using System.IO;
using System.Threading;

class Program
{
    static void Main(string[] args)
    {
        DataAccess.Initial("./data/");
        
        Console.WriteLine("class table:");
        var classData = DataContainer<uint, PClass>.Instance.Items;
        foreach (var item in classData)
        {
            Console.WriteLine($"{item.ID}, {item.Name}, {item.Level}, {item.Type}");
        }

        Console.WriteLine();
        Console.WriteLine("user table:");
        var userdata = DataContainer<uint, User>.Instance.Items;
        foreach (var item in userdata)
        {
            Console.WriteLine($"{item.ID}, {item.Name}, {item.Age}, {item.Head}");
        }

        var lanKey = "1";
        I18N.SetLanguage("cn");
        Console.WriteLine($"\n中文：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");
        I18N.SetLanguage("en");
        Console.WriteLine($"english：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");

        var settings = DataContainer<Settings>.Instance.Data;
        Console.WriteLine($"\n配置：maxconn={settings.MAX_CONNECT}, version={settings.VERSION}");

        Console.WriteLine($"\n转哈希表：");
        var userHT = DataContainer<uint, User>.Instance.GetHashtable(1);
        Console.WriteLine($"用户哈希值：{userHT["ID"]}, {userHT["Name"]}, {userHT["Age"]}, {userHT["Head"]}");

        ConsoleKey key;
        do { key = Console.ReadKey().Key; }
        while (key != ConsoleKey.Q);
    }
}