package capability

import "github.com/nickw444/miio-go/common"

type Info struct {
	Model string
}

func (i *Info) MaybeGetProp(propName string) (handled bool, value interface{}, err error) {
	return false, nil, nil
}

func (i *Info) MaybeHandle(method string, params interface{}) (handled bool, data interface{}, err error) {
	if method == "miIO.info" {
		info := common.DeviceInfo{
			Model: "yeelink.light.color1",
		}
		return true, info, nil
	}

	return false, nil, nil
}
