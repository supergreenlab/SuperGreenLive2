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
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"time"

	"image/jpeg"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
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

func captureHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Info("Taking picture..")
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
	img, _, err := image.Decode(reader)
	if err != nil {
		logrus.Errorf("image.Decode in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info(token)

	plantID, err := kv.GetString("plant")
	if err != nil {
		logrus.Errorf("kv.GetString(plant) in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info(plantID)

	plant := appbackend.Plant{}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api2.supergreenlab.com/plant/%s/", plantID), nil)
	request.Header.Set("Authentication", fmt.Sprintf("Bearer %s", token))

	if err != nil {
		logrus.Errorf("http.NewRequest in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		logrus.Errorf("client.Do in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&plant); err != nil {
		logrus.Errorf("decoder.Decode in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Infof("%+v", plant)

	jpeg.Encode(w, img, nil)
}
