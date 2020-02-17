## Lab 2 - Adding Configs

> ⚠️ **Notice** ⚠️: the following documents may contain some deprecated functionalities that
are still supported by New Relic for backwards compatibility. However, an updated version of this
document is in progress. 

* Flex allows you to add as many configs which are essentially integrations into the "flexConfigs" folder, unless you have specific a different location in the nri-flex-config.yml file.
* Note that:
    * you do not need to restart Flex when adding new Flex Config files/integrations
    * each individual config file is processed asynchrously, whereas everything in a config file itself is processed synchronously. This is intended to give you better control and performance.



### Selecting Example Configs
- To start we can pick some simple configurations to demonstrate the simplicity of adding integrations.
    - Note this can be further simplified by the Git Syncing capability Flex provides.
- Examples:
    - df
    - systemctl
    - nginx
    - redis

To add configs, we have several options to approach this:
* [The Flex release package bundles all available examples so if what you need is there, you can simply copy it over.](#Using-config(s)-that-exists-within-the-release-package)
* [Copy / scp the config files you need over.](Copying-the-files-locally-over-to-the-instance-via-SCP)
* [In the terminal just paste the config file itself via and editor like vi or nano.](#Copying-configs-via-an-editor)
* [Extract the Flex tarball prior, add/keep the configs you want then repackage it (useful for packaging it with specific for a customer), and perform the standard installation.](#Bundling-specific-configs-into-a-new-package)

___
## Using config(s) that exists within the release package
```
change directory to wherever you have extracted the Flex release package
In this case it's extracted within the home directory

cd /home/ec2-user/nri-flex_linux-v0.7.7-pre/examples/flexConfigs
copy or move the file to -> /var/db/newrelic-infra/custom-integrations/flexConfigs
eg.

Move File
sudo mv linux-systemctl-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs

or Copy File
sudo cp linux-systemctl-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs
```

Confirm in Insights if you have some new events, there is no need to restart the agent.
https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/query?query=FROM%20dfSample%20SELECT%20*

---

## Copying the files locally over to the instance via SCP

```
copy file to instance
scp your-config-file.yml flexdemo:/home/ec2-user/

move file to correct directory
ssh flexdemo sudo mv /home/ec2-user/df-cmd-example.yml /var/db/newrelic-infra/custom-integrations/flexConfigs/

### note: unless using a root, account you can't scp directly to a directory owned by root
```

Confirm in Insights if you have some new events, there is no need to restart the agent.
https://insights.newrelic.com/accounts/YOUR_ACCOUNT_ID/query?query=FROM%20dfSample%20SELECT%20*

---

## Copying configs via an editor
```
Given that you have spun up an Amazon Linux or Linux instance with SystemD we can use the following example (this is both within the tarball and the repository) (else use the DF example)

examples/flexConfigs/linux-systemctl-cmd-example.yml
cat or open the file in an editor, and copy the contents
```
```
In your flex demo environment
ssh flexdemo
cd /var/db/newrelic-infra/custom-integrations/flexConfigs/

sudo nano linux-systemctl.yml    (or whatever you want to name it)
paste the contents, and save the file ("ctrl + x",  then "y" and hit "enter")
```

## Bundling specific configs into a new package

```
Extract the Flex release package.
tar -xvf tar -xvf nri-flex_linux-v0.7.7-pre.tar (or whichever version you're using)

Enter the newly created folder
cd nri-flex_linux-v0.7.7-pre/
```
```
Add your own flexConfigs folder within there
mkdir flexConfigs
Within the newly created flexConfigs folder, add the specific configs you want
```
```
Modify the install_linux.sh script

Within that script, you will see this portion commented out
# this will copy all configs, only take what you need
# cp -avr ./flexConfigs /var/db/newrelic-infra/custom-integrations/

uncomment the line that begins with "cp"
```
```
Repackage your modified release
move a directory up
cd ../
tar -cvf nri-flex_linux-v0.7.7-pre-my-modified-pkg.tar nri-flex_linux-v0.7.7-pre/
```

* [Then follow standard steps to install Flex if you haven't already.](w-Lab0-Installing-Flex.md)