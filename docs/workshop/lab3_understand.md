# Understand What Flex Is Doing

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility.

You've added some configs, but now you need to know:

- Is it working?
- Where is the data?
- What is it doing?

## Check the Flex Status Sample

The simplest way to know that Flex is working is checking the Flex Status Sample, as you've done during the confirm step during install.

https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/query?query=SELECT%20*%20FROM%20flexStatusSample%20

Notice that counters within the sample indicate the number of events being generated and also what events are being generated. 

## What events are being created?

Query the event types you've configured, and build dashboards from them as you usually would:

```sql
SELECT * FROM OneOfTheNewEventTypes 
```

## What is Flex doing?

An easy way to see what Flex is doing each time it runs, is to run it directly from where it is installed:

```bash
# change directory
cd /var/db/newrelic-infra/custom-integrations/

# execute
sudo ./nri-flex -verbose -pretty
```

If you have added any additional parameters in `nri-flex-config.yml` you need to pass them through the command line as well to simulate the correct behaviour.