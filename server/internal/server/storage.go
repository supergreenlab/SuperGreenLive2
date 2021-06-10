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
	"archive/zip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func zipHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	storageDir := viper.Get("StorageDir").(string)
	files, err := os.ReadDir(storageDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/zip")
	zw := zip.NewWriter(w)

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", storageDir, file.Name())
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			logrus.Errorf("f.Close in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileInfos, err := f.Stat()
		if err != nil {
			logrus.Errorf("fs.Stat in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fileInfoHeaders, err := zip.FileInfoHeader(fileInfos)
		if err != nil {
			logrus.Errorf("zip.FileInfoHeader in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fz, err := zw.CreateHeader(fileInfoHeaders)
		if err != nil {
			logrus.Errorf("w.CreateHeader in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			logrus.Errorf("fz.Write in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = fz.Write(b)
		if err != nil {
			logrus.Errorf("fz.Write in zipHandler %q", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := zw.Close(); err != nil {
		logrus.Errorf("zw.Close in zipHandler %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
