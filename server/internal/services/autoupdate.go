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
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var CommitDate string

var client = github.NewClient(nil)

func upgradeLiveServer() {
	logrus.Infof("Running upgrade script")
	tagName := "latest"
	if viper.GetBool("Dev") {
		tagName = "beta"
	}
	cmd := exec.Command("/usr/bin/bash", "-c", fmt.Sprintf("cd /tmp && curl -sL https://github.com/supergreenlab/SuperGreenLive2/releases/download/%s/update.sh | sudo bash", tagName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func checkVersion() {
	time.Sleep(1)
	commitDateInt, err := strconv.Atoi(CommitDate)
	if commitDateInt == 0 {
		logrus.Infof("Skipping upgrade, commitDate is 0")
		return
	}
	rr, _, err := client.Repositories.ListReleases(context.Background(), "SuperGreenLab", "SuperGreenLive2", nil)
	if err != nil {
		logrus.Errorf("client.Repositories.GetLatestRelease in checkVersion %q", err)
		return
	}
	for _, r := range rr {
		if (viper.GetBool("Dev") && *r.TagName == "beta") || (!viper.GetBool("Dev") && *r.TagName == "latest") {
			for _, a := range r.Assets {
				if *a.Name == "timestamp" {
					resp, err := http.Get(*a.BrowserDownloadURL)
					if err != nil {
						logrus.Errorf("client.Get in checkVersion %q", err)
						return
					}
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						logrus.Errorf("ioutil.ReadAll in checkVersion %q", err)
						return
					}
					timestamp, err := strconv.Atoi(strings.Trim(string(body), "\n"))
					if err != nil {
						logrus.Errorf("strconv.Atoi in checkVersion %q", err)
						return
					}
					if commitDateInt < timestamp {
						logrus.Infof("Trigger upgrade %d %d", commitDateInt, timestamp)
						upgradeLiveServer()
					}
				}
			}
		}
	}
}

func InitAutoUpdate() {
	go checkVersion()
}
