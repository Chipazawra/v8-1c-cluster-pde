package pool

import (
	"context"
	"errors"
	"github.com/khorevaa/ras-client/messages"
	"github.com/khorevaa/ras-client/serialize/esig"
	uuid "github.com/satori/go.uuid"
	"sync"
	"sync/atomic"
	"time"
)

// Что есть
// 1. Клиент подключения к сервесу.
// 2. На клиенте есть несколько endpoint
//   - требует авторизации на кластере
//   - требует авториазции базы ( 1 раз для каждой базы,
//     если меняешь базы или используется другиой пароль доступа
//     нужна переавторизация
// 3. Деление endpoint по:
//   - общая инфорация (сообщения в целом на кластер)
//   - информация по базам (сообщения с ключем infobase
//

// Надо организовать пул клиентов
// пул endpoint для переключения авторизации по базам
// Клиент читает сообщения для endpoint и отправляем его в ответ.

/*

Алгоритм работы с подключением

1. Создается клиент подключения (далее клиент)
2. Создается 1 проверочное соединение.
2.1. Выполняются начальные команды:
	- NewNegotiateMessage(protocolVersion, c.codec.Version()))
	- &ConnectMessage{params: map[string]interface{}{
		"connect.timeout": int64(2000) Ъ}
2.2. Соединение переходит в ожидание, при отсутствии ошибок. При ошибке клиент не создается.
3. Основной цикл работы
   Запрос на данные - Ответ пользователю
3.1. Открывает точку обмена
3.2. Авторизация на кластере -> Возможна ошибка прав/авторизации -> возврат ошибки
3.3. Если необходимо авторизация на информациооной базе -> Возможна ошибка прав/авторизации -> возврат ошибки
3.4. Выполнение запроса -> Возможна ошибка парсинга -> возврат ошибки
3.5. Ожидание ответа. Для запросов (VIOD_MESSAGE) не ошидания ответа, переход сразу к пункту 3.8
3.6. Разбор ответа. -> возможна ошибка запроса
3.7. Отправка ответа пользователю
3.8. Перевод точки обмена в ожидание. По двум критериям
	- запрос был только на данные кластера (переиспользование для аналогичных запросов)
	- была авторизация по ИБ. (переиспользование для запросов по данной базе)
	  по истечении, н-минут переход в исползование по другим базам, с повторной авторизацией

4. Цикла работы точки обмена
4.1. Открытие
4.2. Отправка собщения
	 Установка блокировки на соединение -> Запись ланных в соединение
4.2. Чтение данных из соединения -> Получение сообщения
     Снятие блокировки на соедиенение
4.3. Ожидание для срока жизни / повторение цикла с пункта 4.2.
4.4. Закрытие точки обмена
4.5. Завершение при закрытии соединения

5. Работа с открытым соединением
5.1. Блокировка использования другими точками обмена
5.2. Запись данных
5.3. Ожидание ответа -> чтение данных. При открытой точке обмена всегда приходит ответ на посланный запрос
	 Даже если он не требует явного ответа, например Авторизация на кластере или в информационной базе
5.4. Разблокировка по таймауту или при получении ответа


*/

var (
	ErrClosed         = errors.New("protocol: pool is closed")
	ErrUnknownMessage = errors.New("protocol: unknown message packet")
	ErrPoolTimeout    = errors.New("protocol: endpoint pool timeout")
)

var timers = sync.Pool{
	New: func() interface{} {
		t := time.NewTimer(time.Hour)
		t.Stop()
		return t
	},
}

var _ EndpointPool = (*endpointPool)(nil)

func NewEndpointPool(opt *Options) EndpointPool {
	p := &endpointPool{
		opt:             opt,
		queue:           make(chan struct{}, opt.PoolSize),
		conns:           make([]*Conn, 0, opt.PoolSize),
		idleConns:       make([]*Conn, 0, opt.PoolSize),
		authInfobaseIdx: make(map[uuid.UUID]struct{ user, password string }),
		authClusterIdx:  make(map[uuid.UUID]struct{ user, password string }),
	}

	p.connsMu.Lock()
	p.checkMinIdleConns()
	p.connsMu.Unlock()

	if opt.IdleTimeout > 0 && opt.IdleCheckFrequency > 0 {
		go p.reaper(opt.IdleCheckFrequency)
	}

	return p
}

type EndpointPool interface {
	NewEndpoint(ctx context.Context) (*Endpoint, error)
	CloseEndpoint(endpoint *Endpoint) error

	Get(ctx context.Context, sig esig.ESIG) (*Endpoint, error)
	Put(ctx context.Context, endpoint *Endpoint)
	Remove(ctx context.Context, endpoint *Endpoint, err error)

	Len() int
	IdleLen() int

	Close() error

	SetAgentAuth(user, password string)
	SetClusterAuth(id uuid.UUID, user, password string)
	SetInfobaseAuth(id uuid.UUID, user, password string)
	GetClusterAuth(id uuid.UUID) (user, password string)
	GetInfobaseAuth(id uuid.UUID) (user, password string)
}

type Pooler interface {
	NewConn(context.Context) (*Conn, error)
	CloseConn(*Conn) error

	Get(context.Context) (*Conn, error)
	Put(context.Context, *Conn)
	Remove(context.Context, *Conn, error)

	Len() int
	IdleLen() int
	Close() error
}

type endpointPool struct {
	opt *Options

	dialErrorsNum uint32 // atomic

	_closed uint32 // atomic

	lastDialErrorMu sync.RWMutex
	lastDialError   error

	queue chan struct{}

	poolSize     int
	idleConnsLen int

	connsMu   sync.Mutex
	conns     []*Conn
	idleConns IdleConns

	authClusterIdx  map[uuid.UUID]struct{ user, password string }
	authInfobaseIdx map[uuid.UUID]struct{ user, password string }
	authAgent       struct{ user, password string }
}

func (p *endpointPool) NewEndpoint(ctx context.Context) (*Endpoint, error) {

	if p.closed() {
		return nil, ErrClosed
	}

	err := p.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	for {
		p.connsMu.Lock()
		endpoint := p.popIdle(esig.ESIG{})
		p.connsMu.Unlock()

		if endpoint == nil {
			break
		}

		if p.isStaleConn(endpoint.conn) {
			_ = p.CloseConn(endpoint.conn)
			continue
		}

		if !endpoint.Inited {
			endpoint, err = p.openEndpoint(ctx, endpoint.conn)

			if err != nil {
				return nil, err
			}

		}

		return endpoint, nil
	}

	newConn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()
		return nil, err
	}

	endpoint, err := p.openEndpoint(ctx, newConn)

	return endpoint, err

}

func (p *endpointPool) Put(ctx context.Context, cn *Endpoint) {
	if !cn.conn.pooled {
		p.Remove(ctx, cn, nil)
		return
	}

	p.connsMu.Lock()
	p.idleConns = append(p.idleConns, cn.conn)
	p.idleConnsLen++
	p.connsMu.Unlock()
	p.freeTurn()
}

// Get returns existed connection from the pool or creates a new one.
func (p *endpointPool) Get(ctx context.Context, sig esig.ESIG) (*Endpoint, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	err := p.waitTurn(ctx)
	if err != nil {
		return nil, err
	}

	for {
		p.connsMu.Lock()
		endpoint := p.popIdle(sig)
		p.connsMu.Unlock()

		if endpoint == nil {
			break
		}

		if p.isStaleConn(endpoint.conn) {
			_ = p.CloseConn(endpoint.conn)

			continue
		}

		if !endpoint.Inited {
			endpoint, err = p.openEndpoint(ctx, endpoint.conn)

			if err != nil {
				return nil, err
			}

		}

		return endpoint, nil
	}

	newConn, err := p.newConn(ctx, true)
	if err != nil {
		p.freeTurn()

		return nil, err
	}

	endpoint, err := p.openEndpoint(ctx, newConn)

	return endpoint, err
}

func (p *endpointPool) Remove(_ context.Context, cn *Endpoint, _ error) {
	p.removeConnWithLock(cn.conn)
	p.freeTurn()
	_ = p.closeConn(cn.conn)
}

func (p *endpointPool) CloseConn(cn *Conn) error {
	p.removeConnWithLock(cn)
	return p.closeConn(cn)
}

func (p *endpointPool) SetAgentAuth(user, password string) {

	p.authAgent = struct{ user, password string }{user: user, password: password}

}

func (p *endpointPool) SetClusterAuth(id uuid.UUID, user, password string) {

	p.authClusterIdx[id] = struct{ user, password string }{user: user, password: password}

}

func (p *endpointPool) SetInfobaseAuth(id uuid.UUID, user, password string) {

	p.authInfobaseIdx[id] = struct{ user, password string }{user: user, password: password}

}

func (p *endpointPool) GetClusterAuth(id uuid.UUID) (user, password string) {

	return p.getAuth(p.authClusterIdx, id)
}

func (p *endpointPool) GetInfobaseAuth(id uuid.UUID) (user, password string) {

	return p.getAuth(p.authInfobaseIdx, id)
}

// Len returns total number of connections.
func (p *endpointPool) Len() int {
	p.connsMu.Lock()
	n := len(p.conns)
	p.connsMu.Unlock()
	return n
}

// IdleLen returns number of idle connections.
func (p *endpointPool) IdleLen() int {
	p.connsMu.Lock()
	n := p.idleConnsLen
	p.connsMu.Unlock()
	return n
}

func (p *endpointPool) Close() error {
	if !atomic.CompareAndSwapUint32(&p._closed, 0, 1) {
		return ErrClosed
	}

	var firstErr error
	p.connsMu.Lock()
	for _, cn := range p.conns {
		if err := p.closeConn(cn); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	p.conns = nil
	p.poolSize = 0
	p.idleConns = nil
	p.idleConnsLen = 0
	p.connsMu.Unlock()

	return firstErr
}

func (p *endpointPool) CloseEndpoint(*Endpoint) error {
	panic("implement me")
}

func (p *endpointPool) ReapStaleConns() (int, error) {
	var n int
	for {
		p.getTurn()

		p.connsMu.Lock()
		cn := p.reapStaleConn()
		p.connsMu.Unlock()

		p.freeTurn()

		if cn != nil {
			_ = p.closeConn(cn)
			n++
		} else {
			break
		}
	}
	return n, nil
}

func (p *endpointPool) openEndpoint(ctx context.Context, conn *Conn) (*Endpoint, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	if !conn.Inited {
		err := p.opt.InitConnection(ctx, conn)

		if err != nil {
			return nil, err
		}
		conn.Inited = true
	}

	openAck, err := p.opt.OpenEndpoint(ctx, conn)
	if err != nil {
		return nil, err
	}

	endpoint := NewEndpoint(openAck)
	endpoint.Inited = true
	endpoint.onRequest = p.onRequest
	endpoint.conn = conn
	conn.endpoints = append(conn.endpoints, endpoint)

	return endpoint, nil
}

func needAgentAuth(req messages.EndpointRequestMessage) bool {
	switch req.(type) {
	case *messages.GetAgentAdminsRequest, *messages.RegAgentAdminRequest, *messages.UnregAgentAdminRequest,
		*messages.RegClusterRequest, *messages.UnregClusterRequest:
		return true
	}

	return false
}

// Get returns existed connection from the pool or creates a new one.
func (p *endpointPool) onRequest(ctx context.Context, endpoint *Endpoint, req messages.EndpointRequestMessage) error {

	if needAgentAuth(req) {
		err := p.setAgentAuth(ctx, endpoint)

		if err != nil {
			return err
		}
	}

	sig := req.Sig()

	if esig.IsNul(sig) {
		return nil
	}

	if esig.Equal(endpoint.sig, sig) {
		return p.updateAuthIfNeed(ctx, endpoint, sig.High(), sig.Low())
	}

	err := p.updateAuthIfNeed(ctx, endpoint, sig.High(), sig.Low())
	if err != nil {
		return err
	}

	endpoint.sig = sig

	return nil
}

func (p *endpointPool) updateAuthIfNeed(ctx context.Context, endpoint *Endpoint, clusterID, infobaseID uuid.UUID) error {

	if user, password := p.GetClusterAuth(clusterID); !endpoint.CheckClusterAuth(user, password) {
		err := p.updateClusterAuth(ctx, endpoint, clusterID, user, password)
		if err != nil {
			return err
		}
	}

	if user, password := p.GetInfobaseAuth(infobaseID); !endpoint.CheckInfobaseAuth(user, password) {
		err := p.updateInfobaseAuth(ctx, endpoint, clusterID, user, password)
		if err != nil {
			return err
		}
	}

	return nil

}

func (p *endpointPool) updateClusterAuth(ctx context.Context, endpoint *Endpoint, clusterID uuid.UUID, user, password string) error {

	authMessage := endpoint.newEndpointMessage(messages.ClusterAuthenticateRequest{
		ClusterID: clusterID,
		User:      user,
		Password:  password,
	})

	message, err := endpoint.sendRequest(ctx, authMessage)

	if err != nil {
		return err
	}

	switch err := message.Message.(type) {

	case *messages.EndpointMessageFailure:

		return err

	}

	endpoint.SetClusterAuth(user, password)

	return nil

}

func (p *endpointPool) setAgentAuth(ctx context.Context, endpoint *Endpoint) error {

	authMessage := endpoint.newEndpointMessage(messages.AuthenticateAgentRequest{
		User:     p.authAgent.user,
		Password: p.authAgent.password,
	})

	message, err := endpoint.sendRequest(ctx, authMessage)

	if err != nil {
		return err
	}

	switch err := message.Message.(type) {

	case *messages.EndpointMessageFailure:

		return err

	}

	return nil

}

func (p *endpointPool) updateInfobaseAuth(ctx context.Context, endpoint *Endpoint, clusterID uuid.UUID, user, password string) error {

	authMessage := endpoint.newEndpointMessage(messages.AuthenticateInfobaseRequest{
		ClusterID: clusterID,
		User:      user,
		Password:  password,
	})

	message, err := endpoint.sendRequest(ctx, authMessage)

	if err != nil {
		return err
	}

	switch err := message.Message.(type) {

	case *messages.EndpointMessageFailure:

		return err

	}

	endpoint.SetInfobaseAuth(user, password)

	return nil

}

func (p *endpointPool) getAuth(idx map[uuid.UUID]struct{ user, password string }, id uuid.UUID) (user, password string) {

	if auth, ok := idx[id]; ok {
		user, password = auth.user, auth.password
		return
	}

	if auth, ok := idx[uuid.Nil]; ok {
		user, password = auth.user, auth.password
	}

	return
}

func (p *endpointPool) checkMinIdleConns() {
	if p.opt.MinIdleConns == 0 {
		return
	}
	for p.poolSize < p.opt.PoolSize && p.idleConnsLen < p.opt.MinIdleConns {
		p.poolSize++
		p.idleConnsLen++
		go func() {
			err := p.addIdleConn()
			if err != nil {
				p.connsMu.Lock()
				p.poolSize--
				p.idleConnsLen--
				p.connsMu.Unlock()
			}
		}()
	}
}

func (p *endpointPool) addIdleConn() error {
	cn, err := p.dialConn(context.TODO(), true)
	if err != nil {
		return err
	}

	p.connsMu.Lock()
	p.conns = append(p.conns, cn)
	p.idleConns = append(p.idleConns, cn)
	p.connsMu.Unlock()
	return nil
}

func (p *endpointPool) newConn(c context.Context, pooled bool) (*Conn, error) {
	cn, err := p.dialConn(c, pooled)
	if err != nil {
		return nil, err
	}

	cn.closer = p.opt.CloseEndpoint

	p.connsMu.Lock()
	p.conns = append(p.conns, cn)
	if pooled {
		// If pool is full remove the cn on next Put.
		if p.poolSize >= p.opt.PoolSize {
			cn.pooled = false
		} else {
			p.poolSize++
		}
	}
	p.connsMu.Unlock()
	return cn, nil
}

func (p *endpointPool) dialConn(c context.Context, pooled bool) (*Conn, error) {
	if p.closed() {
		return nil, ErrClosed
	}

	if atomic.LoadUint32(&p.dialErrorsNum) >= uint32(p.opt.PoolSize) {
		return nil, p.getLastDialError()
	}

	netConn, err := p.opt.Dialer(c)
	if err != nil {
		p.setLastDialError(err)
		if atomic.AddUint32(&p.dialErrorsNum, 1) == uint32(p.opt.PoolSize) {
			go p.tryDial()
		}
		return nil, err
	}

	cn := NewConn(netConn)
	cn.pooled = pooled
	return cn, nil
}

func (p *endpointPool) tryDial() {
	for {
		if p.closed() {
			return
		}

		conn, err := p.opt.Dialer(context.TODO())
		if err != nil {
			p.setLastDialError(err)
			time.Sleep(time.Second)
			continue
		}

		atomic.StoreUint32(&p.dialErrorsNum, 0)
		_ = conn.Close()
		return
	}
}

func (p *endpointPool) setLastDialError(err error) {
	p.lastDialErrorMu.Lock()
	p.lastDialError = err
	p.lastDialErrorMu.Unlock()
}

func (p *endpointPool) getLastDialError() error {
	p.lastDialErrorMu.RLock()
	err := p.lastDialError
	p.lastDialErrorMu.RUnlock()
	return err
}

func (p *endpointPool) getTurn() {
	p.queue <- struct{}{}
}

func (p *endpointPool) waitTurn(c context.Context) error {
	select {
	case <-c.Done():
		return c.Err()
	default:
	}

	select {
	case p.queue <- struct{}{}:
		return nil
	default:
	}

	timer := timers.Get().(*time.Timer)
	timer.Reset(p.opt.PoolTimeout)

	select {
	case <-c.Done():
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return c.Err()
	case p.queue <- struct{}{}:
		if !timer.Stop() {
			<-timer.C
		}
		timers.Put(timer)
		return nil
	case <-timer.C:
		timers.Put(timer)
		//atomic.AddUint32(&p.stats.Timeouts, 1)
		return ErrPoolTimeout
	}
}

func (p *endpointPool) freeTurn() {
	<-p.queue
}

func (p *endpointPool) popIdle(sig esig.ESIG) *Endpoint {
	if len(p.idleConns) == 0 {
		return nil
	}

	endpoint := p.idleConns.Pop(sig, p.opt.MaxOpenEndpoints)

	if endpoint == nil {
		return nil
	}

	p.idleConnsLen--
	p.checkMinIdleConns()
	return endpoint
}

func (p *endpointPool) removeConnWithLock(cn *Conn) {
	p.connsMu.Lock()
	p.removeConn(cn)
	p.connsMu.Unlock()
}

func (p *endpointPool) removeConn(cn *Conn) {
	for i, c := range p.conns {
		if c == cn {
			p.conns = append(p.conns[:i], p.conns[i+1:]...)
			if cn.pooled {
				p.poolSize--
				p.checkMinIdleConns()
			}

			return
		}
	}
}

func (p *endpointPool) closeConn(cn *Conn) error {
	if p.opt.OnClose != nil {
		_ = p.opt.OnClose(cn)
	}
	return cn.Close()
}

func (p *endpointPool) closed() bool {
	return atomic.LoadUint32(&p._closed) == 1
}

func (p *endpointPool) reaper(frequency time.Duration) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for range ticker.C {
		if p.closed() {
			break
		}
		_, err := p.ReapStaleConns()
		if err != nil {
			continue
		}
	}
}

func (p *endpointPool) reapStaleConn() *Conn {
	if len(p.idleConns) == 0 {
		return nil
	}

	cn := p.idleConns[0]
	if !p.isStaleConn(cn) {
		return nil
	}

	p.idleConns = append(p.idleConns[:0], p.idleConns[1:]...)
	p.idleConnsLen--
	p.removeConn(cn)

	return cn
}

func (p *endpointPool) isStaleConn(cn *Conn) bool {

	if cn.closed() {
		return true
	}

	if p.opt.IdleTimeout == 0 && p.opt.MaxConnAge == 0 {
		return false
	}

	now := time.Now()
	if p.opt.IdleTimeout > 0 && now.Sub(cn.UsedAt()) >= p.opt.IdleTimeout {
		return true
	}
	if p.opt.MaxConnAge > 0 && now.Sub(cn.createdAt) >= p.opt.MaxConnAge {
		return true
	}

	return false
}
