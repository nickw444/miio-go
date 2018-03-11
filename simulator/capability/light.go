package capability

type Light struct {
	brightness int
	color      int
}

func (l *Light) MaybeGetProp(propName string) (handled bool, value interface{}, err error) {
	switch propName {
	case "bright":
		return true, l.brightness, nil
	case "hsv":
		return true, 0, nil
	case "rgb":
		return true, 0, nil
	default:
		return false, nil, nil
	}
}

func (l *Light) MaybeHandle(method string, params interface{}) (handled bool, data interface{}, err error) {
	switch method {
	case "set_bright":
		return true, nil, nil
	case "set_rgb":
		return true, nil, nil
	case "set_hsv":
		return true, nil, nil
	default:
		return false, nil, nil
	}
}
