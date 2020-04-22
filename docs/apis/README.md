# Flex APIs

Flex APIs give you access to several data sources, such as HTTP endpoints and raw text output by command-line utilities.

Supported APIs include:

* [commands](commands.md): Retrieve information from any application or shell command. 
* [url](url.md): Retrieve information from any HTTP endpoint.

Flex APIs can be used together in [configuration files](../basics/configure.md). For example:

```yaml
integrations:
 - name: nri-flex # We're telling the Infra agent to run Flex!
   config: # Flex configuration starts here!
     name: sample_source
     apis:
       - event_type: apiCallSample # Name of the event in New Relic
         url: https://jsonplaceholder.typicode.com/posts
       - event_type: commandSample # Name of the event in New Relic
         commands:
           - run: du -c /somedir
             set_header: [dirSizeBytes,dirName]
             regex_match: true
             split: horizontal
             split_by: (\d+)\s+(.*)
           - run: some_other_command
             split_by: \s+
```
The following APIs are still experimental. 'Experimental' here means that New Relic does not yet provides support for them:

- [Database queries](experimental/db.md)
- [Net dial](experimental/dial.md)
- [Git configuration synchronization](experimental/git_sync.md)
