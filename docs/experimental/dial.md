### Net dial

> **Disclaimer**: this function is bundled as alpha. That means that it is not yet supported by New Relic.

Dial is a parameter used under `commands`.

Example of a port test:
```yaml
name: portTestFlex
apis: 
  - timeout: 1000 ### default 1000 ms increase if you'd like
    commands:
    - dial: "google.com:80"
```

Example of sending a message and processing the output:
```yaml
---
name: redisFlex
apis: 
  - name: redis
    commands: 
      - dial: 127.0.0.1:6379
        run: "info\r\n"
        split_by: ":"
```