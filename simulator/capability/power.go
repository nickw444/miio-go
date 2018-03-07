package capability

type Power struct {
	power bool
}

func (p *Power) MaybeGetProp(propName string) (handled bool, value interface{}, err error) {
	if propName == "power" {
		var power string
		if p.power {
			power = "on"
		} else {
			power = "off"
		}
		return true, power, nil
	}

	return false, nil, nil
}

func (p *Power) MaybeHandle(method string, params interface{}) (handled bool, data interface{}, err error) {
	if method == "set_power" {
		value := params.([]interface{})[0].(string)
		if value == "on" {
			p.power = true
		} else if value == "off" {
			p.power = false
		}
		return true, nil, nil
	}
	return false, nil, nil
}
