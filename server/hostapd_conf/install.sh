#!/bin/bash

sudo apt-get update && sudo apt-get upgrade -y

sudo apt-get install -y dnsmasq hostapd

sudo systemctl stop dnsmasq hostapd

for i in `find etc/* -type f`; do
  sudo mkdir -p /$(dirname $i)
  sudo cp $i /$i
done

sudo rfkill unblock wifi
sudo rfkill unblock all

sudo systemctl unmask hostapd
sudo systemctl start dnsmasq hostapd
# sudo reboot
