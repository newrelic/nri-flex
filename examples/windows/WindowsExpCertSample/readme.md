## Overview:
This is an example Flex to retrieve local certificates on a windows machine. The core of this flex is based on the following: 

```
get-childitem -path cert: -recurse -Expiringindays 30
```

You can change the Expiringindays number if you would like the script to pull in more or fewer days.

For installation and configuration of flex please see the [main Flex repo page and docs](https://github.com/newrelic/nri-flex).