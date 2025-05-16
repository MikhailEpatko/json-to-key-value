# json-to-key-value
JSON to key-value parser

1. download JSON from database tables
2. parse JSON to key-value pairs
3. skip trash data
4. write result to json files (a file per table)

JSON example:

```json
{
  "title": "value1-1",
  "type": "value1-2",
  "key2": {
    "text": "value 2-1",
    "description": "value 2-2"
  },
  "color": "#FFFFFF"
}
```

Result example:

```json
{
  "pairs": [
    {"id.title": "value1"},
    {"id.key2.text": "value 2-1"},
    {"id.key2.description": "value 2-2"}
  ]
}
```

'id' - is table row id.

More info: https://ezpkg.io/iter.json/#pkgs