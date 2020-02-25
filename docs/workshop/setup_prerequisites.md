## Prerequisites and setup

> ⚠️ **Notice**: the following documents may contain deprecated functionalities that are still provided for backwards compatibility.

## Web Server for Chrome
Install [Web Server for Chrome](https://chrome.google.com/webstore/detail/web-server-for-chrome/ofhbbkphhbklhfoeikjpcbhemlocgigb?hl=en). It's useful for testing sample payloads and see how to process data locally before deploying.

## AWS Setup
1. Login to the AWS Console and select **Launch Instance**.
2. Select **Amazon Linux 2 AMI (HVM), SSD Volume Type**.
3. Select your instance size, **t2 micro** is fine.
4. Set the memory to **12GiB**.
5. Add tags.
6. Name the security group **Flex Demo** or similar.
7. Launch, then choose an existing keypair, or create one.
8. Your instance should spin up shortly.

## SSH Setup

To simplify SSHing into your AWS Instance, edit your `~/.ssh/config` file to include the following.

#### example ~/.ssh/config
```
Host flexdemo
 HostName 13.210.56.123
 User ec2-user
 IdentityFile ~/.ssh/flexDemo.pem
```
Modify the IP to your AWS instance, and point IdentityFile to where you have stored the key. You need to give your key the correct permissions:
```
chmod 400 ~/.ssh/flexDemo.pem
```

## Installing services for testing
1. Run `ssh flexdemo`. Enter "Yes" if prompted.
2. Follow the steps here to install the Infrastructure Agent

https://infrastructure.newrelic.com/accounts/YOUR_ACCOUNT_ID/install

3. Make sure that you use the command for Amazon Linux 2.
4. Ensure that the agent is reporting data back to New Relic.

```
### install redis & netcat
sudo amazon-linux-extras install redis4.0 -y
sudo systemctl start redis
```

```
### Setup nginx
sudo amazon-linux-extras install nginx1.12 -y
sudo systemctl start nginx

sudo nano /etc/nginx/nginx.conf
Within the server block, after this block
		location / {
		}

		# add the below
		
        location /nginx_status {
               	stub_status;
               	allow 127.0.0.1;        #only allow requests from localhost
               	deny all;               #deny all other hosts
       	}

Confirm your config file was edited correctly with
sudo nginx -t

If everything okay, restart nginx
sudo systemctl restart nginx

confirm nginx stub status metrics are being presented
curl http://localhost/nginx_status
```
```
### Setup redis prometheus integration (exporter)

wget https://github.com/oliver006/redis_exporter/releases/download/v0.32.0/redis_exporter-v0.32.0.linux-amd64.tar.gz

tar -xvf redis_exporter-v0.32.0.linux-amd64.tar.gz

### we can use the "screen" command to run it in the background

screen
./redis_exporter
press "ctrl+a", let go then press "d" (you'll be taken back, with the integration running in the background)

confirm prometheus metrics are being presented here
curl http://localhost:9121/metrics

```

## Kubernetes Setup

For the K8s setup, you can use whatever you like. Below are the instructions for using minikube or kops with AWS.

### Minikube for Local Setup (mac instructions)

1. Download the latest macos release of Minikube.
2. Install [homebrew](https://brew.sh/) and check it's there: `brew -v`.
3. Run `brew install redis`.
4. [Download the Redis Exporter](https://github.com/oliver006/redis_exporter/releases) for local testing (darwin release)
5. Run `brew cask install minikube`.

### KOPS

Note that you will be billed by AWS for having the following environment running, so clean up and delete the environment once you are done.

#### Create the cluster
```
# setup variables
export KOPS_CLUSTER_NAME=imesh.k8s.local
bucket_name=kops-flex-kav ### replace with your own name
KOPS_STATE_STORE=s3://${bucket_name}

# create cluster
kops create cluster --node-count=2 --node-size=t2.medium --zones=us-east-1a --name=${KOPS_CLUSTER_NAME}

# deploy cluster
kops update cluster --name ${KOPS_CLUSTER_NAME} --yes

# validate cluster is running (it can take a few minutes to be ready)
kops validate cluster
```

#### Delete the cluster
```
kops delete cluster --state s3://kops-flex-kav --name imesh.k8s.local --yes
```

