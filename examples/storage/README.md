# Storage Integrations

Collection of flex integrations for the following Storage technologies:

* [Cohesity](https://developer.cohesity.com/apidocs-641.html#/rest/getting-started)
* [Dell Powerscale](https://developer.dell.com/apis/4088/versions/9.2.0.0/docs/1introduction.md)
* [Dell Unity](https://developer.dell.com/apis/3028/versions/5.1.0/models/spec_publish.yml)
* [Dell VMAX](https://infohub.delltechnologies.com/l/dell-emc-powermax-and-vmax-all-flash-snapshot-policies-2/appendix-b-rest-api-examples-2)
* [Dell VPLEX](https://www.delltechnologies.com/asset/en-us/products/storage/technical-support/h18625-dell-vplex-restv2-transition-best-practices-guide.pdf)
* [Pure Storage](https://pure-storage-python-rest-client.readthedocs.io/en/stable/api.html)


The majority of these examples are Flex yaml configs that parse a lookup file (iterates through a list of hosts or clusters) and execute a python script that poll each technology's REST API, using the lookup file's attributes as arguments passed into the python script.

## Getting Started
1. Modify a given json file (lookup file) with a list of hosts/clusters to monitor
2. Modify `INSIGHTS_API_KEY` and `INSIGHTS_URL` with your api key, and add accountId to the string in the url. **NOTE: This is required to remove any entityGuid stamp on the telemetry sent, so entity synthesis can properly take place.**
3. Generate a secret for password - This can be done with Flex via a command such as:

```
./nri-flex -encrypt_pass 'password=****** -pass_phrase 'N3wR3lic!'
```

...or via one of the options [documented here](https://docs.newrelic.com/docs/infrastructure/host-integrations/installation/secrets-management/)


## Entity Synthesis
Each of these integrations will generate unique entities with out of box visualizations, if the eventTypes and telemetry are not customized. To validate the entity synthesis definitions for each technology, please reference them below:

* [Cohesity](https://github.com/newrelic/entity-definitions/tree/main/entity-types/ext-cohesity)
* [Dell Unity VNX](https://github.com/newrelic/entity-definitions/tree/main/entity-types/ext-dell_vnx)
* [Dell VMAX](https://github.com/newrelic/entity-definitions/tree/main/entity-types/ext-dell_vmax)
* [Dell VPLEX](https://github.com/newrelic/entity-definitions/tree/main/entity-types/ext-dell_vplex)
* [Pure Storage](https://github.com/newrelic/entity-definitions/tree/main/entity-types/ext-pure)
