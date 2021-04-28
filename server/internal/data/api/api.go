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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.String("apiurl", "http://192.168.1.87:8080", "SGL Backend api url")
)

func init() {
	viper.SetDefault("ApiUrl", "http://192.168.1.87:8080")
}

func GETSGLObject(url string, obj interface{}) error {
	url = fmt.Sprintf("%s/%s", viper.GetString("ApiUrl"), url)

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
	url = fmt.Sprintf("%s/%s", viper.GetString("ApiUrl"), url)

	token, err := kv.GetString("token")
	if err != nil {
		return err
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	request, err := http.NewRequest("POST", url, nil)
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
