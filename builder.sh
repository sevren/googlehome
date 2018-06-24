#!/bin/sh

###############################
#                             #
# GoogleHomeKodi+SSL Builder  #
#                             #
###############################


if [ "$#" -ne 2 ]; then
	echo "Required 2 arguments DNS_HOST and EMAIL"
	echo "Example using duckdns host Run: ./builder xxxx.duckdns.org xxx@email-domain.example.com"
	exit 1
fi

DNS_HOST_VAR=$1
EMAIL=$2

SCRIPT_DIR=$(dirname $0)

cd googlehomekodi
echo "Cloning out the Master branch for GoogleHomeKodi Node project"
git clone https://github.com/OmerTu/GoogleHomeKodi.git
echo "Replacing the Dockerfile with our own"
cp Dockerfile ./GoogleHomeKodi/Dockerfile
echo "Building and tagging the GoogleHomeKodi docker image"
cd GoogleHomeKodi
docker build -t sinedied/googlehomekodi .
cd ${SCRIPT_DIR}
cd nginx-certbot-proxy
echo "Creating the diffe-helman key.. this will take a while go grab a coffee"
openssl dhparam -out ./dhparams.pem 2048
echo "Now editing the ngnix configuration, inserting DNS_HOST into nginx configuration automatically"
NGNIX_SERVER_NAME_VAR=$(awk '/server_name/{print $2}' nginx | tr -d ';')  # Get the current variable in the server_name configuration and remove the semi-colon at the end
sed -i'' s/${NGNIX_SERVER_NAME_VAR}/$DNS_HOST_VAR/g  nginx #match the occurrence of the server_name variable and modify inplace with the value used in the DNS_HOST parameter
docker build -t sslstuff . --build-arg DNS_HOST=$DNS_HOST_VAR --build-arg YOUR_EMAIL=$EMAIL
cd ${SCRIPT_DIR}


