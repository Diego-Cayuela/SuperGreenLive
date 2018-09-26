/*
 * Copyright (C) 2018  SuperGreenLab <towelie@supergreenlab.com>
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
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/sirupsen/logrus"
)

var dbx files.Client

func init() {
	token := MustGetenv("DBX_TOKEN")
	config := dropbox.Config{
		Token: token,
	}

	dbx = files.New(config)
}

func fu(e error) {
	if e != nil {
		logrus.Fatal(e)
	}
}

func MustGetenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logrus.Fatalf("missing env %s", name)
	}
	return v
}

func takePic() (string, error) {
	name := "latest.jpg"
	cmd := exec.Command("/usr/bin/raspistill", "-vf", "-hf", "-q", "50", "-o", name)
	err := cmd.Run()
	return name, err
}

func uploadPic(name, local, remote string) {
	f, err := os.Open(local)
	fu(err)

	p := fmt.Sprintf("/%s/%s", name, remote)
	ci := files.NewCommitInfo(p)
	ci.Mode.Tag = "overwrite"
	_, err = dbx.Upload(ci, f)
	fu(err)

	logrus.Infof("Uploaded %s", p)
}

func main() {
	name := MustGetenv("NAME")

	remote := fmt.Sprintf("%d.jpg", int32(time.Now().Unix()))
	local, err := takePic()
	fu(err)

	logrus.Info("Uploading files")
	uploadPic(name, local, remote)
	uploadPic(name, local, "latest.jpg")
}
