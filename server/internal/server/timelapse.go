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

	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/kv"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type TimelapseData struct {
	ID              *string `json:"id,omitempty"`
	PlantID         *string `json:"plantID,omitempty"`
	Rotation        *string `json:"rotation,omitempty"`
	SkipNight       *string `json:"skipNight,omitempty"`
	StorageDuration *string `json:"storageDuration,omitempty"`
	RaspiParams     *string `json:"raspiParams,omitempty"`
	FSWebCamParams  *string `json:"fswebcamParams,omitempty"`
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

	if td.Rotation != nil {
		if err := kv.SetString("rotation", *td.Rotation); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.SkipNight != nil {
		if err := kv.SetString("skipnight", *td.SkipNight); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.StorageDuration != nil {
		if err := kv.SetString("storageduration", *td.StorageDuration); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.RaspiParams != nil {
		if err := kv.SetString("raspiparams", *td.RaspiParams); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if td.FSWebCamParams != nil {
		if err := kv.SetString("fswebcamparams", *td.FSWebCamParams); err != nil {
			logrus.Errorf("kv.SetString in timelapseHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "OK")
}

func getTimelapseHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	td := TimelapseData{
		ID:              kv.GetStringOrNil("timelapseid"),
		PlantID:         kv.GetStringOrNil("plantid"),
		Rotation:        kv.GetStringOrNil("rotation"),
		SkipNight:       kv.GetStringOrNil("skipnight"),
		StorageDuration: kv.GetStringOrNil("storageduration"),
		RaspiParams:     kv.GetStringOrNil("raspiparams"),
		FSWebCamParams:  kv.GetStringOrNil("fswebcamparams"),
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(td); err != nil {
		logrus.Errorf("encoder.Encode in getTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
