#!/bin/bash -xe
# Configure SSH to listen on port 443 so we can use
# the proxy to access them.
sed -i 's/\#Port 22/Port 22\nPort 443/' /etc/ssh/sshd_config
service sshd restart || service ssh restart
yum install -y git
pip install ansible
# Setup 2-factor authentication.
# Clone the Google Authenticator project.
sudo yum install -y google-authenticator
# Follow the instructions at: https://www.conspire.com/blog/2014/09/two-factor-auth-on-ec2-with-public-key-and-google/
# (From the point of installing the PAM module in /etc/ssh/sshd_config)
