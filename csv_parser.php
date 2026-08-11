<?php
// csv_parser.php — PHP версия

function detectDelimiter($firstLine) {
    $comma = substr_count($firstLine, ',');
    $semicolon = substr_count($firstLine, ';');
    $tab = substr_count($firstLine, "\t");
    if ($semicolon > $comma && $semicolon > $tab) return ';';
    if ($tab > $comma && $tab > $semicolon) return "\t";
    return ',';
}

function parseValue($value) {
    $value = trim($value);
    if ($value === '') return null;
    $lower = strtolower($value);
    if ($lower === 'true') return true;
    if ($lower === 'false') return false;
    if ($lower === 'null' || $lower === 'none') return null;
    if (is_numeric($value)) {
        if (strpos($value, '.') !== false) return (float)$value;
        return (int)$value;
    }
    return $value;
}

function parseCSV($inputFile, $outputFile = null) {
    if (!file_exists($inputFile)) {
        echo "❌ Файл $inputFile не найден.\n";
        return false;
    }

    $content = file_get_contents($inputFile);
    $lines = array_filter(explode("\n", $content), function($l) {
        return trim($l) !== '';
    });
    $lines = array_values($lines);

    if (empty($lines)) {
        echo "❌ CSV пуст.\n";
        return false;
    }

    $delimiter = detectDelimiter($lines[0]);
    echo "🔍 Определён разделитель: " . ($delimiter === "\t" ? '\\t' : $delimiter) . "\n";

    $headers = array_map('trim', explode($delimiter, $lines[0]));
    echo "📋 Заголовки: " . implode(', ', $headers) . "\n";

    $data = [];
    for ($i = 1; $i < count($lines); $i++) {
        $values = array_map('trim', explode($delimiter, $lines[$i]));
        $record = [];
        foreach ($headers as $j => $header) {
            $val = isset($values[$j]) ? $values[$j] : '';
            $record[$header] = parseValue($val);
        }
        $data[] = $record;
    }

    if (!$outputFile) {
        $outputFile = pathinfo($inputFile, PATHINFO_FILENAME) . '.json';
    }

    file_put_contents($outputFile, json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
    echo "✅ Экспортировано " . count($data) . " записей в $outputFile\n";
    return true;
}

if ($argc < 2) {
    echo "Usage: php csv_parser.php <input.csv> [output.json]\n";
    exit(1);
}

$input = $argv[1];
$output = $argv[2] ?? null;

echo "📊 CSV → JSON Парсер (PHP)\n";
parseCSV($input, $output);
?>
