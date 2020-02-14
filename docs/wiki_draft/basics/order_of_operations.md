# Flex Order of Operations

- Find Start Key
- Strip Keys - Happens before attribute modifiction and auto flattening, useful to get rid of unneeded data and arrays early
- Lazy Flatten
- Standard Flatten
- Remove Keys
- Strip Keys (second round)
- Merge (if used)
- ToLower Case
- Convert Space
- snake_case to camelCase
- Value Parser
- Pluck numbers
- Sub parse
- Value Transformer
- Rename Key // uses regex to find keys to replace
- Store Lookups
- Keep Keys // keeps only keys you want to keep, and removes the rest