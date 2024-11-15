# Exporters

An **exporter** is a tool that transforms all the documents from a processor into a file. There are three available
formats for exporting: CSV, JSON, and YAML. Each format structures the data in a different way, depending on how you
need to use or store the data.

## Available Export Formats

### 1. CSV

The CSV exporter returns two columns:

- **pkey**: The primary key of each document.
- **data**: The document itself, stored in JSON format.

This format is ideal for users who need a simple, spreadsheet-compatible export, with each document's primary key and
JSON data in separate columns.

Example CSV output:
```csv
pkey,data
"[""799""]","{""name"":""John Doe"",""age"":30}"
"[""780""]","{""name"":""Jane Doe"",""age"":25}"
```
### 2. JSON

The JSON exporter returns an array of objects, where each object contains two fields:

- **pkey**: The primary key of the document.
- **data**: The document in its raw JSON format.

This format is perfect for users who want to use the exported data directly in applications or databases that support
JSON.

Example JSON output:

```json
[
  {
    "pkey": "[\"1\"]",
    "data": {
      "name": "John Doe",
      "age": 30
    }
  },
  {
    "pkey": "[\"2\"]",
    "data": {
      "name": "Jane Doe",
      "age": 25
    }
  }
]
```

### 3. YAML

The YAML exporter also returns an array of objects, structured similarly to the JSON format:

- **pkey**: The primary key of the document.
- **data**: The document in its raw YAML format.

YAML is a more human-readable format, suitable for configurations or settings files that need to be edited manually.

Example YAML output:

```yaml
- pkey: "[\"1\"]"
  data:
    name: John Doe
    age: 30
- pkey: "[\"2\"]"
  data:
    name: Jane Doe
    age: 25
```

## Choosing the Right Export Format

* Use `CSV` if you need a simple tabular format that's compatible with spreadsheets or data analysis tools.
* Use `JSON` for structured data that will be processed by systems or applications.
* Use `YAML` for a more human-readable format, suitable for configuration or settings purposes.