package command

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"reflect"
	"time"

	"dcore/codebase/module/http/persistence"
	dcutil "dcore/codebase/util"
)

const (
	RegName  = "reg"
	portName = "port"
)

type Reg struct {
	Name        string
	port        *portDesc
	paramsCount int
}

func NewRegResponse() *Reg {
	return &Reg{Name: RegName, port: newPortDescResponse(), paramsCount: 1}
}

func NewRegRequest(port int) *Reg {
	return &Reg{Name: RegName, port: newPortDescRequest(port), paramsCount: 1}
}

func (this Reg) RequestToURL(pointAddress string) string {
	// ip/reg?port=<int>
	return fmt.Sprintf("%s/%s?%s=%d", pointAddress, this.Name, this.port.name, this.port.value)
}

func (this *Reg) Parse(queryParams url.Values) (err error) {
	if err = checkParamsCount(RegName, queryParams, this.paramsCount); err != nil {
		return err
	}
	if this.port.value, err = tryExtractIntParam(LookName, this.port.name, queryParams); err != nil {
		return err
	}
	return nil
}

func (this Reg) PrepareResponse(params ...interface{}) []byte {
	redis := params[0].(*persistence.RedisModule)
	secret := params[1].(string)
	remoteAddr := params[2].(string)

	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s!@#$^&*()%d!@#$^&*()%s", remoteAddr, time.Now().Nanosecond(), secret)))
	key := hex.EncodeToString(hash.Sum(nil))
	ip := dcutil.RemovePortFromAddressString(remoteAddr)

	redis.AddNode(key, ip, this.port.value)
	return []byte(key)
}

///////////////////////////////////////////////////////////////////////////////

type portDesc struct {
	value    int
	name     string
	portType reflect.Kind // TODO : Not sure about it. Not used now
}

func newPortDescRequest(port int) *portDesc {
	return &portDesc{value: port, name: portName, portType: reflect.Int}
}
func newPortDescResponse() *portDesc {
	return &portDesc{name: portName, portType: reflect.Int}
}
