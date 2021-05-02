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
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"
	"sync"
	"time"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/api"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/tools"
	"github.com/disintegration/imaging"
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
		logrus.Errorf("api.POSTSGLObject(timelapseUploadURL) in captureTimelapse %q", err)
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

	img, err := imaging.Decode(reader, imaging.AutoOrientation(true))
	if err != nil {
		logrus.Errorf("image.Decode in captureHandler %q", err)
		return
	}

	var resized image.Image
	if img.Bounds().Max.X > img.Bounds().Max.Y {
		resized = imaging.Resize(img, 1250, 0, imaging.Lanczos)
	} else {
		resized = imaging.Resize(img, 0, 1250, imaging.Lanczos)
	}

	buff := new(bytes.Buffer)
	err = jpeg.Encode(buff, resized, &jpeg.Options{Quality: 80})
	if err != nil {
		logrus.Errorf("jpeg.Encode in captureHandler %q", err)
		return
	}

	err = api.UploadSGLObject(resp.UploadPath, bytes.NewReader(buff.Bytes()), int64(buff.Len()))
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

	plantID, err := kv.GetString("plantid")
	if err != nil {
		logrus.Errorf("kv.GetString(plant) in captureHandler %q", err)
		return
	}

	plant := appbackend.Plant{}
	if err := api.GETSGLObject(fmt.Sprintf("/plant/%s/", plantID), &plant); err != nil {
		logrus.Errorf("api.GETSGLObject(plant) in captureHandler %q", err)
		return
	}
	box := appbackend.Box{}
	if err := api.GETSGLObject(fmt.Sprintf("/box/%s/", plant.BoxID), &box); err != nil {
		logrus.Errorf("api.GETSGLObject(box) in captureHandler %q", err)
		return
	}
	metaStr := "{}"
	if box.DeviceID.Valid == true {
		device := appbackend.Device{}
		if err := api.GETSGLObject(fmt.Sprintf("/device/%s/", box.DeviceID.UUID), &device); err != nil {
			logrus.Errorf("api.GETSGLObject(device) in captureHandler %q", err)
			return
		}

		getLedBox, err := tools.GetLedBox(box, device)
		if err != nil {
			logrus.Errorf("tools.GetLedBox in captureHandler %q", err)
			return
		}
		t := time.Now()
		from := t.Add(-24 * time.Hour)
		to := t
		meta := appbackend.LoadMetricsMeta(device, box, from, to, appbackend.LoadGraphValue, getLedBox)
		if j, err := json.Marshal(meta); err != nil {
			logrus.Errorf("json.Marshal in captureHandler %q", err)
			metaStr = "{}"
		} else {
			metaStr = string(j)
		}
	}

	uploadPath := strings.Split(resp.UploadPath, "/")
	uploadPath = strings.Split(uploadPath[2], "?")
	frame := appbackend.TimelapseFrame{
		TimelapseID: timelapseIDUUID,
		FilePath:    uploadPath[0],
		Meta:        metaStr,
	}

	if err := api.POSTSGLObject("/timelapseframe", &frame, nil); err != nil {
		logrus.Errorf("api.POSTSGLObject(timelapseframe) in captureTimelapse %q", err)
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
