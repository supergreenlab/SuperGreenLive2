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

package kv

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	_ = pflag.String("leveldbdir", "/tmp/sgllive.leveldb", "LevelDB directory location")
)

func init() {
	viper.SetDefault("LevelDBDir", "/tmp/sgllive.leveldb")
}

var db *leveldb.DB

func GetString(key string) (string, error) {
	data, err := db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SetString(key, value string) error {
	err := db.Put([]byte(key), []byte(value), nil)
	return err
}

func Init() {
	var err error
	db, err = leveldb.OpenFile(viper.GetString("LevelDBDir"), nil)
	if err != nil {
		logrus.Fatalf("leveldb.OpenFile in kv.Init %q", err)
	}
}
