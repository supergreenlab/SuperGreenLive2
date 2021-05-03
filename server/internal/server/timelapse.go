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
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/services"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type TimelapseData struct {
	ID      *string `json:"id,omitempty"`
	PlantID *string `json:"plantID,omitempty"`
	Cron    *string `json:"cron,omitempty"`
	Rotate  *string `json:"rotate,omitempty"`
}

func timelapseHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	td := TimelapseData{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&td); err != nil {
		logrus.Errorf("json.Unmarshal in timelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if td.ID != nil {
		if err := kv.SetString("timelapseid", *td.ID); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.PlantID != nil {
		if err := kv.SetString("plantid", *td.PlantID); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.Cron != nil {
		if err := kv.SetString("cron", *td.Cron); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.Rotate != nil {
		if err := kv.SetString("rotate", *td.Rotate); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	services.ScheduleTimelapse()

	fmt.Fprintf(w, "OK")
}
