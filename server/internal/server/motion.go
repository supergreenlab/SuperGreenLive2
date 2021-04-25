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
	"net/http"
	"os"
	"os/exec"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var cmd *exec.Cmd

func startMotionHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if cmd != nil {
		fmt.Fprintf(w, "OK")
		return
	}
	cmd = exec.Command("/usr/bin/motion")
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		logrus.Fatal(err)
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
		log.Fatal("failed to kill process: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	logrus.Info("Motion stopped")
	fmt.Fprintf(w, "OK")
	cmd = nil
}
