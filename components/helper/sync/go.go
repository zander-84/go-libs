package sync

func Go(fn func()) {
	go func() {
		defer func() {
			defer func() {
				if err := recover(); err != nil {
					return
				}
			}()
		}()
		fn()
	}()
}
