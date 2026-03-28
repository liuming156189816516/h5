package baselib

var g_arrConfigFunc []func() error = nil

func init() {
	g_arrConfigFunc = make([]func() error, 0, 2)
	sig := NewSignalHandler()
	go func() {
		for {
			select {
			case <-sig.ReloadSignal():
				ReloadServerConfig(nil)
			}
		}
	}()
}
func ReloadServerConfig([]byte) {
	for _, fb := range g_arrConfigFunc {
		fb()
	}
}
func InitConfig(LoadConfig func() error) error {
	if LoadConfig == nil {
		return nil
	}
	if err := LoadConfig(); err != nil {
		return err
	}

	RegisterReloadFunc(LoadConfig)
	return nil
}

func RegisterReloadFunc(f func() error) {
	g_arrConfigFunc = append(g_arrConfigFunc, f)
}
