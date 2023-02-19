package types

import (
	"reflect"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/unix"
)

type Epoll struct {
	fd int
	// connections map[int]*websocket.Conn
	lock *sync.RWMutex
}

func MkEpoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &Epoll{
		fd:   fd,
		lock: &sync.RWMutex{},
		// connections: make(map[int]*websocket.Conn),
	}, nil
}

func (e *Epoll) Add(conn *websocket.Conn) (int, error) {
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return 0, err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	// e.connections[fd] = conn
	// if len(e.connections)%100 == 0 {
	// 	log.Printf("Total number of connections: %v", len(e.connections))
	// }
	return fd, nil
}

func (e *Epoll) Remove(conn *websocket.Conn) error {
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	// delete(e.connections, fd)
	// if len(e.connections)%100 == 0 {
	// 	log.Printf("Total number of connections: %v", len(e.connections))
	// }
	return nil
}

func (e *Epoll) Wait() ([]int, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	// var connections []*websocket.Conn
	var fds []int
	for i := 0; i < n; i++ {
		// conn := e.connections[int(events[i].Fd)]
		// connections = append(connections, conn)
		fds = append(fds, int(events[i].Fd))
	}
	return fds, nil
}

func websocketFD(conn *websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
