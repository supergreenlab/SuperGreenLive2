/*
 * Copyright (C) 2022  SuperGreenLab <towelie@supergreenlab.com>
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

package services

import (
	"fmt"
	"time"

	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/kv"
	wifi "github.com/mark2b/wpa-connect"
	"github.com/sirupsen/logrus"
)

func InitWifi() {
	ssid, err := kv.GetString("ssid")
	if err != nil {
		logrus.Errorf("kv.GetString(ssid) in CaptureFrame %q", err)
	}
	password, err := kv.GetString("wpassword")
	if err != nil {
		logrus.Errorf("kv.GetString(password) in CaptureFrame %q", err)
	}

	wifi.SetDebugMode()
	if ssid != "" && password != "" {
		if conn, err := wifi.ConnectManager.Connect(ssid, password, time.Second*60); err == nil {
			fmt.Println("Connected", conn.NetInterface, conn.SSID, conn.IP4.String(), conn.IP6.String())
		} else {
			fmt.Println(err)
		}
	}
}
