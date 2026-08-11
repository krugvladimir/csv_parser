// csv_parser.rs — Rust версия

use serde_json::json;
use std::env;
use std::fs::File;
use std::io::{BufRead, BufReader, Write};

fn detect_delimiter(line: &str) -> char {
    let comma = line.matches(',').count();
    let semicolon = line.matches(';').count();
    let tab = line.matches('\t').count();
    if semicolon > comma && semicolon > tab {
        ';'
    } else if tab > comma && tab > semicolon {
        '\t'
    } else {
        ','
    }
}

fn parse_value(value: &str) -> serde_json::Value {
    let v = value.trim();
    if v.is_empty() {
        return serde_json::Value::Null;
    }
    match v.to_lowercase().as_str() {
        "true" => return serde_json::Value::Bool(true),
        "false" => return serde_json::Value::Bool(false),
        "null" | "none" => return serde_json::Value::Null,
        _ => {}
    }
    if let Ok(i) = v.parse::<i64>() {
        return serde_json::Value::Number(i.into());
    }
    if let Ok(f) = v.parse::<f64>() {
        if let Some(n) = serde_json::Number::from_f64(f) {
            return serde_json::Value::Number(n);
        }
    }
    serde_json::Value::String(v.to_string())
}

fn main() -> std::io::Result<()> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Usage: cargo run -- <input.csv> [output.json]");
        return Ok(());
    }

    let input_file = &args[1];
    let output_file = if args.len() > 2 {
        args[2].clone()
    } else {
        "output.json".to_string()
    };

    println!("📊 CSV → JSON Парсер (Rust)");

    let file = File::open(input_file)?;
    let reader = BufReader::new(file);
    let lines: Vec<String> = reader
        .lines()
        .filter_map(|line| {
            let l = line.ok()?;
            if l.trim().is_empty() { None } else { Some(l) }
        })
        .collect();

    if lines.is_empty() {
        println!("❌ CSV пуст.");
        return Ok(());
    }

    let delimiter = detect_delimiter(&lines[0]);
    println!("🔍 Определён разделитель: {}", delimiter);

    let headers: Vec<String> = lines[0]
        .split(delimiter)
        .map(|h| h.trim().to_string())
        .collect();
    println!("📋 Заголовки: {:?}", headers);

    let mut data = Vec::new();
    for line in &lines[1..] {
        let values: Vec<&str> = line.split(delimiter).map(|v| v.trim()).collect();
        let mut record = serde_json::Map::new();
        for (i, header) in headers.iter().enumerate() {
            let val = if i < values.len() { values[i] } else { "" };
            record.insert(header.clone(), parse_value(val));
        }
        data.push(serde_json::Value::Object(record));
    }

    let json = serde_json::to_string_pretty(&data)?;
    let mut file_out = File::create(&output_file)?;
    file_out.write_all(json.as_bytes())?;

    println!("✅ Экспортировано {} записей в {}", data.len(), output_file);
    Ok(())
}
