Handling errors in goroutines:

1. Using errorChannel


```
 func myGoroutine(done chan<- error) {
    // Perform some work
    // If an error occurs:
    done <- myError
}

func main() {
    done := make(chan error)
    go myGoroutine(done)

    // Wait for the goroutine to complete
    err := <-done
    if err != nil {
        // Handle the error
    }
}
 
 ```

2. Using sharedError

```

var wg sync.WaitGroup
var sharedError error
var mutex sync.Mutex

func myGoroutine() {
    defer wg.Done()

    // Perform some work
    // If an error occurs:
    mutex.Lock()
    sharedError = myError
    mutex.Unlock()
}

func main() {
    // Start goroutines
    wg.Add(1)
    go myGoroutine()

    // Wait for all goroutines to complete
    wg.Wait()

    // Check for errors
    mutex.Lock()
    defer mutex.Unlock()
    if sharedError != nil {
        // Handle the error
    }
}

```
3. Using callback functions.

``` 
 type ErrorHandler interface {
    HandleError(err error)
}

func myGoroutine(errorHandler ErrorHandler) {
    defer wg.Done()

    // Perform some work
    // If an error occurs:
    errorHandler.HandleError(myError)
}

type MyErrorHandler struct {
    // Implement the HandleError method
}

func (eh *MyErrorHandler) HandleError(err error) {
    // Handle the error
}

func main() {
    errorHandler := &MyErrorHandler{}
    wg.Add(1)
    go myGoroutine(errorHandler)

    // Wait for all goroutines to complete
    wg.Wait()
}
 ```