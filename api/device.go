package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lanthora/cacao/model"
)

func DeviceShow(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	devices := model.GetDevicesByUserID(user.ID)

	type devinfo struct {
		DevID          uint   `json:"devid"`
		NetID          uint   `json:"netid"`
		Vmac           string `json:"vmac"`
		IP             string `json:"ip"`
		EIP           string `json:"eip"`
		EPR           string `json:"epr"`
		Online         bool   `json:"online"`
		RX             uint64 `json:"rx"`
		TX             uint64 `json:"tx"`
		OS             string `json:"os"`
		Version        string `json:"version"`
		Hostname       string `json:"hostname"`
		Country        string `json:"country"`
		Region         string `json:"region"`
		LastActiveTime string `json:"lastActiveTime"`
	}

	response := make([]devinfo, 0)
	for _, d := range devices {
		response = append(response, devinfo{
			DevID:          d.ID,
			NetID:          d.NetID,
			Vmac:           d.VMac,
			IP:             d.IP,
			EIP:            d.EIP,
			EPR:            d.EPR,
			Online:         d.Online,
			RX:             d.RX,
			TX:             d.TX,
			OS:             d.OS,
			Version:        d.Version,
			Hostname:       d.Hostname,
			Country:        d.Country,
			Region:         d.Region,
			LastActiveTime: d.UpdatedAt.Format(time.DateTime),
		})
	}

	setResponseData(c, gin.H{
		"devices": response,
	})
}

func DeviceDelete(c *gin.Context) {
	var request struct {
		DevID uint `json:"devid"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		setErrorCode(c, InvalidRequest)
		return
	}

	deviceModel := model.GetDeviceByDevID(request.DevID)
	if deviceModel.Online {
		setErrorCode(c, CannotDeleteOnlineDevice)
		return
	}

	netModel := model.GetNetByNetID(deviceModel.NetID)

	user := c.MustGet("user").(*model.User)
	if user.ID != netModel.UserID {
		setErrorCode(c, DeviceNotExists)
		return
	}

	deviceModel.Delete()
	setResponseData(c, nil)
}
