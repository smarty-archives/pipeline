package httpx

type WaitGroup interface {
	Add(delta int)
	Wait()
	Done()
}
