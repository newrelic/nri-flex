## Lab 0 - Install Flex

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility.

If you've followed the [prerequisites](setup_prerequisites.md), and have a Linux instance with the Infrastructure agent installed with some services to test against, you're ready to roll!

### Installation 
* Flex provides a one liner install for Linux.
* Alternatively, you can also deploy manually.

### Steps
1. Choose the install method:
    * [One liner](#one-liner-install)
    * [Manual](#manual-install)
2. [Confirm](#confirm)

#### One-liner install
```bash
sudo bash -c "$(curl -L https://newrelic-flex.s3-ap-southeast-2.amazonaws.com/install_linux_s3.sh)"

# Unpacked in /tmp/nri-flex-linux-$VERSION/...
```

#### Manual install

Download the latest Linux Flex package from [Releases](https://github.com/newrelic/nri-flex/releases).

```bash
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

You can confirm Flex is running and installed correctly if `flexStatusSample` are being generated in your account.

For example: https://insights.newrelic.com/accounts/YourAccountID/query?query=SELECT%20*%20FROM%20flexStatusSample%20
