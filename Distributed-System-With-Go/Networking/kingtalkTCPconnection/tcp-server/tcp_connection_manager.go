package tcp

import (
	"fmt"
	cmap "kingtalk_tcp/concurrent-map"
	"sync"
	"sync/atomic"
)

type connectionMap struct {
	sync.RWMutex
	conns    map[uint64]Connection
	disposed bool
}

// NewConnectionMap func;
func NewConnectionMap() *connectionMap {
	return &connectionMap{
		RWMutex: sync.RWMutex{},
		conns:   make(map[uint64]Connection),
	}
}

const connectionMapNum = 32

// ConnectionManager type
type ConnectionManager struct {
	connectionMaps    sync.Map
	connectionMapsOld cmap.ConcurrentMap[string, *connectionMap]
	// disposeFlag       bool
	disposeFlagInt *int32
	disposeOnce    sync.Once
	disposeWait    sync.WaitGroup
}

// NewConnectionManager func
func NewConnectionManager() *ConnectionManager {
	manager := &ConnectionManager{
		disposeFlagInt: new(int32),
		disposeOnce:    sync.Once{},
		disposeWait:    sync.WaitGroup{},
	}
	for idx := 0; idx < connectionMapNum; idx++ {
		connectMap := NewConnectionMap()
		manager.connectionMaps.Store(idx, connectMap)
	}

	return manager
}

// Dispose func
func (manager *ConnectionManager) Dispose() {
	manager.disposeOnce.Do(func() {
		// manager.disposeFlag = true
		atomic.StoreInt32(manager.disposeFlagInt, 1)
		for i := 0; i < connectionMapNum; i++ {
			var conns *connectionMap
			connsValue, ok := manager.connectionMaps.Load(i)
			if !ok {
				conns = NewConnectionMap()
				manager.connectionMaps.Store(i, conns)
			} else {
				conns = connsValue.(*connectionMap)
			}

			conns.Lock()
			conns.disposed = true
			for _, conn := range conns.conns {
				conn.Close()
			}
			conns.Unlock()
		}
		manager.disposeWait.Wait()
	})
}

// GetConnection func
func (manager *ConnectionManager) GetConnection(connID uint64) Connection {
	var conns *connectionMap

	connID = connID & 0xffffffffffffff
	connsValue, ok := manager.connectionMaps.Load(connID % connectionMapNum)
	fmt.Println("ConnectionManager GetConnection connId: %d and conn: %v and ok: %v", connID, connsValue, ok)
	if !ok {
		conns = NewConnectionMap()
		manager.connectionMaps.Store(connID%connectionMapNum, conns)
		return nil
	}

	conns = connsValue.(*connectionMap)

	conns.RLock()
	defer conns.RUnlock()

	conn, _ := conns.conns[connID]
	fmt.Println("ConnectionManager GetConnection connId: %d and conn: %v", connID, conn)
	return conn
}

// SendMessageToConnection func;
// Get a connection by ConnID and send msg to it;
func (manager *ConnectionManager) SendMessageToConnection(connID uint64, msg interface{}) error {
	var conns *connectionMap

	connID = connID & 0xffffffffffffff
	connsValue, ok := manager.connectionMaps.Load(connID % connectionMapNum)
	fmt.Println("ConnectionManager SendMessageToConnection connId: %d and conn: %v and ok: %v", connID, connsValue, ok)
	if !ok {
		conns = NewConnectionMap()
		manager.connectionMaps.Store(connID%connectionMapNum, conns)

		err := fmt.Errorf("Not found connection by id: %d", connID)
		return err
	}

	conns, ok = connsValue.(*connectionMap)
	if !ok {
		err := fmt.Errorf("Can not parse to connectionMap: %d", connID)
		fmt.Println("ConnectionManager::SendMessageToConnection - Error: %+v", err)
		return err
	}

	conns.RLock()
	defer conns.RUnlock()

	conn, _ := conns.conns[connID]
	fmt.Println("ConnectionManager SendMessageToConnection connId: %d and conn: %v", connID, conn)
	return conn.Send(msg)
}

// SendMessageToAllConnection func;
// Send msg to connection that was connected;
func (manager *ConnectionManager) SendMessageToAllConnection(msg interface{}) {
	fmt.Println("ConnectionManager::SendMessageToAllConnection - Message: %+v", msg)
	manager.connectionMaps.Range(func(mapID, connsValue interface{}) bool {
		conns, ok := connsValue.(*connectionMap)
		if !ok {
			err := fmt.Errorf("Can not parse to connectionMap: %+v", mapID)
			fmt.Println("ConnectionManager::SendMessageToConnection - Parser Error: %+v", err)
			return true
		}

		conns.RLock()
		for _, conn := range conns.conns {
			if err := conn.Send(msg); err != nil {
				fmt.Println("ConnectionManager::SendMessageToConnection - Send via conn Error: %+v", err)
			}

		}
		conns.RUnlock()

		return true
	})
}

// putConnection func
func (manager *ConnectionManager) PutConnection(conn Connection) {
	connID := conn.GetConnID()
	var conns *connectionMap

	connID = connID & 0xffffffffffffff
	connsValue, ok := manager.connectionMaps.Load(connID % connectionMapNum)
	fmt.Println("putConnection connId: %d and conn: %v and ok: %v", connID, connsValue, ok)
	if !ok {
		conns = NewConnectionMap()
		manager.connectionMaps.Store(connID%connectionMapNum, conns)
	} else {
		conns, ok = connsValue.(*connectionMap)
		if !ok {
			err := fmt.Errorf("Can not parse to connectionMap: %d", connID)
			fmt.Println("ConnectionManager::SendMessageToConnection - Error: %+v", err)
			return
		}
	}

	conns.Lock()
	defer conns.Unlock()

	if conns.disposed {
		conn.Close()
		return
	}
	fmt.Println("putConnection connId: %d and conn: %v setToMap", connID, conn)
	conns.conns[connID] = conn
	manager.disposeWait.Add(1)
}

// delConnection func
func (manager *ConnectionManager) delConnection(conn Connection) {
	if atomic.LoadInt32(manager.disposeFlagInt) == 1 {
		manager.disposeWait.Done()
		return
	}

	connID := conn.GetConnID()
	connID = connID & 0xffffffffffffff
	var conns *connectionMap

	connsValue, ok := manager.connectionMaps.Load(connID % connectionMapNum)
	fmt.Println("delConnection connId: %d and conn: %v and ok: %v", connID, connsValue, ok)
	if !ok {
		conns = NewConnectionMap()
		manager.connectionMaps.Store(connID%connectionMapNum, conns)

		return
	}

	conns, ok = connsValue.(*connectionMap)
	if !ok {
		err := fmt.Errorf("Can not parse to connectionMap: %d", connID)
		fmt.Println("ConnectionManager::SendMessageToConnection - Error: %+v", err)
		return
	}

	conns.Lock()
	defer conns.Unlock()
	fmt.Println("delConnection connId: %d delete", connID)
	delete(conns.conns, connID)
	manager.disposeWait.Done()
}

// delConnection func
func (manager *ConnectionManager) delConnectionByID(connID uint64) {
	connID = connID & 0xffffffffffffff
	if atomic.LoadInt32(manager.disposeFlagInt) == 1 {
		manager.disposeWait.Done()
		return
	}

	var conns *connectionMap

	connsValue, ok := manager.connectionMaps.Load(connID % connectionMapNum)
	fmt.Println("delConnectionByID connId: %d and conn: %v and ok: %v", connID, connsValue, ok)
	if !ok {
		conns = NewConnectionMap()
		manager.connectionMaps.Store(connID%connectionMapNum, conns)

		return
	}

	conns, ok = connsValue.(*connectionMap)
	if !ok {
		err := fmt.Errorf("Can not parse to connectionMap: %d", connID)
		fmt.Println("ConnectionManager::SendMessageToConnection - Error: %+v", err)
		return
	}

	conns.Lock()
	defer conns.Unlock()
	fmt.Println("delConnectionByID connId: %d delete", connID)
	delete(conns.conns, connID)
	manager.disposeWait.Done()
}
