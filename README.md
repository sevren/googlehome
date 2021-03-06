# Googlehome Projects repo

The following repo is a collection of scripts and docker images for Google Home automation.

Presently these scripts are running on a RPI model 3

## Repo Structure

* GoogleHomeKodi Folder (Dockerfile and configuration file)
* SSLStuff Folder (Dockerfile and ngnix configuration file)

## GoogleHome Kodi Integration

Thanks and all the credit goes to OmerTu. see https://github.com/OmerTu/GoogleHomeKodi for his fantastic work on the NodeJs application

We can setup Google Home to work with Kodi version 17.3+. This setup has been tested with Libreelec.

This setup conists of 2 running docker containers on the raspberry PI. 
1. Running the Google Home Kodi Node Application, listening for requests on port 8099 (by default)
2. Running a preconfigured Nginx webserver to redirect url requests from IFFTT into the Google Home Kodi docker container

## Building and running the Google Home Kodi stack automatically

*OBS!: BEFORE YOU RUN THE BUILDER PLEASE UPDATE THE CONFIGURATION FILE FOUND IN googlehomekodi FOLDER!*

*OBS!: BEFORE YOU RUN THE BUILDER PLEASE ENSURE THAT YOU HAVE CREATED A DNS HOST (SEE nsupdate.info or duckdns.org)*

*OBS!: THIS BUILDER WILL AUTOMATICALLY CONFIGURE THE ngnix FILE IN sslstuff WITH THE NAME PROVIDED AS THE FIRST PARAMETER IN THE SCRIPT*


From the root folder there is two scripts
1. builder.sh 
2. runghkodi.sh

You can use `builder.sh` to automatically go into each folder and build/tag the docker images.

```bash
./builder.sh ghome.duckdns.org yourEmailAddr@email.domain.com
```

Once everything has completed successfully you can run the stack using the following:

```bash
./runghkodi.sh
```

### Building && Running the Google Home Kodi container (Seperately)

A modified dockerfile which runs on Armv7 can be found in the googlehomekodi folder.

The docker container is built using the following command:

We clone out the GoogleHomeKodi git repo and replace the Dockerfile with our own.
Probably we can just point context to the GoogleHomeKodi folder instead of replacing the dockerfile..

```bash
cd googlehomekodi
git clone https://github.com/OmerTu/GoogleHomeKodi.git
cp Dockerfile ./GoogleHomeKodi
cd GoogleHomeKodi
docker build -t sinedied/googlehomekodi .
```
And can be run using the following:

```bash
docker run --network=ghome  -d  -p 8099:8099 --restart always -v /storage/googlehomekodiConfig:/config --name googlehomekodi sinedied/googlehomekodi
``` 

### SSL setup for googlehome kodi

Since the code for GoogleHomeKodi is not set up for SSL. We can setup a nginx proxy in a docker container and forward requests
to the same docker network. The proxy will enable us to use SSL through the use of a letsencrypt certificate.

#### DNS Setup

1. Signup here: https://www.nsupdate.info/
2. Login into your account and click Overview
3. Add Host and choose a proper subdomain name (remember this domain as you will need it for your webhooks)
4. Update your router settings to use this DynDNS with the host created

OR you can use https://duckdns.org/ 

#### DiffeHellman key creation

Run this once and save the file and put it in the folder nginx-certbot-proxy
```bash
openssl dhparam -out ./dhparams.pem 2048
```

#### Using certbot to obtain certificates for your dns

```bash
certbot certonly --standalone -d $DNS_HOST --email $YOUR_EMAIL --agree-tos --quiet
```
#### Building a small nginx alpine image with certbot

The folder nginx-certbot-proxy contains the Dockerfile and configuration required to make SSL work.

*OBS!: BEFORE YOU RUN THE BUILD PLEASE UPDATE THE CONFIGURATION FILE FOUND IN sslstuff/ngnix. 

YOU MUST UPDATE THE VARIABLE: server_name xxx.xxx.xxx; TO THE NAME YOU USED FOR THE DNS SETUP*

```bash
cd nginx-certbot-proxy
docker build -t sslstuff . --build-arg DNS_HOST=xxxx --build-arg YOUR_EMAIL=xxxx@xx.com
```

## Google Home with Kodi and IFTTT

Your webserver needs to be publicly accessible
Following the steps above to create a dns name for your webserver. Either nsupdate.info works or you can use duckdns.org


The system still needs the IFTT integration set up.

1. Sign in to IFTT
2. Create New Applet
3. Select GoogleAssitant
4. Select IF
5. Fill in key phrase and response
6. Select Webhook
7. Fill in url to your webserver responsible for processing requests: https://dns-name:port/path/rest/resource
8. Select Application/json
9. Fill in the token as part of the request body





