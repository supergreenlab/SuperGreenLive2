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
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	_ = pflag.String("motionurl", "http://localhost:8082", "Motion url")
)

func init() {
}

var cmd *exec.Cmd

func motionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	resp, err := http.Get("http://localhost:8082")
	if err != nil {
		logrus.Errorf("http.Get in motionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func startMotionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd != nil {
		fmt.Fprintf(w, "OK")
		return
	}
	cmd = exec.Command("/usr/bin/motion")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		logrus.Errorf("cmd.Start in startMotionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("Motion started")
	fmt.Fprintf(w, "OK")
}

func stopMotionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd == nil {
		fmt.Fprintf(w, "OK")
		return
	}
	if err := cmd.Process.Kill(); err != nil {
		log.Errorf("cmd.Process.Kill in stopMotionHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logrus.Info("Motion stopped")
	fmt.Fprintf(w, "OK")
	cmd = nil
}
