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

package main

import (
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/config"
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/data/kv"
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/server"
	"github.com/SuperGreenLab/SuperGreenLive2/server/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	_ = pflag.Bool("Dev", false, "")
)

func init() {
	viper.SetDefault("Dev", false)
}

func main() {
	config.Init()
	kv.Init()
	services.InitCron()
	services.InitAutoUpdate()

	server.Start()

	logrus.Info("Liveserver started")

	select {}
}
