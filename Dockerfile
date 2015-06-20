FROM ubuntu
MAINTAINER Anchal Agrawal
RUN apt-get update && apt-get install -y python-pip python-dev build-essential
RUN pip install vaurien
