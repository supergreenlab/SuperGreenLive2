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
	"net/http"

	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type PlantData struct {
	ID string `json:"id"`
}

func plantHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pl := PlantData{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&pl); err != nil {
		logrus.Errorf("decoder.Decode in tokenHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := kv.SetString("plant", pl.ID); err != nil {
		logrus.Errorf("kv.SetString in plantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")
}
