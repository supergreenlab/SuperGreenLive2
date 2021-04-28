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

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("apiurl", "http://192.168.1.87:8080", "SGL Backend api url")
	_ = pflag.String("storageurl", "http://192.168.1.87:9000", "SGL Backend storage url")
	_ = pflag.String("storagehost", "minio:9000", "SGL Backend storage host name")
)

func init() {
	viper.SetDefault("ApiUrl", "http://192.168.1.87:8080")
	viper.SetDefault("StorageUrl", "http://192.168.1.87:9000")
	viper.SetDefault("StorageHost", "minio:9000")
}

func GETSGLObject(url string, obj interface{}) error {
	url = fmt.Sprintf("%s%s", viper.GetString("ApiUrl"), url)

	token, err := kv.GetString("token")
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authentication", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}

func POSTSGLObject(url string, obj interface{}, respObj interface{}) error {
	url = fmt.Sprintf("%s%s", viper.GetString("ApiUrl"), url)

	token, err := kv.GetString("token")
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	request.Header.Set("Authentication", fmt.Sprintf("Bearer %s", token))
	request.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if respObj != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(respObj); err != nil {
			return err
		}
	}
	return nil
}

func UploadSGLObject(url string, obj io.Reader, length int64) error {
	url = fmt.Sprintf("%s%s", viper.GetString("StorageUrl"), url)

	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("PUT", url, obj)
	if err != nil {
		return err
	}
	request.Host = viper.GetString("StorageHost")
	request.ContentLength = length

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("%q", err)
		return err
	}
	if len(content) != 0 {
		err := errors.New(string(content))
		logrus.Errorf("Upload error: %q", err)
		return err
	}

	return nil
}
