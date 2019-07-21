package persistence

import "log"

type HDDPersist interface {
	Save(data string)
	CheckExists(id string) (exists bool)
	Close()
}

type MockPersistModule struct{}

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

type MockRedisModule struct {
	nodes  map[string]*ConnectionID
	points map[string]*ConnectionID // not used now
}

func NewMockRedisModule() *MockRedisModule {
	return &MockRedisModule{
		nodes:  make(map[string]*ConnectionID, 16),
		points: make(map[string]*ConnectionID, 16)}
}

func (this *MockRedisModule) AddNode(key string, ip string, port int) (out *ConnectionID) {
	out = NewConnectionID(key, ip, port)
	this.nodes[key] = out
	return out
}

func (this *MockRedisModule) GetNode(key string) (result *ConnectionID, has bool) {
	// LOL I just can't do: "return this.nodes[key]", 'cos "has" not returns by default
	result, has = this.nodes[key]
	return result, has
}

func (this *MockRedisModule) GetRandomNodes(count int) map[string]*ConnectionID {
	return this.nodes
}

func (this *MockRedisModule) GetAllNodes() map[string]*ConnectionID {
	return this.nodes
}

func (this *MockRedisModule) GetAllNodesAsSlice() []*ConnectionID {
	out := make([]*ConnectionID, 0, len(this.nodes))
	for _, v := range this.nodes {
		out = append(out, v)
	}

	return out
}

func (this *MockRedisModule) RemoveNode(key string) {
	delete(this.nodes, key)
}

func (this *MockRedisModule) AddPoint(key string, ip string, port int) (out *ConnectionID) {
	out = NewConnectionID(key, ip, port)
	this.points[key] = out
	return out
}

func (this *MockRedisModule) GetPoint(key string) (result *ConnectionID, has bool) {
	result, has = this.points[key]
	return result, has
}

func (this *MockRedisModule) GetAllPoints() []*ConnectionID {
	out := make([]*ConnectionID, 0, len(this.points))
	for _, v := range this.points {
		out = append(out, v)
	}

	return out
}

func (this *MockRedisModule) RemovePoint(key string) {
	delete(this.points, key)
}
