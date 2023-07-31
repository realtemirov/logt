# Logt
Color logger for golang


# Usage


```go
func main() {
    // Create a new logger
    // nameSpace: The name of the logger
    // true: Enable show intro message
    l := logt.NewLogData("repository",true)

    // Initialize writer logger
    w := l.NewWriter("item.Create()")
    defer w.Close()

    // Write a message
    w.Data("Hello world")
    w.Debug("Hello world")
    w.Error("Hello world")
    w.Info("Hello world")
    w.Msg("Hello world")
    w.Succes("Hello world")
    w.Warning("Hello world")
    w.Write("Hello world")

}
```

# Output
![Alt text](image.png)

# TODO

- Add writers with format
- Add writers to file
- ...