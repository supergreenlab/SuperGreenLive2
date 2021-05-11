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
	"log"
	"os"
	"strings"
	"sync"
	"time"

	appbackend "github.com/SuperGreenLab/AppBackend/pkg"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/data/kv"
	"github.com/SuperGreenLab/SuperGreenLivePI2/server/internal/tools"
	"github.com/disintegration/imaging"
	"github.com/gofrs/uuid"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	c                     *cron.Cron
	timelapseEntryID      *cron.EntryID
	timelapseEntryIDMutex sync.Mutex
	_                     = pflag.String("storagedir", "/tmp/storage", "location for the latest pics")
)

func init() {
	viper.SetDefault("StorageDir", "/tmp/storage")
}

type timelapseUploadURLRequest struct{}

type timelapseUploadURLResult struct {
	UploadPath string `json:"uploadPath"`
}

func captureTimelapse() {
	token, err := kv.GetString("token")
	if err != nil {
		logrus.Errorf("kv.GetString(token) in captureTimelapse %q", err)
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

	resp := timelapseUploadURLResult{}
	if err := appbackend.POSTSGLObject(token, "/timelapseUploadURL", &timelapseUploadURLRequest{}, &resp); err != nil {
		logrus.Errorf("appbackend.POSTSGLObject(timelapseUploadURL) in captureTimelapse %q", err)
		return
	}

	err = appbackend.UploadSGLObject(resp.UploadPath, bytes.NewReader(buff.Bytes()), int64(buff.Len()))
	if err != nil {
		logrus.Errorf("appbackend.UploadSGLObject in captureTimelapse %q", err)
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
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/plant/%s/", plantID), &plant); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(plant) in captureHandler %q", err)
		return
	}
	box := appbackend.Box{}
	if err := appbackend.GETSGLObject(token, fmt.Sprintf("/box/%s/", plant.BoxID), &box); err != nil {
		logrus.Errorf("appbackend.GETSGLObject(box) in captureHandler %q", err)
		return
	}
	meta := appbackend.MetricsMeta{Date: time.Now()}
	if box.DeviceID.Valid == true {
		device := appbackend.Device{}
		if err := appbackend.GETSGLObject(token, fmt.Sprintf("/device/%s/", box.DeviceID.UUID), &device); err != nil {
			logrus.Errorf("appbackend.GETSGLObject(device) in captureHandler %q", err)
			return
		}

		skipNight, err := kv.GetStringWithDefault("skipnight", "true")
		if err != nil {
			logrus.Errorf("kv.GetStringWithDefault(skipnight) in CaptureFrame %q", err)
			return
		}
		if skipNight == "true" {
			deviceParams := tools.DeviceParamsResult{}
			url := fmt.Sprintf("/device/%s/params?params=BOX_%d_*_HOUR&params=BOX_%d_*_MIN", box.DeviceID.UUID, *box.DeviceBox, *box.DeviceBox)
			if err := appbackend.GETSGLObject(token, url, &deviceParams); err != nil {
				logrus.Errorf("appbackend.GETSGLObject(device/params) in captureHandler %q", err)
				return
			}
			onHour, _ := deviceParams.GetInt(device, fmt.Sprintf("BOX_%d_ON_HOUR", *box.DeviceBox))
			onMin, _ := deviceParams.GetInt(device, fmt.Sprintf("BOX_%d_ON_MIN", *box.DeviceBox))
			offHour, _ := deviceParams.GetInt(device, fmt.Sprintf("BOX_%d_OFF_HOUR", *box.DeviceBox))
			offMin, _ := deviceParams.GetInt(device, fmt.Sprintf("BOX_%d_OFF_MIN", *box.DeviceBox))
			t := time.Now()
			on := time.Date(t.Year(), t.Month(), t.Day(), onHour, onMin, 0, 0, time.UTC)
			off := time.Date(t.Year(), t.Month(), t.Day(), offHour, offMin, 0, 0, time.UTC)
			isOnNow := (on.Before(off) && t.After(on) && t.Before(off)) ||
				(on.After(off) && (t.Before(off) || t.After(on)))
			if !isOnNow {
				logrus.Infof("Skipping night time")
				return
			}
		}

		getLedBox, err := tools.GetLedBox(box, device)
		if err != nil {
			logrus.Errorf("tools.GetLedBox in captureHandler %q", err)
			return
		}
		t := time.Now()
		from := t.Add(-24 * time.Hour)
		to := t
		meta = appbackend.LoadMetricsMeta(device, box, from, to, appbackend.LoadGraphValue, getLedBox)
	}

	var metaStr string
	if j, err := json.Marshal(meta); err != nil {
		logrus.Errorf("json.Marshal in captureHandler %q", err)
		return
	} else {
		metaStr = string(j)
	}

	uploadPath := strings.Split(resp.UploadPath, "/")
	uploadPath = strings.Split(uploadPath[2], "?")
	frame := appbackend.TimelapseFrame{
		TimelapseID: timelapseIDUUID,
		FilePath:    uploadPath[0],
		Meta:        metaStr,
	}

	if err := appbackend.POSTSGLObject(token, "/timelapseframe", &frame, nil); err != nil {
		logrus.Errorf("appbackend.POSTSGLObject(timelapseframe) in captureTimelapse %q", err)
		return
	}

	if err := storePic(buff, frame); err != nil {
		logrus.Errorf("storePic in captureTimelapse %q", err)
		return
	}
}

func storePic(buff *bytes.Buffer, frame appbackend.TimelapseFrame) error {
	storageDir := viper.Get("StorageDir")
	name := strings.Split(frame.FilePath, "/")[1]
	path := fmt.Sprintf("%s/%s", storageDir, name)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err := f.Write(buff.Bytes()); err != nil {
		return err
	}

	return removeOldFiles()
}

func removeOldFiles() error {
	storageDir := viper.Get("StorageDir").(string)
	storageDuration, err := kv.GetIntWithDefault("storageDuration", 86400)
	files, err := os.ReadDir(storageDir)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now().Add(time.Duration(-storageDuration) * time.Second)
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", storageDir, file.Name())
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		s, err := f.Stat()
		if err != nil {
			return err
		}
		if s.ModTime().Before(t) {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ScheduleTimelapse() {
	if cron, err := kv.GetString("cron"); err != nil {
		logrus.Errorf("kv.GetString in ScheduleTimelapse %q", err)
		return
	} else {
		timelapseEntryIDMutex.Lock()
		defer timelapseEntryIDMutex.Unlock()
		if timelapseEntryID != nil {
			c.Remove(*timelapseEntryID)
		}
		entryID, err := c.AddFunc(cron, captureTimelapse)
		if err != nil {
			logrus.Errorf("c.AddFunc in ScheduleTimelapse %q", err)
			return
		}
		timelapseEntryID = &entryID
	}
}

func Init() {
	if _, err := os.Stat(viper.GetString("StorageDir")); os.IsNotExist(err) {
		if err := os.Mkdir(viper.GetString("StorageDir"), 0755); err != nil {
			logrus.Fatalf("%q", err)
		}
	}
	c = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))
	ScheduleTimelapse()
	c.Start()
}
