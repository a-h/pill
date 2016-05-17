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
