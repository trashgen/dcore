package persistance

import "log"

type BlackListModule struct {
    // Postgres entry
}

func NewBlackListModule() *BlackListModule {
    return &BlackListModule{}
}

func (this *BlackListModule) Save(ip string) error {
    // TODO : Persist to Black list
    log.Printf("Save address [%s] to black list\n", ip)
    return nil
}

func (this *BlackListModule) Close() {}