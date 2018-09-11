package command

import (
    "net/url"
    "reflect"
    "dcore/codebase/module/http/persistence"
    "errors"
    "fmt"
    "strconv"
)

const (
    keyName    = "key"
    targetName = "target"
)

// Контракт не для полиморфизма, а идеологический - чтобы логика работы всех Команд была максимально стандартизирована.
// Так же стоит учитывать, что каждая команда уникальна, и помимо этих методов всегда может быть что-то еще в любом количестве.
// Однако интерфейс использования должен быть единым для всех.
// Этот интерфейс предполагает:
// 1. Формирование запроса на сервер;
// 2. Обработка на сервере;
// 3. Формирование ответа.
// Обработка ответа на клиенте происходит на уровне, который совершал вызов на формирование запроса.
type Command interface {
    // Клиент: Формирует строку запроса на заданный поинт учитывая весь набор параметров строки запроса
    RequestToURL(pointAddress string) string
    // Сервер: Парсит параметры запроса, валидирует их, заполняет внутренние структуры значениями
    Parse(queryParams url.Values) (err error)
    // Работа с полученными параметрами - хендлер. В итоге формирует слайс байтов для отправки ответа.
    // Даже если логика обработки одной команды похожа на другую - не надо выделять её в метод - они уникальны,
    // и потому должны полностью быть независимы от других. Во избежание.
    // call like: Command.PrepareResponse(net.Conn.RemoteAddr as string, ID as int)
    PrepareResponse(params... interface{}) []byte
    // use like (human readable) (you know the type, you know the order):
    // remoteAddress := params[0].(string)
    // key := params[1].(int)
}

                                            ///////////////////////
                                            // Common ParamDescs //
                                            ///////////////////////

type keyDesc struct {
    name     string
    value    string
    kindType reflect.Kind
}
func newKeyDescRequest(key string) *keyDesc {
    return &keyDesc{value: key, name: keyName, kindType: reflect.String}
}
func newKeyDescResponse() *keyDesc {
    return &keyDesc{name: keyName, kindType: reflect.String}
}

type targetDesc struct {
    name     string
    value    string
    kindType reflect.Kind
}
func newTargetDescRequest(target string) *targetDesc {
    return &targetDesc{value: target, name: targetName, kindType: reflect.String}
}
func newTargetDescResponse() *targetDesc {
    return &targetDesc{name: targetName, kindType: reflect.String}
}

                                              ///////////////////
                                              // Common Checks //
                                              ///////////////////

func checkParamsCount(cmdName string, toCheck url.Values, fromCheck int) error {
    if len(toCheck) != fromCheck {
        return errors.New(fmt.Sprintf("request '%s' not valid: bad params count - wait [%d] has [%d]", cmdName, fromCheck, len(toCheck)))
    }
    return nil
}

func tryExtractKeyParam(cmdName string, keyParamName string, queryParams url.Values, redis *persistence.RedisModule) (out string, err error) {
    out = queryParams.Get(keyParamName)
    if len(out) == 0 {
        return "", errors.New(fmt.Sprintf("request '%s' not valid: bad param name - wait [%s]", cmdName, keyParamName))
    }
    if _, has := redis.GetNode(out); ! has {
        return "", errors.New(fmt.Sprintf("request '%s' not valid: %s not found", keyParamName, cmdName))
    }
    return out, nil
}

func tryExtractIntParam(cmdName string, paramName string, queryParams url.Values) (out int, err error) {
    outValue := queryParams.Get(paramName)
    if len(outValue) == 0 {
        return 0, errors.New(fmt.Sprintf("request '%s' not valid: bad param name - wait [%s]", cmdName, paramName))
    }
    if out, err = strconv.Atoi(outValue); err != nil {
         return 0, errors.New(fmt.Sprintf("request '%s' not valid: bad param [%s] type - wait [int] has [%s]", cmdName, paramName, reflect.TypeOf(outValue).Name()))
    }
    return out, nil
}

func tryExtractStringParam(cmdName string, paramName string, queryParams url.Values) (out string, err error) {
    out = queryParams.Get(paramName)
    if len(out) == 0 {
        return "", errors.New(fmt.Sprintf("request '%s' not valid: bad param name - wait [%s]", cmdName, paramName))
    }
    return out, nil
}