# E2E features matrix

## General 

- [x] custom_attributes
- [] base_url
 
## APIs configuration

### Command 

| Name | Test |  
| ------ |------|
|                `run` | X    |
|              `shell` | TODO |
|              `split` | X    |
|           `split_by` | X    |
|        `regex_match` | X    |
|         `row_header` | TODO |
|          `row_start` | TODO |
|         `line_start` | TODO |
|           `line_end` | TODO |
|         `set_header` | X    |
| `header_regex_match` | TODO |
|    `header_split_by` | TODO |
|       `split_output` | TODO |
|            `timeout` | TODO |
|             `assert` | TODO |

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
| `jconvert_space`           | TODO |
| `ignore_output`            | TODO |
| `jq`                       | TODO |
| `keep_keys`                | TODO |
| `lazy_flatten`              | TODO |
| `lookup_file`               | TODO |
| `math`                      | TODO |
| `perc_to_decimal`           | TODO |
| `remove_keys`               | TODO |
| `rename_keys / replace_keys` | TODO |
| `sample_filter`             | TODO |
| `sample_include_filter`     | TODO |
| `sample_exclude_filter`     | TODO |
| `snake_to_camel`            | TODO |
| `split_array (leaf_array)`  | TODO |
| `split_objects`             | TODO |
| `start_key`                 | TODO |
| `store_variables`           | TODO |
| `lookups`                   | TODO |
| `dedupe_lookups`            | TODO |
| `store_lookups`             | TODO |
| `strip_keys`                | TODO |
| `timestamp`                 | TODO |
| `to_lower`                  | TODO |
| `value_parser`              | TODO |
| `value_transformer`         | TODO |
| `timestamp_conversion`      | TODO |