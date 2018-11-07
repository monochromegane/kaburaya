package kaburaya

type elasticChannel struct {
	receive <-chan struct{}
	send    chan<- struct{}
	queue   *queue
}
