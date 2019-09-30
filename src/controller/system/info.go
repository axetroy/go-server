// Copyright 2019 Axetroy. All rights reserved. MIT license.
package system

import (
	"errors"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"os/user"
	"runtime"
	"time"
)

type Info struct {
	Username         string         `json:"username"`            // 当前用户名
	Host             host.InfoStat  `json:"host"`                // 操作系统信息
	Avg              load.AvgStat   `json:"avg"`                 // 负载信息
	Arch             string         `json:"arch"`                // 系统架构, 32/64位
	CPU              []cpu.InfoStat `json:"cpu"`                 // CPU信息
	RAMAvailable     uint64         `json:"ram_available"`       // 系统内存是否可供程序使用
	RAMTotal         uint64         `json:"ram_total"`           // 总内存大小
	RAMFree          uint64         `json:"ram_free"`            // 目前可用内存
	RAMUsedBy        uint64         `json:"ram_used_by"`         // 程序占用的内存
	RAMUsedByPercent float64        `json:"ram_used_by_percent"` // 程序占用的内存百分比
	UploadUsageStat  disk.UsageStat `json:"upload_usage_stat"`   // 上传目录的使用量统计
	Time             string         `json:"time"`                // 系统当前时间
	Timezone         string         `json:"timezone"`            // 当前服务器所在的时区
}

func GetSystemInfo() (res schema.Response) {
	var (
		err             error
		data            Info
		hostInfo        *host.InfoStat
		CPUInfo         []cpu.InfoStat
		avgStat         *load.AvgStat
		uploadUsageStat *disk.UsageStat
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}

	}()

	v, _ := mem.VirtualMemory()

	if CPUInfo, err = cpu.Info(); err != nil {
		return
	}

	if hostInfo, err = host.Info(); err != nil {
		return
	}

	if avgStat, err = load.Avg(); err != nil {
		return
	}

	if uploadUsageStat, err = disk.Usage(config.Upload.Path); err != nil {
		return
	}

	var u *user.User

	if u, err = user.Current(); err != nil {
		return
	}

	t := time.Now()

	data = Info{
		Username:         u.Username,
		Host:             *hostInfo,
		Arch:             runtime.GOARCH,
		Avg:              *avgStat,
		CPU:              CPUInfo,
		RAMAvailable:     v.Available,
		RAMTotal:         v.Total,
		RAMFree:          v.Free,
		RAMUsedBy:        v.Used,
		RAMUsedByPercent: v.UsedPercent,
		UploadUsageStat:  *uploadUsageStat,
		Time:             t.Format(time.RFC3339Nano),
		Timezone:         t.Location().String(),
	}

	return
}

func GetSystemInfoRouter(context *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetSystemInfo()
}
