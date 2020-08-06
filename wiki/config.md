# Configuration Map

It is good practice to use constants in your code instead of [magic numbers](https://en.wikipedia.org/wiki/Magic_number_(programming))
and it's even better to set them outside your code either by providing a config file or reading from Environment variable.
Mortar have a `Config` interface that is used everywhere to read external configurations.
While Mortar can be configured explicitly and that gives you total control over it. It is much comfortable to use its defaults.
To read them Mortar expects a dedicated Configuration key called **mortar**

```yaml
mortar:
  name: "tutorial"
  server:
    grpc:
      port: 5380
    rest:
      external:
        port: 5381
      internal:
        port: 5382
...
```

Every project should have a configuration file, and your is no exception. You can put all your external configuration values
in it (or them).

The concept is simple, you use the `Get` function that accepts a key. A key is actually a path within the configuration map.
Looking at the above example to access gRPC server port you should use the following key.

`mortar.server.grpc.port`

> Default delimiter is `.` but if needed it can be changed. (TODO)

Once you `Get` a value with a provided key you can

- Check if it was set `value.IsSet() bool`
- Cast it to a type
  - `Bool() bool`
  - `Int() int`
  - `StringMapStringSlice() map[string][]string`
  - ...

## Environment variables

While it depends on the implementation you should assume that if there is an Environment Variable with a *matching* name
its value going to be used first.

### *Matching* Environment Variable names

As mentioned previously default delimiter is `.` however when naming Environment variables you can't use `.`.
It is expected, that chosen Implementation will allow you to configure a delimiter *replacer*.
If you choose to use [viper](https://github.com/spf13/viper) using [brick wrapper](https://github.com/go-masonry/bviper).
By default, there is a replacer that will replace `_` delimiter to `.` used in our code.

This is better explained with an example. Look at the map below.

```yaml
mortar:
  server:
    grpc:
      port: 5380
```

Let's say you want to change port value from 5380 to 7777, you can change the file itself. However, you can also override it.
Viper allows you to override configuration values with a matching Environment Variable. In our case:

```shell script
export MORTAR_SERVER_GRPC_PORT="7777"
```

When you want to read gRPC server port value in your code you should write

`mortar.server.grpc.port`

Viper will look for Environment variable by replacing `_` with `.` case-insensitive **first** and return its value if set.

## Mortar Keys

Mortar expects different keys in it's configuration map to enable or configure different abilities.
In [this](../mortar/keys.go) file we try to show what the configuration map should look like and expose all the keys.

## Config format

While in this example we showed you `config.yml` in YAML format you can choose what ever works for you as long as the provided `Config` implementation knows how to read it and will abstract every key to be queried in this form:

`root.child.childOfChild`
