package tcp

// Пара запрос-ответ для инициализации создания лайна.
type Request1013 struct {
	ID     int
	Target string
}
type Response1013 struct {
	ID      int
	Status  bool
	Address string
	Target  string
}

// Пакет смерти. Посылается читерам, которые не прошли проверку на Point
// Так же их IP заносится в черный список Point'a
type Command777 struct {
	ID     int
	Status bool
}

// Команда от клиента регистрации о хосте для подключения.
// На данный момент обе ноды проверили друг друга на Поинте, так что поидее все хорошо.
// Последний аккорд в создании Лайна.
type Command88 struct {
	ID           int
	HostAddr     string
	ThoseNodeKey string
}
