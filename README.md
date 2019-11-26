# umutex

[![Build Status](https://github.com/kawasin73/umutex/workflows/Go/badge.svg)](https://github.com/kawasin73/umutex/actions)

Upgradable RWMutex implementation for Go.

`Upgrade()` is given priority to `Lock()`.

## How to use

```go
import "github.com/kawasin73/umutex"

var mu umutex.UMutex

// goroutine 1
mu.RLock()

// goroutine 2
mu.RLock()

// goroutine 1 : Upgrade waits until RUnlock at goroutine 2
mu.Upgrade()

// goroutine 2 : Upgrade at several goroutine makes dead-locks and return false
mu.Upgrade() == false
mu.RUnlock()

// goroutine 1 : Upgrade succeeds
```

```go
import "github.com/kawasin73/umutex"

var mu umutex.UMutex

// goroutine 1
mu.RLock()

// goroutine 2 : Lock waits until unlocked at goroutine 1
mu.Lock()

// goroutine 1 : Upgrade is given priority to Lock and success
mu.Upgrade() == true

// goroutine 1
mu.Unlock()

// goroutine 2 : Lock succeeds
```

## API

- `RLock()`
- `RUnlock()`
- `Lock()`
- `Unlock()`
- `Upgrade() -> bool`

## LICENSE

MIT
