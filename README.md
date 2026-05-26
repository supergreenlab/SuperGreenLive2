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
  * [Install the liveserver](#install-the-liveserver)
    * [USB cameras](#usb-cameras)
* [Upgrade](#upgrade)


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

Supported boards: **Raspberry Pi 2 and newer** (including Pi Zero 2 W), running current **32-bit or 64-bit Raspberry Pi OS** (Bookworm or newer). Raspberry Pi 1 is not supported.

- [Raspberry Pi](https://www.raspberrypi.com/products/) — Pi 3, Pi 4, or Pi Zero 2 W recommended for camera use; Pi 2 works on 32-bit OS
- [Camera module](https://www.raspberrypi.com/products/camera-module-v2/) (or USB camera — see below)
- [Power supply](https://www.raspberrypi.com/products/type-c-power-supply/)

# Installation

Flash the latest [Raspberry Pi OS](https://www.raspberrypi.com/software/) image, enable the camera in `raspi-config` if needed, and verify capture works:

```sh
rpicam-still -o /tmp/test.jpg
```

(`libcamera-still -o /tmp/test.jpg` also works on older Bookworm images.)

Open a terminal (local or SSH), then:

## Install the liveserver

```sh
curl -sL https://github.com/SuperGreenLab/SuperGreenLive2/releases/download/latest/install.sh | sudo bash
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

# Upgrade

To upgrade the timelapse installation, run this command in a terminal:

```sh
curl -sL https://github.com/supergreenlab/SuperGreenLive2/releases/download/latest/update.sh | sudo bash
```

# Reset to default

You might want to reset the installation to default.
To do so, run this command:

```sh

sudo rm -rf /var/liveserver

```
