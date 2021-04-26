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

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	db "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// TODO temporary
var rotate = false

func takePic() (string, error) {
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

func loadSGLObject(url string, obj interface{}) error {
	token, err := kv.GetString("token")
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authentication", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func captureHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Info("Taking picture..")

	plantID, err := kv.GetString("plant")
	if err != nil {
		logrus.Errorf("kv.GetString(plant) in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plant := db.Plant{}
	if err := loadSGLObject(fmt.Sprintf("https://api2.supergreenlab.com/plant/%s/", plantID), &plant); err != nil {
		logrus.Errorf("loadSGLObject(plant) in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	box := db.Box{}
	if err := loadSGLObject(fmt.Sprintf("https://api2.supergreenlab.com/box/%s/", plant.BoxID), &box); err != nil {
		logrus.Errorf("loadSGLObject(box) in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var device *db.Device = nil
	if box.DeviceID.Valid == true {
		device = &db.Device{}
		if err := loadSGLObject(fmt.Sprintf("https://api2.supergreenlab.com/device/%s/", box.DeviceID.UUID), device); err != nil {
			logrus.Errorf("loadSGLObject(device) in captureHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	cam, err := takePic()
	if err != nil {
		logrus.Errorf("takePic in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, err := os.Open(cam)
	if err != nil {
		logrus.Errorf("os.Open in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	buffBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.Errorf("ioutil.ReadAll in captureHandler %q - device: %+v", err, device)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buff := bytes.NewBuffer(buffBytes)
	buff, err = db.AddSGLOverlays(box, plant, device, buff)
	if err != nil {
		logrus.Errorf("addSGLOverlays in captureHandler %q - device: %+v", err, device)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buff.Bytes())
}
