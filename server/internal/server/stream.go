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
	"text/template"
	"time"

	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/tools"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_    = pflag.String("videodev", "video0", "Video device")
	_    = pflag.Int("streamport", 8082, "Stream port")
	tmpl *template.Template
)

func init() {
	viper.SetDefault("VideoDev", "video0")
	viper.SetDefault("StreamPort", 8082)

}

var cmd *exec.Cmd

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url, err := url.Parse(fmt.Sprintf("http://localhost:%d", viper.GetInt("StreamPort")))
	if err != nil {
		logrus.Errorf("url.Parse in streamHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func startStreamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd != nil {
		fmt.Fprintf(w, "OK")
		return
	}

	tools.WaitCamAvailable()
	if tools.UseLegacy() {
		log.Info("Starting stream via picamera-streamer")
		cmd = exec.Command("/usr/local/bin/picamera-streamer", "--height", "720", "--width", "960", "--port", "8082")
	} else if tools.USBCam() {
		log.Info("Starting stream via usbcam-streamer")
		cmd = exec.Command("/usr/local/bin/usbcam-streamer", "--height", "720", "--width", "960", "--port", "8082", "--device", "/dev/video0")
	} else {
		log.Info("Starting stream via libcamera-streamer")
		cmd = exec.Command("/usr/local/bin/libcamera-streamer", "--height", "720", "--width", "960", "--port", "8082")
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
