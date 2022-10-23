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

//go:embed motion.conf
var motionConf string

var (
	_    = pflag.String("videodev", "video0", "Video device")
	_    = pflag.Int("motionport", 8082, "Motion port")
	tmpl *template.Template
)

func init() {
	viper.SetDefault("VideoDev", "video0")
	viper.SetDefault("MotionPort", 8082)

	var err error
	tmpl, err = template.New("motionConf").Parse(motionConf)
	if err != nil {
		logrus.Fatal(err)
	}
}

var cmd *exec.Cmd

func motionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	url, err := url.Parse(fmt.Sprintf("http://localhost:%d", viper.GetInt("MotionPort")))
	if err != nil {
		logrus.Errorf("url.Parse in motionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(w, r)
}

func startMotionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd != nil {
		fmt.Fprintf(w, "OK")
		return
	}

	tools.WaitCamAvailable()
	if tools.UseLegacy() || tools.USBCam() {

		motionConfigPath := fmt.Sprintf("/tmp/motion-%d.conf", os.Getpid())
		mcp, err := os.OpenFile(motionConfigPath, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			logrus.Errorf("os.OpenFile in startMotionHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(mcp, struct {
			VideoDev   string
			MotionPort int
		}{
			viper.GetString("VideoDev"), viper.GetInt("MotionPort"),
		}); err != nil {
			mcp.Close()
			logrus.Errorf("tmpl.Execute in startMotionHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mcp.Close()

		cmd = exec.Command("/usr/bin/motion", "-c", motionConfigPath)
	} else {
		cmd = exec.Command("/usr/local/bin/libcamera-streamer", "--height", "720", "--width", "960", "--port", "8082")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logrus.Errorf("cmd.Start in startMotionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("Motion started")
	fmt.Fprintf(w, "OK")
	go func() {
		time.Sleep(5 * time.Minute)
		if err := stopMotion(); err != nil {
			log.Errorf("stopMotion in startMotionHandle timeout gorouting %q", err)
		}
	}()
}

func stopMotion() error {
	if cmd == nil {
		return nil
	}
	if err := cmd.Process.Kill(); err != nil {
		return err
	}
	logrus.Info("Motion stopped")
	cmd = nil
	return nil
}

func stopMotionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := stopMotion(); err != nil {
		log.Errorf("stopMotion in stopMotionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "OK")
}
