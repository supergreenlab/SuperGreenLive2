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
	_ "embed"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/tools"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("VideoDev", "video0")
	viper.SetDefault("StreamPort", 18082)
	viper.SetDefault("StreamHeight", 720)
	viper.SetDefault("StreamWidth", 960)
}

var cmd *exec.Cmd

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	logrus.Debug(fmt.Sprintf("start proxy for port %d", viper.GetInt("StreamPort")))
	proxyUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d", viper.GetInt("StreamPort")))
	if err != nil {
		logrus.Errorf("proxyUrl.Parse in streamHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.ServeHTTP(w, r)
}

func startStreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd != nil {
		fmt.Fprintf(w, "OK")
		return
	}

	tools.WaitCamAvailable()

	if tools.USBCam() {
		log.Debug("Starting stream via usbcam-streamer")
		cmd = exec.Command("/usr/local/bin/usbcam-streamer", "--height", viper.GetString("StreamHeight"), "--width", viper.GetString("StreamWidth"), "--port", viper.GetString("StreamPort"), "--device", fmt.Sprintf("/dev/%s", viper.GetString("VideoDev")))
	} else if tools.UseLegacy() {
		log.Debug("Starting stream via picamera-streamer")
		cmd = exec.Command("/usr/local/bin/picamera-streamer", "--height", viper.GetString("StreamHeight"), "--width", viper.GetString("StreamWidth"), "--port", viper.GetString("StreamPort"))
	} else {
		log.Debug("Starting stream via libcamera-streamer")
		cmd = exec.Command("/usr/local/bin/libcamera-streamer", "--height", viper.GetString("StreamHeight"), "--width", viper.GetString("StreamWidth"), "--port", viper.GetString("StreamPort"))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logrus.Errorf("cmd.Start in startStreamHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("Stream started")
	fmt.Fprintf(w, "OK")
	go func() {
		time.Sleep(5 * time.Minute)
		if err := stopStream(); err != nil {
			log.Errorf("stopStream in startStreamHandle timeout gorouting %q", err)
		}
	}()
}

func stopStream() error {
	if cmd == nil {
		return nil
	}
	logrus.Infof("%+v", cmd)
	if err := cmd.Process.Kill(); err != nil {
		return err
	}
	logrus.Info("Stream stopped")
	cmd = nil
	return nil
}

func stopStreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := stopStream(); err != nil {
		log.Errorf("stopStream in stopStreamHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "OK")
}
