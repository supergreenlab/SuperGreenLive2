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

type TokenData struct {
	Token string `json:"token"`
}

func tokenHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	td := TokenData{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&td); err != nil {
		logrus.Errorf("json.Unmarshal in tokenHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := kv.SetString("token", td.Token); err != nil {
		logrus.Errorf("kv.SetString in tokenHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK")
}

func loggedInHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if _, err := kv.GetString("token"); err == nil {
		fmt.Fprintf(w, "true")
		return
	}
	fmt.Fprintf(w, "false")
}
