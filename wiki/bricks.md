# Bricks

When we build software most of us always try to rely on something we or others build before us.
Since we don't want to reinvent the wheel, again. Go standard library is an excellent example here.
There are `strings`, `time`, `http` and many other build-in libraries that we use.
While this example is great, it doesn't scale to 3rd party libraries.

Let's look at Logger libraries for example, there are

- [Logrus](https://github.com/sirupsen/logrus)
- [Apex](https://github.com/apex/log)
- [Zap](https://github.com/uber-go/zap)
- [Zerolog](https://github.com/rs/zerolog)
- ...

You are encouraged to have a look at each library, but if not I can assure you, all of their API are different from each other.
Mortar after all is a library, and it's purpose to be used in many projects.
That is why we defined different [Interfaces](../interfaces).

## We define *Brick* as an implementation of an Interface defined in Mortar using *external* library. 

- [Mortar Logger](https://github.com/go-masonry/mortar/blob/master/interfaces/log/interfaces.go) 
    - [Implementation using Zerolog](https://github.com/go-masonry/bzerolog)
- [Mortar Config](https://github.com/go-masonry/mortar/blob/master/interfaces/cfg/interfaces.go)
    - [Implementation using Viper](https://github.com/go-masonry/bviper)
- [Mortar Tracing](https://github.com/go-masonry/mortar/blob/master/interfaces/trace/interfaces.go)
    > This is a special case, since [open tracing](https://github.com/opentracing/opentracing-go) is already an abstraction. 
    - [Implementation using Jaeger](https://github.com/go-masonry/bjaeger)
- ...

To easily differentiate an actual library from its Brick wrapper, every Brick package starts with a `b`.

- **b**viper
- **b**zerolog
- ...



