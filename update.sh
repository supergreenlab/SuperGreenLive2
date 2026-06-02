#!/bin/bash

set -e

TAG=${1:-latest}

configure_swap() {
  local VAR="CONF_SWAPSIZE"
  local VALUE="1024"
  local FILE="/etc/dphys-swapfile"

  # Check total RAM in MB
  local TOTAL_RAM=$(free -m | awk '/^Mem:/{print $2}')

  if [ "$TOTAL_RAM" -gt 1024 ]; then
    echo "System has more than 1GB of RAM ($TOTAL_RAM MB). Skipping swap size modification."
    return
  fi

  # Ensure the file exists, otherwise create it
  if [ ! -f "$FILE" ]; then
    echo "$FILE not found. Creating file."
    sudo touch "$FILE"
  fi

  # Modify the CONF_SWAPSIZE value or append it if not present
  if grep -q "^#\?\s*$VAR=" "$FILE"; then
    sudo sed -i "s/^#\?\s*$VAR=.*/$VAR=$VALUE/" "$FILE"
  else
    echo "$VAR=$VALUE" | sudo tee -a "$FILE" > /dev/null
  fi

  echo "Swap size configured to $VALUE MB."
}

dphys-swapfile swapoff
configure_swap
dphys-swapfile setup
dphys-swapfile swapon

apt-get --allow-releaseinfo-change update

apt-get install --yes \
        fswebcam ffmpeg libmagickwand-7.q16-10 \
        python3-opencv \
        python3-libcamera python3-picamera2
apt --reinstall install --yes libcamera-apps-lite


curl --remote-name \
     --location \
     https://github.com/SuperGreenLab/SuperGreenLive2/releases/download/$TAG/liveserver.zip
rm -r liveserver
unzip liveserver.zip

systemctl stop liveserver
cp -r liveserver/assets/* /usr/local/share/appbackend
cp -r liveserver/static/* /usr/local/share/appbackend_static
if [ "$(dpkg --print-architecture)" = "arm64" ]; then
  cp liveserver/liveserver_arm64 /usr/local/bin/liveserver
elif [ "$(uname -m)" = "armv6l" ]; then
  cp liveserver/liveserver_arm32v6 /usr/local/bin/liveserver
else
  cp liveserver/liveserver_arm32 /usr/local/bin/liveserver
fi
cp liveserver/tools/* /usr/local/bin/
systemctl start liveserver
