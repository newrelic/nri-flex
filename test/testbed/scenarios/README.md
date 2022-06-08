# E2E features matrix

## General 

- [x] custom_attributes
- [] base_url
 
## APIs configuration

### Command 

| Name | Test |  
| ------ |------|
|                `run` | X    |
|              `shell` | X    |
|              `split` | X    |
|           `split_by` | X    |
|        `regex_match` | X    |
|         `row_header` | X    |
|          `row_start` | X    |
|         `line_start` | X    |
|           `line_end` | X    |
|         `set_header` | X    |
| `header_regex_match` | TODO |
|    `header_split_by` | X    |
|       `split_output` | TODO |
|            `timeout` | TODO |
|             `assert` | X    |

### URL

| Name             | Test |  
|------------------|------|
| `headers`        | X    |
| `tls config`     | TODO |

### File

| Name          | Test |  
|---------------|------|
| `set_headers` | TODO |

## Supported transformation functions

| Name                       | Test |  
|----------------------------|------|
| `add_attributes`           | X    |
| `convert_space`           | TODO |
| `ignore_output`            | X    |
| `jq`                       | X    |
| `keep_keys`                | X    |
| `lazy_flatten`              | X    |
| `lookup_file`               | TODO |
| `math`                      | X    |
| `perc_to_decimal`           | X    |
| `remove_keys`               | X    |
| `rename_keys / replace_keys` | X    |
| `sample_filter`             | X    |
| `sample_include_filter`     | X    |
| `sample_exclude_filter`     | X    |
| `snake_to_camel`            | X    |
| `split_array (leaf_array)`  | TODO |
| `split_objects`             | X    |
| `start_key`                 | X    |
| `store_variables`           | TODO |
| `lookups`                   | TODO |
| `dedupe_lookups`            | TODO |
| `store_lookups`             | TODO |
| `strip_keys`                | X    |
| `timestamp`                 | TODO |
| `to_lower`                  | X    |
| `value_parser`              | X    |
| `value_transformer`         | X    |
| `timestamp_conversion`      | TODO |