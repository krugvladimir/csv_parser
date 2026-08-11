

### 1. `csv_parser.py` (Python)

```python
# csv_parser.py — Python версия

import csv
import json
import sys
import os
import re
from datetime import datetime

def detect_delimiter(file_path):
    """Автоопределение разделителя"""
    with open(file_path, 'r', encoding='utf-8-sig') as f:
        first_line = f.readline()
        if '\t' in first_line:
            return '\t'
        elif ';' in first_line:
            return ';'
        else:
            return ','

def parse_value(value):
    """Умное преобразование типов"""
    value = value.strip()
    if value == '':
        return None
    if value.lower() == 'true':
        return True
    if value.lower() == 'false':
        return False
    if value.lower() == 'null' or value.lower() == 'none':
        return None
    # Число
    try:
        if '.' in value:
            return float(value)
        else:
            return int(value)
    except ValueError:
        pass
    # Дата (попытка)
    try:
        for fmt in ['%Y-%m-%d', '%d-%m-%Y', '%Y/%m/%d']:
            try:
                return datetime.strptime(value, fmt).isoformat()
            except ValueError:
                continue
    except:
        pass
    return value

def parse_nested_keys(header):
    """Обработка точечной нотации: user.name -> {'user': {'name': ...}}"""
    result = {}
    for key in header:
        parts = key.split('.')
        current = result
        for part in parts[:-1]:
            if part not in current:
                current[part] = {}
            current = current[part]
        current[parts[-1]] = None
    return result

def build_nested_dict(row, keys):
    """Сборка вложенного словаря из плоского списка значений"""
    result = {}
    for i, key in enumerate(keys):
        parts = key.split('.')
        current = result
        for part in parts[:-1]:
            if part not in current:
                current[part] = {}
            current = current[part]
        current[parts[-1]] = parse_value(row[i]) if i < len(row) else None
    return result

def parse_csv(input_file, output_file=None, delimiter=None, nested=False):
    if not os.path.exists(input_file):
        print(f"❌ Файл {input_file} не найден.")
        return False

    if delimiter is None:
        delimiter = detect_delimiter(input_file)
        print(f"🔍 Определён разделитель: {repr(delimiter)}")

    data = []
    with open(input_file, 'r', encoding='utf-8-sig') as f:
        reader = csv.reader(f, delimiter=delimiter)
        rows = list(reader)
        if not rows:
            print("❌ CSV пуст.")
            return False
        headers = [h.strip() for h in rows[0]]
        print(f"📋 Заголовки: {headers}")

        for row in rows[1:]:
            if nested:
                # Вложенная структура
                # Создаём пустую структуру
                record = {}
                for i, header in enumerate(headers):
                    parts = header.split('.')
                    current = record
                    for part in parts[:-1]:
                        if part not in current:
                            current[part] = {}
                        current = current[part]
                    current[parts[-1]] = parse_value(row[i]) if i < len(row) else None
                data.append(record)
            else:
                # Плоская структура
                record = {}
                for i, header in enumerate(headers):
                    record[header] = parse_value(row[i]) if i < len(row) else None
                data.append(record)

    if output_file is None:
        output_file = os.path.splitext(input_file)[0] + '.json'

    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

    print(f"✅ Экспортировано {len(data)} записей в {output_file}")
    return True

def main():
    if len(sys.argv) < 2:
        print("Usage: python csv_parser.py <input.csv> [output.json] [--nested]")
        print("  --nested  : использовать вложенную структуру из точечной нотации")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = None
    nested = False

    for arg in sys.argv[2:]:
        if arg == '--nested':
            nested = True
        elif not arg.startswith('--'):
            output_file = arg

    print("📊 CSV → JSON Парсер (Python)")
    parse_csv(input_file, output_file, nested=nested)

if __name__ == "__main__":
    main()
