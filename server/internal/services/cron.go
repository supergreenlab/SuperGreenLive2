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
	"os"
	"strings"
	"sync"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/api"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/tools"
	"github.com/gofrs/uuid"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	c                     *cron.Cron
	timelapseEntryID      *cron.EntryID
	timelapseEntryIDMutex sync.Mutex
)

type timelapseUploadURLRequest struct{}

type timelapseUploadURLResult struct {
	UploadPath string `json:"uploadPath"`
}

func captureTimelapse() {
	resp := timelapseUploadURLResult{}
	if err := api.POSTSGLObject("/timelapseUploadURL", &timelapseUploadURLRequest{}, &resp); err != nil {
		logrus.Errorf("api.POSTSGLObject in captureTimelapse %q", err)
		return
	}

	cam, err := tools.TakePic()
	if err != nil {
		logrus.Errorf("takePic in captureTimelapse %q", err)
		return
	}
	reader, err := os.Open(cam)
	if err != nil {
		logrus.Errorf("os.Open in captureTimelapse %q", err)
		return
	}
	defer reader.Close()

	fi, err := reader.Stat()
	if err != nil {
		logrus.Errorf("reader.Stat in captureTimelapse %q", err)
		return
	}

	err = api.UploadSGLObject(resp.UploadPath, reader, fi.Size())
	if err != nil {
		logrus.Errorf("api.UploadSGLObject in captureTimelapse %q", err)
		return
	}

	timelapseID, err := kv.GetString("timelapseid")
	if err != nil {
		logrus.Errorf("kv.GetString in captureTimelapse %q", err)
		return
	}

	timelapseIDUUID, err := uuid.FromString(timelapseID)
	if err != nil {
		logrus.Errorf("uuid.FromString in captureTimelapse %q", err)
		return
	}

	uploadPath := strings.Split(resp.UploadPath, "?")
	frame := appbackend.TimelapseFrame{
		TimelapseID: timelapseIDUUID,
		FilePath:    uploadPath[0],
		Meta:        "{}",
	}

	if err := api.POSTSGLObject("/timelapseframe", &frame, nil); err != nil {
		logrus.Errorf("api.POSTSGLObject in captureTimelapse %q", err)
		return
	}
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
