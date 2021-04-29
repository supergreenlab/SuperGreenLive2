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

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/api"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
)

// TODO temporary
var rotate = false

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

	buff, err = appbackend.AddSGLOverlays(box, plant, device, buff)
	if err != nil {
		logrus.Errorf("addSGLOverlays in captureHandler %q - device: %+v", err, device)
		return nil, err
	}

	return buff, nil
}
