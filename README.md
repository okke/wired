# WIRED
Wired is a dependency injection framework for Go which supports both constructor argument injection as property injection. It has a clear and simple to use API that promotes an elegant coding style. Injection behavior is extensible. Singletons and factories come out of the box. As well as support for initialization of primitive fields from configuration (like environment variables).

# Simple example
Let's showcase Wired with a simple straightforward example. In this case, two structs (a Room and a Table) with a dependency between them are created. Both structs have functions that are registered as constructors so Wired knows how to create Rooms and Tables. Using a call to *Inject*, which takes a function with an arbitrary amount of arbitrary typed arguments, required objects are created and can be used.

```Go
import (
  "fmt"

  "github.com/okke/wired"
)

type Table struct {  

}

type Room struct {
  table *Table
}

func NewTable() *Table { return &Table{} }

func NewRoom(table *Table) *Room { return &Room{table: table} }

func Work() {

  wired.Global().Go(func(scope wired.Scope) {

    scope.Register(NewTable)
    scope.Register(NewRoom)

    scope.Inject(func(room *Room) {
      // do something with room and its table
      fmt.Println("room has a table", room.table)
    })
  })
}
```

## Auto-wire
Previous examples used argument injection to construct objects. The NewRoom function required a Table pointer which is provided by Wired. Another way of achieving the same is through the concept of auto-wiring. Simply specify in a struct definition that Wired may try to initialize it's fields. For this, the framework is using special (empty) tagging structs. In this case *wired.AutoWire*.

```Go
type Table struct {

}

type Room struct {
  wired.AutoWire      // This will drive auto-wiring

  Table *Table        // so this Table can be injected after a Room is created
}

func NewTable() *Table { return &Table{} }
func NewRoom() *Room   { return &Room{} }
```

Note, in this example the table field in the Room struct has been made public. This is needed otherwise Wired can not access the field and initialize it.

## Auto-config
Wired also supports the wiring of primitive values that can be used as initialization values. It comes with a simple string template parser that will lookup key-value pairs. This can be used to inject for example environment variables. But also from other sources like configuration servers.

```Go
type ServerConfiguration struct {

  wired.AutoConfig  // this will tell Wired to look for config tags on fields
  
  Port     int    `autoconfig:"${server_port:8080}"`    // lookup server_port or use default value 8080
  Prefix   string `autoconfig:"${server_prefix:/api/}"` // lookup server_prefix or use '/api/'

}
```

All primitive types are supported. And hooking in your own configurators is as easy as registering a Configurator like any other type:

```Go
func NewMyConfig() wired.Configurator {
  
  // return an object that implements the wired.Configurator interface
  // that'll do the variable lookup ( the 'ConfigValue(key string) string' method )
  // 
  return &myConfig{} 
}
```

```Go
wired.Global().Go(func(scope wired.Scope) {

    scope.Register(NewMyConfig)
}
``` 

## Constructing slices
Wired can construct slices of known types. When multiple registered constructors return the same type, and somewhere a slice of this type is required, all constructors are called to fill the slice.

```Go
type Listener interface {
  Listen(message Message)
}

func ConsumeAndLog() Listener {
  return .... // construct listener
}

func ConsumeAndForward() Listener {
  return .... // construct listener
}
```

```Go
wired.Global().Go(func(scope wired.Scope) {

  scope.Register(ConsumeAndLog)     // register Listener
  scope.Register(ConsumeAndForward) // register Listener

  // and use listeners
  //
  scope.Inject(func(listeners []Listener) {
    for _, listener := range listeners {
      listener.Listen(Produce())
    }
  })
})
```

## Constructing maps
Wired can construct maps of known types. When multiple registered constructors return the same type and this type has a *Key()* method defined, a map from key type to constructed type will be created. So when for example you have a driver struct with a *Key()* method returning the driver's name as a string, a map from string to driver struct will be available.

```Go
type driver struct {
	name string
}

func (driver *driver) Key() string {
  return driver.name
}

func newAWSDriver() *driver {
  return &driver{name: "aws"}
}

func newAzureDriver() *driver {
  return &driver{name: "azure"}
}

func newGoogleDriver() *driver {
  return &driver{name: "google"}
}
```

```Go
wired.Go(func(scope wired.Scope) {
  scope.Register(newAWSDriver)
  scope.Register(newAzureDriver)
  scope.Register(newGoogleDriver)

  scope.Inject(func(drivers map[string]*driver) {
    // do something with all drivers
  })
}
``` 

## Singletons
Singletons are supported by embedding a *wired.Singleton* 'tag' inside a struct. This tells Wired that within given scope, only one instance of this struct will be constructed.

```Go
type Father struct {
  wired.Singleton   // This will tell Wired that Father is a Singleton

  // rest of struct definition
}

type Son struct {
  wired.AutoWire    // This will auto-wire the Father singleton

  father *Father
}
```

## Factories
Factories are objects that create other objects. Factories are recognized through the usage of a *wired.Factory* 'tag' and should have a *Construct* method that can be used to construct objects. 

```Go
type pepperFactory struct {
  wired.Factory   // This will tell Wired to treat this as a factory

  count int       // factory state
}

type pepper struct {
  nr int
}

// Construct will construct pepper objects
// so whenever such an object is required, this method will be called
//
func (pepperFactory *pepperFactory) Construct() *pepper {
  pepperFactory.count = pepperFactory.count + 1
  return &pepper{nr: pepperFactory.count}
}
```

Example usage of a factory object:

```Go
wired.Global().Go(func(scope wired.Scope) {
  scope.Register(newPepperFactory)
  
  // instead of injecting the pepperFromFactory into a function,
  // ask Wired to construct one for us
  //
  scope.Inject(func (pepper *pepper) {
    // do something with pepper
  })
})
```

Note, factories are singletons but, in contrast to other objects, are created as soon as they are registered.

## Scopes
Wired knows the concept of scopes to hold references to constructor functions and singleton objects. Scopes are always inherited from another scope. Everything that is accessible in a parent scope, is accessible within a child scope. But everything in a child scope is not known to the parent and will override the parent.

Creating a new top level scope can be done through the *wired.Go* function:

```Go
wired.Go(func (scope wired.Scope) {

})
```

Creating a child (or inner or sub) scope can be done by using a scope's *Go* method:

```Go
wired.Go(func (scope wired.Scope) {
  
  // register things on a top level

  scope.Go(func (inner wired.Scope) {

    // so something on a sub level
  })

})
```

Wired comes with a global scope that can be used to register constructors and factories on a global level. Doing this in an init function and by always using the global scope as a starting point for object graph construction, ensures Wired knows how to create the types you need. The global scope is accessible through the *wired.Global()* function.

```Go
func init() {

  wired.Global().Register(MyConstructor)

}

func someWhereElse() {

  wired.Global().Go(func (scope wired.Scope) {
    // have fun with the instance(s) MyConstructor can create
  })
  
}
```
