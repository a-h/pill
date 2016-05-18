# pill
A skills matrix Web application, written in Go and targeting MongoDB for data storage. Authentication handled by Google.

# Building a local Docker Image
To build a local Docker image:

`docker build -t pill .`

Proxies are supported by the use of the `--build-arg` parameter:

`docker build --build-arg http_proxy=http://10.0.2.2:3128 --build-arg https_proxy=http://10.0.2.2:3128 -t pill .`

# Docker Compose
To run the stack locally, use the `docker-compose up` command which uses the `docker-compose.yml` file to create the infrastructure.

# Rebuilding
To rebuild the application stack in the container, use `docker-compose build pill`.

# Accessing the Website.
* The `docker-compose.yml` contains a port-forwarding rule to forward 8080 on the container host to port 8080 on the guest.
  * If you're running on Windows, you will also need to configure VirtualBox to setup port-forwarding on your boot2docker VirtualBox instance.
* The default connection string is to connect to the `mongo` service on `mongodb://mongo:27017`.

# Configuring the service in AWS.
* Setup Elastic Container Service (AWS ECS) from the console in
* Setup an Elastic Container Repository (AWS ECR) in AWS.
* Create a role you will use to manage the "pill" service.
  * It will be a "Role for Cross-Account Access"
* Attach the following policies to the role:
  * AmazonEC2ContainerRegistryFullAccess
  * AmazonEC2ContainerServiceRole
* Assign your management user to the role.

# Build the container and put it in the repository
* Login
  * `aws ecr get-login --region eu-west-1`
    * If you get a missing pipe, you're probably using the Windows command shell, switch to the docker quickstart terminal.
    * If you have trouble with proxies, see https://github.com/docker/toolbox/issues/102#issuecomment-154900769
* Build
  * `docker build -t pill .` (see above)
* Tag
  * `docker tag pill:latest 180466524585.dkr.ecr.us-east-1.amazonaws.com/pill:latest`
* Push
  * `docker push 180466524585.dkr.ecr.us-east-1.amazonaws.com/pill:latest`
