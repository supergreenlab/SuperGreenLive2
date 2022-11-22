#!/bin/bash

set -e

apt-get --allow-releaseinfo-change update

apt-get install --yes \
        fswebcam ffmpeg libmagickwand-dev \
        python3-opencv python3-picamera

if [ "$(/usr/bin/lsb_release -rs)" -ge "11" ]; then
  echo "running on debian bullseye or newer"
  apt-get install --yes \
          python3-libcamera python3-picamera2
  curl --location \
       --output /usr/local/bin/libcamera-streamer \
       https://raw.githubusercontent.com/black-161-flag/libcamera-streamer/main/bin/libcamera-streamer
  chmod +x /usr/local/bin/libcamera-streamer
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

curl --location \
     --output /usr/local/bin/picamera-streamer \
     https://raw.githubusercontent.com/black-161-flag/picamera-streamer/main/bin/picamera-streamer
chmod +x /usr/local/bin/picamera-streamer

curl --location \
     --output /usr/local/bin/usbcam-streamer \
     https://raw.githubusercontent.com/black-161-flag/usbcam-streamer/main/bin/usbcam-streamer
chmod +x /usr/local/bin/usbcam-streamer

# curl -OL https://github.com/supergreenlab/SuperGreenLive2/releases/download/latest/liveserver.zip
curl --remote-name \
     --location \
     https://github.com/black-161-flag/SuperGreenLive2/releases/latest/liveserver.zip
unzip -o liveserver.zip

mkdir --parents /usr/local/share/appbackend /usr/local/share/appbackend_static

cp --recursive liveserver/assets/* /usr/local/share/appbackend
cp --recursive liveserver/static/* /usr/local/share/appbackend_static
cp liveserver/liveserver /usr/local/bin/liveserver

#if [ "$(dpkg --print-architecture)" = "arm64" ]; then
#  cp liveserver/liveserver_arm64 /usr/local/bin/liveserver
#else
#  cp liveserver/liveserver_arm32 /usr/local/bin/liveserver
#fi

mkdir --parents /etc/liveserver
cp liveserver/liveserver.toml /etc/liveserver/liveserver.toml

cp liveserver/liveserver.service /etc/systemd/system/
systemctl enable liveserver
systemctl start liveserver
