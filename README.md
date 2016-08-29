# Thunder
Golang toolkit for the [ZigZag](https://github.com/valerykalashnikov/thunder) storage API

##Usage

~~~go
  import "github.com/valerykalashnikov/thunder"
~~~

Construct a new Thunder client:

~~~go
  client := thunder.Client{"http://localhost:8082", "password"}
~~~

Invoke methods on the client to access different parts of the ZigZag API

To write value to the storage

~~~go
  value := "value" // also you can use arrays and maps as values
  err := client.Set("key", value) // error can contain http status code and message
  
  // or if you want to set expiration value (in minutes)
  
  value = "value will expire in 1 minute"
  err := client.Set("key_with_1_minute_expiration", value, 1)
~~~

To obtain value from storage

~~~go
  value, err := client.Get("key") // error can contain http status code and message
~~~

To update storage value (usefull when you have to update value without breaking expiration)

~~~go
  value := "New value"
  err := client.Update("key", value) // error can contain http status code and message
~~~

To delete value from storage

~~~go
  client.Delete("key")
~~~

To obtain keys matching pattern

~~~go
  keys, err := client.Keys("^[a-z]")
~~~

