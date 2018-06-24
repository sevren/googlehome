#!/bin/sh

###########################
#                         #
# GoogleHomeKodi Runner   #
#                         #
###########################


# First attempt to create a network since all containers need to be able to see each other
echo "Attempting to create the internal docker network"
docker network create ghome || true

#Run the GoogleHomeKodi container first (ensure the configuration file exists in a folder which you mount)
echo "Now running the GoogleHomeKodi node application in a docker container, subsequent reboots of the system will automatically spin up this container"
docker run --network=ghome  -d  -p 8099:8099 --restart always -v $PWD/googlehomekodi:/config --name googlehomekodi sinedied/googlehomekodi

# Run the Proxy server
echo "Now running the Nginx web server in a docker container, subsequent reboots of the system will automatically spin up this container"
docker run  --network=ghome -d -p 8443:443 --restart always --name nginx sslstuff
