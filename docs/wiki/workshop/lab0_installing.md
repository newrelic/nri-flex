## Lab 0 - Install Flex

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility. However, an updated version of this document is in progress. 

If you've followed the prereqs, and have a linux instance with the infrastructure Agent installed with some services to test against, you're ready to roll.
* Downloading the darwin(mac) release for local testing can also be quite helpful.

### Installation 
* Flex provides a one liner install for Linux.
* Alternatively you can also deploy manually.

### Steps
* Install
    * [One Liner](#one-liner-install) or;
    * [Manual](#manual-install)
* [Confirm](#confirm)

#### One Liner Install
```
sudo bash -c "$(curl -L https://newrelic-flex.s3-ap-southeast-2.amazonaws.com/install_linux_s3.sh)"

# Unpacked in /tmp/nri-flex-linux-$VERSION/...
```

#### Manual Install

* Download the latest Linux Flex package from the Releases section (do not download the source code zip)
* https://github.com/newrelic/nri-flex/releases

```
# copy file over scp to your instance
eg. (note: could be a different version)
scp nri-flex_linux-v0.7.7-pre.tar flexdemo:/home/ec2-user 

SSH to instance
ssh flexdemo

Extract package
tar -xvf nri-flex_linux-v0.7.7-pre.tar

Enter directory
cd nri-flex_linux-v0.7.7-pre/

Install Flex
sudo ./install_linux.sh

```

#### Confirm
* You can confirm Flex is running and installed correctly, if "flexStatusSample"s are being generated in your account.

* eg. https://insights.newrelic.com/accounts/YourAccountID/query?query=SELECT%20*%20FROM%20flexStatusSample%20
