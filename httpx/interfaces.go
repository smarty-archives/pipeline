package httpx

type WaitGroup interface {
	Add(delta int)
	Done()
}
