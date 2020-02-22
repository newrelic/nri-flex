### Net Dial

Dial is a parameter used under commands.

port test eg.
```
name: portTestFlex
apis: 
  - timeout: 1000 ### default 1000 ms increase if you'd like
    commands:
    - dial: "google.com:80"
```

sending a message and processing the output eg.
```
---
name: redisFlex
apis: 
  - name: redis
    commands: 
      - dial: 127.0.0.1:6379
        run: "info\r\n"
        split_by: ":"
```