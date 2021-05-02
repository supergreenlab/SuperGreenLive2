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
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/api"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

// TODO temporary
var rotate = false

type DeviceParamsResult struct {
	Params map[string]interface{}
}

func GetLedBox(box appbackend.Box, device appbackend.Device) (appbackend.GetLedBox, error) {
	deviceParams := DeviceParamsResult{}
	keys := []string{}
	for i := 0; i < 6; i += 1 {
		keys = append(keys, fmt.Sprintf("params=LED_%d_BOX"))
	}
	if err := api.GETSGLObject(fmt.Sprintf("/device/%s/params?", box.DeviceID.UUID, strings.Join(keys, "&")), &deviceParams); err != nil {
		logrus.Errorf("api.GETSGLObject(device/params) in captureHandler %q", err)
		return nil, err
	}
	return func(i int) (int, error) {
		v := deviceParams.Params[fmt.Sprintf("%s.KV.LED_%d_BOX", box.DeviceID.UUID, i)].(string)
		return strconv.Atoi(v)
	}, nil
}

func TakePic() (string, error) {
	logrus.Info("Taking picture..")

	var cmd *exec.Cmd
	name := "/tmp/cam.jpg"
	if rotate {
		cmd = exec.Command("/usr/bin/raspistill", "-vf", "-hf", "-q", "50", "-o", name)
	} else {
		cmd = exec.Command("/usr/bin/raspistill", "-q", "50", "-o", name)
	}
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return name, err
}

func CaptureFrame() (*bytes.Buffer, error) {
	plantID, err := kv.GetString("plantid")
	if err != nil {
		logrus.Errorf("kv.GetString(plant) in captureHandler %q", err)
		return nil, err
	}

	plant := appbackend.Plant{}
	if err := api.GETSGLObject(fmt.Sprintf("/plant/%s/", plantID), &plant); err != nil {
		logrus.Errorf("api.GETSGLObject(plant) in captureHandler %q", err)
		return nil, err
	}
	box := appbackend.Box{}
	if err := api.GETSGLObject(fmt.Sprintf("/box/%s/", plant.BoxID), &box); err != nil {
		logrus.Errorf("api.GETSGLObject(box) in captureHandler %q", err)
		return nil, err
	}
	var device *appbackend.Device = nil
	if box.DeviceID.Valid == true {
		device = &appbackend.Device{}
		if err := api.GETSGLObject(fmt.Sprintf("/device/%s/", box.DeviceID.UUID), device); err != nil {
			logrus.Errorf("api.GETSGLObject(device) in captureHandler %q", err)
			return nil, err
		}
	}

	cam, err := TakePic()
	if err != nil {
		logrus.Errorf("takePic in captureHandler %q", err)
		return nil, err
	}

	reader, err := os.Open(cam)
	if err != nil {
		logrus.Errorf("os.Open in captureHandler %q", err)
		return nil, err
	}
	defer reader.Close()

	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		logrus.Errorf("image.Decode in captureHandler %q", err)
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
		logrus.Errorf("jpeg.Encode in captureHandler %q", err)
		return nil, err
	}

	var meta *appbackend.MetricsMeta
	if device != nil {
		getLedBox, err := GetLedBox(box, *device)
		if err != nil {
			logrus.Errorf("tools.GetLedBox in captureHandler %q", err)
			return nil, err
		}

		t := time.Now()
		from := t.Add(-24 * time.Hour)
		to := t
		m := appbackend.LoadMetricsMeta(*device, box, from, to, appbackend.LoadGraphValue, getLedBox)
		meta = &m
	}

	buff, err = appbackend.AddSGLOverlays(box, plant, meta, buff)
	if err != nil {
		logrus.Errorf("addSGLOverlays in captureHandler %q - device: %+v", err, device)
		return nil, err
	}

	return buff, nil
}
