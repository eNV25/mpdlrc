package internal

type noopMutex struct{}

func (*noopMutex) Lock()    {}
func (*noopMutex) Unlock()  {}
func (*noopMutex) RLock()   {}
func (*noopMutex) RUnlock() {}
