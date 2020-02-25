## Lab 2 - Add Flex configurations

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility.

Flex allows you to add multiple configurations in `/flexConfigs`, unless you have specified a different location in the `nri-flex-config.yml` file.

Note that:
* You don't need to restart Flex when adding new Flex config file/integration.
* Each configuration file is processed asynchrously, whereas the content of the configuration files is processed synchronously.  

### Selecting example configurations

We are going to pick some simple configurations to show how simply it is to add new integrations.

Some examples:
- df
- systemctl
- nginx
- redis

There's several approaches to adding configurations:

* The [Flex release package](#Use-config(s)-that-exists-within-the-release-package) bundles all available examples so if what you need is there, you can simply copy it over.
* You can [copy / SCP](Copy-the-files-locally-over-to-the-instance-via-SCP) the config files you need.
* You can [paste in the terminal](#Copy-configs-via-an-editor) the config file itself using editors like vi or nano.
* Finally, you can also [extract the Flex tarball, add/keep the configs you want, and then repackage it](#Bundle-specific-configs-into-a-new-package).

___
## Use config(s) that exists within the release package

Change directory to wherever you have extracted the Flex release package; in this case it's extracted within the home directory:

```bash
cd /home/ec2-user/nri-flex_linux-v0.7.7-pre/examples/flexConfigs
```

Copy or move the file to `/var/db/newrelic-infra/custom-integrations/flexConfigs`:

```bash
sudo mv linux-systemctl-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs
```
```bash
sudo cp linux-systemctl-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs
```

Confirm that you have some new events; there is no need to restart the agent:

https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/query?query=FROM%20dfSample%20SELECT%20*

## Copy the files locally over to the instance via SCP

Copy the file to an instance using SCP:
```bash
scp your-config-file.yml flexdemo:/home/ec2-user/
```
Move the file to the right directory:
```bash
ssh flexdemo sudo mv /home/ec2-user/df-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs/
```
Note that you can't SCP directly to a directory owned by root if you're not logged as root.

Confirm in Insights if you have some new events; there is no need to restart the agent: 

https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/query?query=FROM%20dfSample%20SELECT%20*

## Copy configs via an editor

Given that you have spun up an Amazon Linux or Linux instance with SystemD, we can use the following example:
```bash
examples/flexConfigs/linux-systemctl-cmd-example.yml
```

Cat or open the file in an editor, then copy the contents. In your Flex demo environment, run:
```bash
ssh flexdemo
cd /var/db/newrelic-infra/custom-integrations/flexConfigs/
sudo nano linux-systemctl.yml
```
Paste the contents, then save the file (**Ctrl + X**,  then **Y** and **Enter**)

## Bundle specific configs into a new package

Extract the Flex release package.
```bash
tar -xvf tar -xvf nri-flex_linux-v0.7.7-pre.tar
````
Navigate to the newly created folder:
```bash
cd nri-flex_linux-v0.7.7-pre/
```
Add your own flexConfigs folder:
```bash
mkdir flexConfigs
```
Within the newly created flexConfigs folder, add the specific configs you want. 

Then modify the `install_linux.sh` script. This will copy all configs, only keep what you need.

```bash
cp -avr ./flexConfigs /var/db/newrelic-infra/custom-integrations/
```

Finally, repackage your modified release by moving a directory up 
```bash
cd ../
tar -cvf nri-flex_linux-v0.7.7-pre-my-modified-pkg.tar nri-flex_linux-v0.7.7-pre/
```
Now you can follow the standard steps to [install Flex](w-Lab0-Installing-Flex.md).