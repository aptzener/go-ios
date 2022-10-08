package api

import (
	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/instruments"
	"github.com/danielpaulus/go-ios/ios/screenshotr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Info gets device info
// Info                godoc
// @Summary      Get lockdown info for a device by udid
// @Description  Returns all lockdown values and additional instruments properties for development enabled devices.
// @Tags         general_device_specific
// @Produce      json
// @Param        udid  path      string  true  "device udid"
// @Success      200  {object}  map[string]interface{}
// @Router       /device/{udid}/info [get]
func Info(c *gin.Context) {
	device := c.MustGet(IOS_KEY).(ios.DeviceEntry)

	allValues, err := ios.GetValuesPlist(device)
	if err != nil {
		print(err)
	}
	svc, err := instruments.NewDeviceInfoService(device)
	if err != nil {
		log.Debugf("could not open instruments, probably dev image not mounted %v", err)
	}
	if err == nil {
		info, err := svc.NetworkInformation()
		if err != nil {
			log.Debugf("error getting networkinfo from instruments %v", err)
		} else {
			allValues["instruments:networkInformation"] = info
		}
		info, err = svc.HardwareInformation()
		if err != nil {
			log.Debugf("error getting hardwareinfo from instruments %v", err)
		} else {
			allValues["instruments:hardwareInformation"] = info
		}
	}
	c.IndentedJSON(http.StatusOK, allValues)
}

// Screenshot grab screenshot from a device
// Screenshot                godoc
// @Summary      Get screenshot for device
// @Description Takes a png screenshot and returns it.
// @Tags         general_device_specific
// @Produce      png
// @Param        udid  path      string  true  "device udid"
// @Success      200  {object}  []byte
// @Router       /device/{udid}/screenshot [get]
func Screenshot(c *gin.Context) {
	device := c.MustGet(IOS_KEY).(ios.DeviceEntry)
	conn, err := screenshotr.New(device)
	log.Error(err)
	b, _ := conn.TakeScreenshot()

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "application/octet-stream", b)
}
