# csv_parser.rb — Ruby версия

require 'csv'
require 'json'

def detect_delimiter(first_line)
  comma = first_line.count(',')
  semicolon = first_line.count(';')
  tab = first_line.count("\t")
  if semicolon > comma && semicolon > tab
    ';'
  elsif tab > comma && tab > semicolon
    "\t"
  else
    ','
  end
end

def parse_value(value)
  value = value.to_s.strip
  return nil if value.empty?
  return true if value.downcase == 'true'
  return false if value.downcase == 'false'
  return nil if value.downcase == 'null' || value.downcase == 'none'
  if value =~ /^[-+]?\d+$/
    return value.to_i
  end
  if value =~ /^[-+]?\d+\.\d+$/
    return value.to_f
  end
  value
end

def parse_csv(input_file, output_file = nil)
  unless File.exist?(input_file)
    puts "❌ Файл #{input_file} не найден."
    return false
  end

  content = File.read(input_file)
  lines = content.split("\n").reject { |l| l.strip.empty? }

  if lines.empty?
    puts "❌ CSV пуст."
    return false
  end

  delimiter = detect_delimiter(lines[0])
  puts "🔍 Определён разделитель: #{delimiter == "\t" ? '\\t' : delimiter}"

  headers = lines[0].split(delimiter).map(&:strip)
  puts "📋 Заголовки: #{headers.join(', ')}"

  data = []
  lines[1..-1].each do |line|
    values = line.split(delimiter).map(&:strip)
    record = {}
    headers.each_with_index do |header, i|
      record[header] = parse_value(values[i] || '')
    end
    data << record
  end

  output_file ||= File.basename(input_file, '.*') + '.json'
  File.write(output_file, JSON.pretty_generate(data))
  puts "✅ Экспортировано #{data.size} записей в #{output_file}"
  true
end

if ARGV.length < 1
  puts "Usage: ruby csv_parser.rb <input.csv> [output.json]"
  exit 1
end

input = ARGV[0]
output = ARGV[1]

puts "📊 CSV → JSON Парсер (Ruby)"
parse_csv(input, output)
