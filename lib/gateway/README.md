## User Guide


## Default JSON Mapping

| frugal   | JSON          | JSON example                     | Notes                                                                                                                                                                               |
| ------   | ----          | ------------                     | -----                                                                                                                                                                               |
| bool     | boolean       | true, false                      |                                                                                                                                                                                     |
| byte     | base64 string | "Y"                              | JSON value will be the data encoded as a string using standard base64 encoding with paddings. Either standard or URL-safe base64 encoding with/without paddings are accepted.       |
| i16, i32 | number        | 1, -10, 0                        | JSON value will be a decimal number. Either numbers or strings are accepted.                                                                                                        |
| i64      | string        | 1, -10, 0                        | JSON value will be a decimal string. Either numbers or strings are accepted.                                                                                                        |
| double   | number        | 1.1, -10.0, 0, "NaN", "Infinity" | JSON value will be a number or one of the special string values "NaN", "Infinity", and "-Infinity". Either numbers or strings are accepted. Exponent notation is also accepted.     |
| binary   | base64 string | "YWJjMTIzIT8kKiYoKSctPUBY+"      | JSON value will be the data encoded as a string using standard base64 encoding with paddings. Either standard or URL-safe base64 encoding with/without paddings are accepted.       |
| string   | string        | "Hello World!"                   |                                                                                                                                                                                     |
| list<V>  | array         | [v, …]                           | null is accepted as the empty list [].                                                                                                                                              |
| map<K,V> | object        | {"k": v, …}                      | All keys are converted to strings. Values are converted as per rules in this table.                                                                                                 |
| struct   | object        | {"fBar": v, "g": null, …}        | Generates JSON objects. Message field names are mapped to lowerCamelCase and become JSON object keys. A null value is treated as the default value of the corresponding field type. |
| enum     | string        | "FOO_BAR"                        | The enum value’s name as specified in the IDL.                                                                                                                                      |

1. Thrift allows map keys to be any type (including structs and
   exceptions). For maximum compatibility, only base types are allowed for
   map keys.
2. Sets cannot be serialized with the default marshaler, if you use sets,
   another marshaler must be used (to be developed).

## Custom JSON Mapping

You can customize the mapping for your types using annotations in the Thrift IDL.

## Interpreting Path Parameters

## Interpreting Query Parameters

1. For multiple query parameters, only the first is respected.
