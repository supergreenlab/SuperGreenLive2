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

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/kv"
	"github.com/gofrs/uuid"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func getAPIPlantHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getPlantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pid := p.ByName("id")

	plant := appbackend.Plant{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/plant/%s", pid), &plant); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/plant/:id) in getPlantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(plant); err != nil {
		logrus.Errorf("encoder.Encode in getPlantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAPIBoxHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getBoxHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bid := p.ByName("id")

	box := appbackend.Box{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/box/%s", bid), &box); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/box/:id) in getBoxHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(box); err != nil {
		logrus.Errorf("encoder.Encode in getBoxHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAPITimelapseHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tid := p.ByName("id")

	timelapse := appbackend.Timelapse{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/timelapse/%s", tid), &timelapse); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/timelapse/:id) in getTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(timelapse); err != nil {
		logrus.Errorf("encoder.Encode in getTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetPlantsParams struct {
	Offset int
	Limit  int
}

type GetPlantsResult struct {
	Plants []appbackend.Plant `json:"plants"`
}

func getAPIPlantsHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getPlantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gpp := GetPlantsParams{}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&gpp, r.URL.Query()); err != nil {
		logrus.Errorf("DecodeQuery %q for %s", err.Error(), r.URL.Query())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	plants := GetPlantsResult{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/plants?offset=%d&limit=%d", gpp.Offset, gpp.Limit), &plants); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/plants) in getPlantHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(plants); err != nil {
		logrus.Errorf("encoder.Encode in getPlantsHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetBoxesParams struct {
	Offset int
	Limit  int
}

type GetBoxesResult struct {
	Boxes []appbackend.Box `json:"boxes"`
}

func getAPIBoxesHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getBoxesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gbp := GetBoxesParams{}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&gbp, r.URL.Query()); err != nil {
		logrus.Errorf("DecodeQuery %q for %s", err.Error(), r.URL.Query())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	boxes := GetBoxesResult{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/boxes?offset=%d&limit=%d", gbp.Offset, gbp.Limit), &boxes); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/boxes) in getBoxesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(boxes); err != nil {
		logrus.Errorf("encoder.Encode in getBoxesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetTimelapsesParams struct {
	Offset int
	Limit  int
}

type SelectTimelapsesResult struct {
	appbackend.Timelapse

	NFrames int `json:"nFrames"`
}

type GetTimelapsesResult struct {
	Timelapses []SelectTimelapsesResult `json:"timelapses"`
}

func getAPITimelapsesHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in getTimelapsesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	gpp := GetTimelapsesParams{}

	var decoder = schema.NewDecoder()
	if err := decoder.Decode(&gpp, r.URL.Query()); err != nil {
		logrus.Errorf("DecodeQuery %q for %s", err.Error(), r.URL.Query())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timelapses := GetTimelapsesResult{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/timelapses?addNFrames=true&offset=%d&limit=%d", gpp.Offset, gpp.Limit), &timelapses); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(/timelapses) in getTimelapsesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(timelapses); err != nil {
		logrus.Errorf("encoder.Encode in getTimelapsesHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type CreateTimelapseResponse struct {
	ID string `json:"id"`
}

func createAPITimelapseHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in createTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := appbackend.Timelapse{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&t); err != nil {
		logrus.Errorf("dec.Decode in createTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctr := CreateTimelapseResponse{}
	if err := appbackend.POSTSGLObject(token, "/timelapse", &t, &ctr); err != nil {
		logrus.Errorf("appbackend.POSTSGLObject(/timelapse) in createTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ctr); err != nil {
		logrus.Errorf("encoder.Encode in createTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UpdateTimelapseResponse struct {
	Status string `json:"status"`
}

func updateAPITimelapseHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	timelapseID, err := kv.GetString("timelapseid")
	if err != nil {
		logrus.Errorf("kv.GetString(timelapseid) in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := appbackend.Timelapse{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&t); err != nil {
		logrus.Errorf("dec.Decode in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := uuid.FromString(timelapseID)
	if err != nil {
		logrus.Errorf("uuid.FromString in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.ID.UUID = u
	t.ID.Valid = true

	utr := UpdateTimelapseResponse{}
	if err := appbackend.PUTSGLObject(token, "/timelapse", &t, &utr); err != nil {
		logrus.Errorf("appbackend.PUTSGLObject(/timelapse) in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(utr); err != nil {
		logrus.Errorf("encoder.Encode in updateTimelapseHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
