package capability

type Capability interface {
	MaybeGetProp(propName string) (handled bool, value interface{}, err error)
	MaybeHandle(method string, params interface{}) (handled bool, data interface{}, err error)
}
