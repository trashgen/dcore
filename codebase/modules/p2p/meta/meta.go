package meta

import "net"

/*
 * Эти интерфейсы специально выделены в отдельный пакадж, дабы импортом тянуть минимум.
 * Смысл этих интерфейсов в обработке запросов/ответов на П2П канале.
 * Задаются реализации (паттерн Стратегия) в момент создания Node.
 * Принцип работы такой же как стандартный net.http.ListenAndServe() (исключая возможность 'nil').
 *
 * TODO : DEPRECATED. выпилить их и сделать reghost/regclient handlers без этого изврата. Идея провалилась.
 */

// TCP. Используется сервером в ответ на поступающие с клиента запросы.
type RequestHandler interface {
    Run(data string, conn net.Conn) (response []byte, hasResponse bool, err error)
}

// TCP. Когда сервер ответил на запрос - надо обработать, что он там придумал.
type ResponseHandler interface {
    Run(data string, conn net.Conn) (response []byte, hasResponse bool, err error)
}