# Install terraform
wget -nc https://releases.hashicorp.com/terraform/0.6.14/terraform_0.6.14_linux_amd64.zip
sudo yum install -y unzip
unzip terraform_0.6.14_linux_amd64.zip
sudo mv terraform-* /usr/local/bin
sudo mv terraform /usr/local/bin
rm terraform_0.6.14_linux_amd64.zip
# Install awscli
sudo yum install -y python-setuptools
sudo easy_install pip
sudo pip install awscli
# Ensure that the clock is set properly.
# Sync the time
yum install -y chrony
sudo systemctl enable chronyd
sudo systemctl start chronyd
