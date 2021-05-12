package test

// a place to declare test shared resources

import "sync"

// lock - the Mutex for multi-thread / go-routine resource locking
var lock sync.Mutex
