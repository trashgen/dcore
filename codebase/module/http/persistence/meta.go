package persistence

import "log"

type HDDPersist interface {
    Save(data string)
    CheckExists(id string) (exists bool)
    Close()
}

type MockPersistModule struct {}

func NewMockPersistModule() *MockPersistModule {
    return &MockPersistModule{}
}

func (this *MockPersistModule) Save(ip string) {
    log.Printf("Save [%s] to Database\n", ip)
}

func (this *MockPersistModule) CheckExists(id string) (exists bool) {
    log.Printf("Checking key [%s] exists is true\n", id)
    return true
}

func (this *MockPersistModule) Close() {
    log.Print("Close MockModule\n")
}