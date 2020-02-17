## Understanding What Flex Is Doing

> ⚠️ **Notice** ⚠️: the following documents may contain some deprecated functionalities that
are still provided for backwards compatibility. However, an updated version of this
document is in progress. 


So you've added some configs, but now you need to know:
- is it working?
- where is the data?
- how much is it doing?
- what is it doing?

#### Checking the Flex Status Sample

The simplest way is to check the Flex Status Sample, as you've done during the confirm step during install.

Update with your account id:
* https://insights.newrelic.com/accounts/YourAccountID/query?query=SELECT%20*%20FROM%20flexStatusSample%20

You will notice within the sample, that counters are kept for the number of events being generated and also what events are being generated. 

#### What events are being created?

So if you've checked the status sample, or the config files you would have deploy you would see the event types being created.

You can then query those same event types, and build dashboards from them as you usual would.

eg. SELECT * FROM OneOfTheNewEventTypes 

#### What is Flex doing?

An easy way to see what Flex is doing each time it runs, is to run it directly from where it is installed.

```
# change directory
cd /var/db/newrelic-infra/custom-integrations/

# execute
sudo ./nri-flex -verbose -pretty
```

Note if you have added any additional parameters in your `nri-flex-config.yml` you will need to pass them through the command line as well to simulate the correct behaviour.