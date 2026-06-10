/*
 * Copyright (C) 2021  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg/api"
	sglimage "github.com/SuperGreenLab/AppBackend/pkg/image"
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/kv"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const captureCmdTimeout = 60 * time.Second

var (
	camMutex sync.Mutex

	StopStreamFunc func() error
)

func init() {
	viper.SetDefault("USBCam", false)
}

func findExecutable(names ...string) (string, error) {
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no executable found (tried %s)", strings.Join(names, ", "))
}

func stillCaptureExecutable() (string, error) {
	return findExecutable("rpicam-still", "libcamera-still")
}

func StreamExecutable() (string, error) {
	if USBCam() {
		return findExecutable("usbcam-streamer")
	}
	return findExecutable("libcamera-streamer")
}

// TODO move this to api
type DeviceParamsResult struct {
	Params map[string]interface{} `json:"params"`
}

func (dpr DeviceParamsResult) GetInt(device appbackend.Device, key string) (int, error) {
	k := fmt.Sprintf("%s.KV.%s", device.Identifier, key)
	v, ok := dpr.Params[k].(string)
	if !ok {
		return 0, errors.New("Not found")
	}
	return strconv.Atoi(v)
}

func GetLedBox(box appbackend.Box, device appbackend.Device) (appbackend.GetLedBox, error) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in GetLedBox %q", err)
		return nil, err
	}

	deviceParams := DeviceParamsResult{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/device/%s/params?params=LED_*_BOX", box.DeviceID.UUID.String()), &deviceParams); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(device/params) in GetLedBox %q", err)
		return nil, err
	}
	return func(i int) (int, error) {
		k := fmt.Sprintf("LED_%d_BOX", i)
		return deviceParams.GetInt(device, k)
	}, nil
}

var lastPic time.Time

func WaitCamAvailable() {
	camMutex.Lock()
	defer camMutex.Unlock()
}

func TakePic() (string, error) {
	if StopStreamFunc != nil {
		if err := StopStreamFunc(); err != nil {
			logrus.Warningf("StopStreamFunc in TakePic %q", err)
		}
	}

	camMutex.Lock()
	defer camMutex.Unlock()
	logrus.Info("Taking picture..")

	rotation, err := kv.GetString("rotation")
	if err != nil {
		rotation = "0"
	}

	tmpFile, err := os.CreateTemp("/tmp", "cam-*.jpg")
	if err != nil {
		logrus.Errorf("os.CreateTemp in TakePic %q", err)
		return "", err
	}
	name := tmpFile.Name()
	tmpFile.Close()

	var execPath string
	params := []string{}

	if USBCam() == false {
		execPath, err = stillCaptureExecutable()
		if err != nil {
			os.Remove(name)
			logrus.Errorf("stillCaptureExecutable in TakePic %q", err)
			return "", err
		}
		logrus.Infof("Using %s for still capture", execPath)
		raspiParams, err := kv.GetString("raspiparams")
		if err != nil {
			logrus.Errorf("kv.GetString(raspiparams) in TakePic %q", err)
		}

		params = strings.FieldsFunc(raspiParams, func(c rune) bool {
			return c == ' '
		})
		params = append(params, []string{"--rotation", rotation, "--quality", "100", "--output", name}...)
	} else {
		execPath, err = findExecutable("fswebcam")
		if err != nil {
			os.Remove(name)
			logrus.Errorf("findExecutable(fswebcam) in TakePic %q", err)
			return "", err
		}
		fswebcamParams, err := kv.GetString("fswebcamparams")
		if err != nil {
			logrus.Errorf("kv.GetString(fswebcamparams) in TakePic %q", err)
		}

		params = strings.FieldsFunc(fswebcamParams, func(c rune) bool {
			return c == ' '
		})
		params = append(params, []string{"-d", fmt.Sprintf("/dev/%s", viper.GetString("VideoDev")), "--rotate", rotation, "--resolution", "2592x1944", name}...)
	}

	ctx, cancel := context.WithTimeout(context.Background(), captureCmdTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, execPath, params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		os.Remove(name)
		logrus.Errorf("cmd.Run in TakePic %q", err)
		return "", err
	}
	return name, nil
}

func CaptureFrame() (*bytes.Buffer, error) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in CaptureFrame %q", err)
		return nil, err
	}

	plantID, err := kv.GetString("plantid")
	if err != nil {
		logrus.Errorf("kv.GetString(plant) in CaptureFrame %q", err)
		return nil, err
	}

	plant := appbackend.Plant{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/plant/%s", plantID), &plant); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(plant) in CaptureFrame %q", err)
		return nil, err
	}
	box := appbackend.Box{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/box/%s", plant.BoxID), &box); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(box) in CaptureFrame %q", err)
		return nil, err
	}
	var device *appbackend.Device = nil
	if box.DeviceID.Valid == true {
		device = &appbackend.Device{}
		if err := appbackend.GETSGLObject(token, fmt.Sprintf("/device/%s", box.DeviceID.UUID), device); err != nil {
			logrus.Errorf("appbackend.GETSGLObject(device) in CaptureFrame %q", err)
			return nil, err
		}
	}

	cam, err := TakePic()
	if err != nil {
		logrus.Errorf("takePic in CaptureFrame %q", err)
		return nil, err
	}
	defer os.Remove(cam)

	reader, err := os.Open(cam)
	if err != nil {
		logrus.Errorf("os.Open in CaptureFrame %q", err)
		return nil, err
	}
	defer reader.Close()

	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		logrus.Errorf("image.Decode in CaptureFrame %q", err)
		return nil, err
	}
	var resized image.Image
	if img.Bounds().Max.X > img.Bounds().Max.Y {
		resized = imaging.Resize(img, 1250, 0, imaging.Lanczos)
	} else {
		resized = imaging.Resize(img, 0, 1250, imaging.Lanczos)
	}

	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, resized, &jpeg.Options{Quality: 80})
	if err != nil {
		logrus.Errorf("jpeg.Encode in CaptureFrame %q", err)
		return nil, err
	}

	// TODO DRY this with timelapse service
	meta := appbackend.MetricsMeta{Date: time.Now()}
	if device != nil {
		getLedBox, err := GetLedBox(box, *device)
		if err != nil {
			logrus.Errorf("tools.GetLedBox in CaptureFrame %q", err)
			return nil, err
		}

		t := time.Now()
		from := t.Add(-24 * time.Hour)
		to := t
		meta = appbackend.LoadMetricsMeta(*device, box, from, to, appbackend.LoadGraphValue, getLedBox)
	}

	buff, err = sglimage.AddSGLOverlays(box, plant, meta, buff)
	if err != nil {
		logrus.Errorf("addSGLOverlays in CaptureFrame %q - device: %+v", err, device)
		return nil, err
	}

	return buff, nil
}

func USBCam() bool {
	usbCam := viper.GetBool("USBCam")
	return usbCam
}
