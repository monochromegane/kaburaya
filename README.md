# Kaburaya [![Build Status](https://travis-ci.org/monochromegane/kaburaya.svg?branch=master)](https://travis-ci.org/monochromegane/kaburaya)

WIP.

Kaburaya optimize the number of goroutines by feedback control. It provides elastic semaphore.

## Usage

```go
sem := kaburaya.NewSem(100 * time.Millisecond)
var wg sync.WaitGroup
for // Something condition {
	wg.Add(1)
	sem.Wait()
	go func() {
		defer sem.Signal()
		defer wg.Done()
		// Something job
	}()
}
wg.Wait()
sem.Stop()
```

## License

[MIT](https://github.com/monochromegane/kaburaya/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
