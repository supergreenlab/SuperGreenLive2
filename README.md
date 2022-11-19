![SuperGreenLab](docs/sgl.png?raw=true)

Table of Contents
=================

* [SuperGreenLivePI2](#supergreenlivepi2)
   * [Why timelapses?](#why-timelapses)
   * [Features](#features)
      * [Linked to your SGL account](#linked-to-your-sgl-account)
      * [Select you plant diary](#select-you-plant-diary)
      * [Easy installation](#easy-installation)
      * [Live status view](#live-status-view)
* [Hardware requirements](#hardware-requirements)
* [Installation](#installation)
  * [Debian bulsseye ](#debian-bulsseye)
  * [Debian buster ](#debian-buster)
  * [Install the liveserver](#install-the-liveserver)
    * [USB cameras](#usb-cameras)


There is now a complete and detailed guide on the website, please follow it [here](https://www.supergreenlab.com/guide/how-to-setup-a-remote-live-camera) :)
===

# SuperGreenLivePI2

Remote live cam companion that goes with the [SuperGreenApp2](https://github.com/supergreenlab/SuperGreenApp2)

Check your grow in real time. Posts daily and weekly timelapses to your diary feed.

![Example](docs/screenshot-live.png?raw=true)

## Why timelapses?

One thing we tend to ignore for obvious reasons, is the plant's movements.
Plants actually move a lot during the day, it's too slow for us to notice, but becomes very clear once in high speed.

So that's the point of this, take a pic every 10 minutes, then compile all those pics into videos, daily and weekly.

One of the thing with movements, is they can allow to spot something wrong before it shows up.

In normal conditions, the plant kind of "breathes", as seen in this [timelapse](https://www.instagram.com/p/BvMcC_oH94E/).

[This plant](https://www.instagram.com/p/BvZReZBHzrO/) on the other hand is thirsty, the leaves start to go down slowly, and the breathing has stopped, notice the [next day](https://www.instagram.com/p/Bvb2ULdn1_5/) how it was in really bad condition, and how it bounced back when fed water.

## Features

- Easy installation for raspberry pi
- Web interface for easier setup
- Takes pictures every 10 minutes
- Adds controller sensor metrics as meta data for later analysis
- Posts daily and weekly timelapses to your plant diary feed

### Linked to your SGL account

![Login screen](docs/screen-login.png?raw=true)

### Select you plant diary

![Plant screen](docs/screen-plant.png?raw=true)

### Easy installation

Low latency live feed for easier installation and focus tuning.

![Camera screen](docs/screen-camera.png?raw=true)

### Live status view

![Index screen](docs/screen-index.png?raw=true)

# Hardware requirements

- [RaspberryPI](https://www.raspberrypi.org/products/) + [Wifi (optional, most rpi have integrated wifi now)](https://www.raspberrypi.org/products/raspberry-pi-usb-wifi-dongle/)
- [Camera](https://www.raspberrypi.org/products/camera-module-v2/), I got [those](https://www.amazon.com/SainSmart-Fish-Eye-Camera-Raspberry-Arduino/dp/B00N1YJKFS) for the wide angle lens, but that's only for small spaces (this is the one used for the pic above).
- [Power supply](https://www.raspberrypi.org/products/raspberry-pi-universal-power-supply/)

# Installation

First follow the raspbian [official quickstart](https://projects.raspberrypi.org/en/projects/raspberry-pi-getting-started).
You'll need an interface connection setup with wifi or ethernet.
Open a terminal either through a screen+keyboard or a ssh session.

## Debian bulsseye 

In the current debian buster for raspian pi, the handling for the camera has [changed](https://www.raspberrypi.com/news/bullseye-camera-system/). 
It is currently still possible to use the deprecated camera system under debian bullseye, but the support for it will eventually be dropped.

Actually cameras work right out of the box after flashing debian bullseye.
To check if the camera basically works under debian bullseye with the new libcamera, type the following into a terminal:

```sh
libcamera-still -o /tmp/test.jpg
```

`libcamera-still` is the replacement for the deprecated `raspistill` and the above command takes a test image and saves it.

## Debian buster

First thing is to enable camera interface, this is done through `raspi-config`, type in the terminal:

```sh
sudo raspi-config
```

Then, with the arrow keys, go to `Interface Options` then `Camera`, enable it and then say `Yes` when it proposes to reboot.

Once the raspberrypi has reboot, open a terminal, and type:

## Install the liveserver

```sh
curl -sL https://raw.githubusercontent.com/black-161-flag/SuperGreenLive2/switch_to_libcamera/install.sh | sudo bash
```

### USB cameras

To use an usb camera, add the following lines to the liveserver config under `/etc/liveserver/liveserver.toml`:

```
USBCam=true
VideoDev="video0"
```
if you are not using the first camera under `/dev/video0`, replace with the appropriate video device (for exampe `video1` for `/dev/video1`).

And restart the liveserver service:

```
sudo systemctl restart liveserver.service
``` 

Once this is done, open the page at http://localhost:8081 if using a pi with screen+keyboard, or http://raspberrypi.local:8081 from another computer (to get a live view).

On windows you might need to install [the Bonjour protocol from Apple](https://support.apple.com/kb/DL999?locale=en_US) to be able to find by name (needs reboot).

You can also replace the raspberrypi.local part by the rpi's IP address if you can find it from your router's interface.

# Reset to default

You might want to reset the installation to default.
To do so, run this command:

```sh

sudo rm -rf /var/liveserver

```
