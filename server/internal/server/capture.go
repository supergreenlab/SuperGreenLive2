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
	"net/http"

	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/tools"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func captureHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	buff := new(bytes.Buffer)

	if err := tools.CaptureFrame(buff); err != nil {
		logrus.Errorf("tools.CaptureFrame in captureHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(buff.Bytes())
}
