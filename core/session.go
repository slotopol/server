package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LumenResearch/uasurfer"
	"github.com/gin-gonic/gin"
)

type Device struct {
	FormFactor  string `xorm:"'ff' varchar(8) notnull" json:"ff" yaml:"ff" xml:"ff,attr"`
	DeviceModel string `xorm:"'dm' varchar(24) notnull" json:"dm" yaml:"dm" xml:"dm,attr"`
	OSName      string `xorm:"'os' varchar(12) notnull" json:"os" yaml:"os" xml:"os,attr"`
	OSVersion   string `xorm:"'pv' varchar(4) notnull" json:"pv" yaml:"pv" xml:"pv,attr"`
	Browser     string `xorm:"'br' varchar(16) notnull" json:"br" yaml:"br" xml:"br,attr"`
}

func (d *Device) IsZero() bool {
	return d.FormFactor == "" && d.DeviceModel == "" && d.OSName == "" && d.OSVersion == "" && d.Browser == ""
}

func (d *Device) IsCongruent(d2 *Device) bool {
	if d.FormFactor != d2.FormFactor &&
		d.FormFactor != "" && d.FormFactor != "Unknown" &&
		d2.FormFactor != "" && d2.FormFactor != "Unknown" {
		return false
	}
	if d.DeviceModel != d2.DeviceModel &&
		d.DeviceModel != "" && d.DeviceModel != "Unknown" &&
		d2.DeviceModel != "" && d2.DeviceModel != "Unknown" {
		return false
	}
	if d.OSName != d2.OSName &&
		d.OSName != "" && d.OSName != "Unknown" &&
		d2.OSName != "" && d2.OSName != "Unknown" {
		return false
	}
	if d.OSVersion != d2.OSVersion &&
		d.OSVersion != "" && d.OSVersion != "Unknown" && d.OSVersion != "0" &&
		d2.OSVersion != "" && d2.OSVersion != "Unknown" && d2.OSVersion != "0" {
		return false
	}
	if d.Browser != d2.Browser &&
		d.Browser != "" && d.Browser != "Unknown" &&
		d2.Browser != "" && d2.Browser != "Unknown" {
		return false
	}
	return true
}

func (d *Device) Setup(c *gin.Context) {
	var ua uasurfer.UserAgent
	var rawua = c.GetHeader("User-Agent")
	uasurfer.ParseUserAgentWithHints(rawua, &uasurfer.Hints{
		Mobile:          c.GetHeader("Sec-CH-UA-Mobile"),
		Platform:        c.GetHeader("Sec-CH-UA-Platform"),
		PlatformVersion: c.GetHeader("Sec-CH-UA-Platform-Version"),
		FormFactors:     c.GetHeader("Sec-CH-UA-Form-Factors"),
	}, &ua)
	d.FormFactor = ua.DeviceType.StringTrimPrefix()
	d.DeviceModel = strings.Trim(c.GetHeader("Sec-CH-UA-Model"), `"`)
	d.OSName = ua.OS.Name.StringTrimPrefix()
	d.OSVersion = strconv.Itoa(ua.OS.Version.Major)
	d.Browser = ua.Browser.Name.StringTrimPrefix()
}

func (d *Device) String() string {
	return fmt.Sprintf("ff:%s|dm:%s|os:%s|pv:%s|br:%s",
		d.FormFactor, d.DeviceModel, d.OSName, d.OSVersion, d.Browser)
}

type Session struct {
	SID    uint64    `xorm:"pk autoincr" json:"sid" yaml:"sid" xml:"sid,attr"`
	UID    uint64    `xorm:"notnull index" json:"uid" yaml:"uid" xml:"uid,attr"` // user ID
	CTime  time.Time `xorm:"created 'ctime' notnull default CURRENT_TIMESTAMP" json:"ctime" yaml:"ctime" xml:"ctime,attr"`
	UTime  time.Time `xorm:"updated 'utime' notnull default CURRENT_TIMESTAMP" json:"utime" yaml:"utime" xml:"utime,attr"`
	Status byte      `xorm:"notnull default 1" json:"status" yaml:"status" xml:"status"`
	Device `xorm:"extends" yaml:",inline"`
}

func (Session) TableName() string {
	return "session"
}
