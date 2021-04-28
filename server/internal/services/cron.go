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

package services

import (
	"sync"

	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/api"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/tools"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	c                     *cron.Cron
	timelapseEntryID      *cron.EntryID
	timelapseEntryIDMutex sync.Mutex
)

func captureTimelapse() {
	if _, err := tools.CaptureFrame(); err != nil {
		logrus.Errorf("tools.CaptureFrame in captureTimelapse %q", err)
		return
	}

	api.LoadSGLObject("/timelapseUploadURL")
}

func ScheduleTimelapse() {
	if cron, err := kv.GetString("cron"); err != nil {
		logrus.Errorf("kv.GetString in ScheduleTimelapse %q", err)
		return
	} else {
		timelapseEntryIDMutex.Lock()
		if timelapseEntryID != nil {
			c.Remove(*timelapseEntryID)
		}
		if entryID, err := c.AddFunc(cron, captureTimelapse); err != nil {
			logrus.Errorf("c.AddFunc in ScheduleTimelapse %q", err)
		} else {
			timelapseEntryID = &entryID
		}
		timelapseEntryIDMutex.Unlock()
	}
}

func Init() {
	c = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	ScheduleTimelapse()
	c.Start()
}
