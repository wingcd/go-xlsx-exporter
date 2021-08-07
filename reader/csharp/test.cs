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
        var classData = DataContainer<int, PClass>.Instance.Items;
        foreach (var item in classData)
        {
            Console.WriteLine($"{item.ID}, {item.Name}, {item.Level}, {item.Type}");
        }

        Console.WriteLine();
        Console.WriteLine("user table:");
        var userdata = DataContainer<int, User>.Instance.Items;
        foreach (var item in userdata)
        {
            Console.WriteLine($"{item.ID}, {item.Name}, {item.Age}, {item.Head}");
        }

        var lanKey = "1";
        I18N.SetLanguage("cn");
        Console.WriteLine($"中文：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");
        I18N.SetLanguage("en");
        Console.WriteLine($"english：tanslate key={lanKey}, text={I18N.Translate(lanKey)}");

        ConsoleKey key;
        do { key = Console.ReadKey().Key; }
        while (key != ConsoleKey.Q);
    }
}