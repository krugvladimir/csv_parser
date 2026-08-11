// csv_parser.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.regex.*;

public class csv_parser {
    public static void main(String[] args) throws IOException {
        if (args.length < 1) {
            System.out.println("Usage: java csv_parser <input.csv> [output.json]");
            return;
        }

        String inputFile = args[0];
        String outputFile = args.length > 1 ? args[1] : "output.json";

        System.out.println("📊 CSV → JSON Парсер (Java)");

        String content = new String(Files.readAllBytes(Paths.get(inputFile)));
        String[] lines = content.split("\\n");
        if (lines.length == 0) {
            System.out.println("❌ CSV пуст.");
            return;
        }

        String delimiter = detectDelimiter(lines[0]);
        System.out.println("🔍 Определён разделитель: " + delimiter);

        String[] headers = lines[0].split(delimiter);
        for (int i = 0; i < headers.length; i++) headers[i] = headers[i].trim();
        System.out.println("📋 Заголовки: " + String.join(", ", headers));

        List<Map<String, Object>> data = new ArrayList<>();
        for (int i = 1; i < lines.length; i++) {
            String line = lines[i].trim();
            if (line.isEmpty()) continue;
            String[] values = line.split(delimiter);
            Map<String, Object> record = new LinkedHashMap<>();
            for (int j = 0; j < headers.length; j++) {
                String val = j < values.length ? values[j].trim() : "";
                record.put(headers[j], parseValue(val));
            }
            data.add(record);
        }

        String json = toJson(data);
        Files.write(Paths.get(outputFile), json.getBytes());

        System.out.println("✅ Экспортировано " + data.size() + " записей в " + outputFile);
    }

    private static String detectDelimiter(String line) {
        int commaCount = line.split(",").length - 1;
        int semicolonCount = line.split(";").length - 1;
        int tabCount = line.split("\t").length - 1;
        if (semicolonCount > commaCount && semicolonCount > tabCount) return ";";
        if (tabCount > commaCount && tabCount > semicolonCount) return "\t";
        return ",";
    }

    private static Object parseValue(String value) {
        if (value.isEmpty()) return null;
        if (value.equalsIgnoreCase("true")) return true;
        if (value.equalsIgnoreCase("false")) return false;
        if (value.equalsIgnoreCase("null") || value.equalsIgnoreCase("none")) return null;
        try {
            if (value.contains(".")) return Double.parseDouble(value);
            return Integer.parseInt(value);
        } catch (NumberFormatException e) {
            return value;
        }
    }

    private static String toJson(List<Map<String, Object>> data) {
        StringBuilder sb = new StringBuilder();
        sb.append("[\n");
        for (int i = 0; i < data.size(); i++) {
            Map<String, Object> record = data.get(i);
            sb.append("  {");
            int j = 0;
            for (Map.Entry<String, Object> entry : record.entrySet()) {
                sb.append("\"").append(entry.getKey()).append("\": ");
                Object val = entry.getValue();
                if (val instanceof String) {
                    sb.append("\"").append(val).append("\"");
                } else {
                    sb.append(val);
                }
                if (j < record.size() - 1) sb.append(", ");
                j++;
            }
            sb.append("}");
            if (i < data.size() - 1) sb.append(",");
            sb.append("\n");
        }
        sb.append("]");
        return sb.toString();
    }
}
