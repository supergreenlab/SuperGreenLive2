#!/bin/bash

set -e

TAG=${1:-latest}

apt-get install --yes \
        fswebcam ffmpeg libmagickwand-dev \
        python3-opencv python3-picamera

if [ "$(/usr/bin/lsb_release -rs)" -ge "11" ]; then
  echo "running on debian bullseye or newer"
  apt-get install --yes \
          python3-libcamera python3-picamera2
  apt --reinstall install --yes libcamera-apps-lite
else
  echo "running on debian buster or older"
  apt-get install --yes \
          python3-pip libatlas-base-dev
  pip3 install simplejpeg
  pip3 install numpy \
       --upgrade \
       --index-url https://www.piwheels.org/simple
fi


curl --remote-name \
     --location \
     https://github.com/SuperGreenLab/SuperGreenLive2/releases/download/$TAG/liveserver.zip
rm -r liveserver
unzip liveserver.zip

systemctl stop liveserver
cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver
cp liveserver/tools/* /usr/local/bin/
systemctl start liveserver
