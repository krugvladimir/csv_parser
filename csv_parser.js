// csv_parser.js — JavaScript версия

const fs = require('fs');
const path = require('path');

function detectDelimiter(firstLine) {
    const delimiters = [',', ';', '\t', '|'];
    let bestDelim = ',';
    let maxFields = 0;
    for (const delim of delimiters) {
        const fields = firstLine.split(delim);
        if (fields.length > maxFields) {
            maxFields = fields.length;
            bestDelim = delim;
        }
    }
    return bestDelim;
}

function parseValue(value) {
    value = value.trim();
    if (value === '') return null;
    if (value.toLowerCase() === 'true') return true;
    if (value.toLowerCase() === 'false') return false;
    if (value.toLowerCase() === 'null' || value.toLowerCase() === 'none') return null;
    if (!isNaN(value) && value !== '') {
        if (value.includes('.')) return parseFloat(value);
        return parseInt(value, 10);
    }
    return value;
}

function parseCSV(inputFile, outputFile = null) {
    if (!fs.existsSync(inputFile)) {
        console.error(`❌ Файл ${inputFile} не найден.`);
        return false;
    }

    const content = fs.readFileSync(inputFile, 'utf-8');
    const lines = content.split('\n').filter(line => line.trim() !== '');

    if (lines.length === 0) {
        console.error('❌ CSV пуст.');
        return false;
    }

    const delimiter = detectDelimiter(lines[0]);
    console.log(`🔍 Определён разделитель: ${delimiter === '\t' ? '\\t' : delimiter}`);

    const headers = lines[0].split(delimiter).map(h => h.trim());
    console.log(`📋 Заголовки: ${headers.join(', ')}`);

    const data = [];
    for (let i = 1; i < lines.length; i++) {
        const row = lines[i].split(delimiter);
        const record = {};
        for (let j = 0; j < headers.length; j++) {
            const val = j < row.length ? row[j] : '';
            record[headers[j]] = parseValue(val);
        }
        data.push(record);
    }

    if (!outputFile) {
        outputFile = path.basename(inputFile, path.extname(inputFile)) + '.json';
    }

    fs.writeFileSync(outputFile, JSON.stringify(data, null, 2), 'utf-8');
    console.log(`✅ Экспортировано ${data.length} записей в ${outputFile}`);
    return true;
}

if (process.argv.length < 3) {
    console.log('Usage: node csv_parser.js <input.csv> [output.json]');
    process.exit(1);
}

const inputFile = process.argv[2];
const outputFile = process.argv[3] || null;

console.log('📊 CSV → JSON Парсер (JavaScript)');
parseCSV(inputFile, outputFile);
