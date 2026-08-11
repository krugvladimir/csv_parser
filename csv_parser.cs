// csv_parser.cs — C# версия

using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;

class Program
{
    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine("Usage: dotnet run <input.csv> [output.json]");
            return;
        }

        string inputFile = args[0];
        string outputFile = args.Length > 1 ? args[1] : "output.json";

        Console.WriteLine("📊 CSV → JSON Парсер (C#)");

        if (!File.Exists(inputFile))
        {
            Console.WriteLine($"❌ Файл {inputFile} не найден.");
            return;
        }

        string content = File.ReadAllText(inputFile);
        var lines = content.Split('\n').Where(l => !string.IsNullOrWhiteSpace(l)).ToArray();
        if (lines.Length == 0)
        {
            Console.WriteLine("❌ CSV пуст.");
            return;
        }

        string delimiter = DetectDelimiter(lines[0]);
        Console.WriteLine($"🔍 Определён разделитель: {delimiter}");

        var headers = lines[0].Split(delimiter).Select(h => h.Trim()).ToArray();
        Console.WriteLine($"📋 Заголовки: {string.Join(", ", headers)}");

        var data = new List<Dictionary<string, object>>();
        for (int i = 1; i < lines.Length; i++)
        {
            var values = lines[i].Split(delimiter).Select(v => v.Trim()).ToArray();
            var record = new Dictionary<string, object>();
            for (int j = 0; j < headers.Length; j++)
            {
                string val = j < values.Length ? values[j] : "";
                record[headers[j]] = ParseValue(val);
            }
            data.Add(record);
        }

        string json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(outputFile, json);

        Console.WriteLine($"✅ Экспортировано {data.Count} записей в {outputFile}");
    }

    static string DetectDelimiter(string line)
    {
        int comma = line.Count(c => c == ',');
        int semicolon = line.Count(c => c == ';');
        int tab = line.Count(c => c == '\t');
        if (semicolon > comma && semicolon > tab) return ";";
        if (tab > comma && tab > semicolon) return "\t";
        return ",";
    }

    static object ParseValue(string value)
    {
        if (string.IsNullOrEmpty(value)) return null;
        string v = value.Trim().ToLower();
        if (v == "true") return true;
        if (v == "false") return false;
        if (v == "null" || v == "none") return null;
        if (int.TryParse(value, out int i)) return i;
        if (double.TryParse(value, out double d)) return d;
        return value;
    }
}
