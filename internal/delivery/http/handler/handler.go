package handler

type Handlers struct {
	MenuHandler *MenuHandlerImpl
	UserHandler *UserHandlerImpl
	// Tambahkan handler lain di sini
	// OrderHandler *OrderHandler
	// UserHandler  *UserHandler
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	// Tambahkan parameter handler lain
	// orderHandler *OrderHandler,
	// userHandler *UserHandler,
) *Handlers {
	return &Handlers{
		MenuHandler: menuHandler,
		UserHandler: userHandler,
		// OrderHandler: orderHandler,
	}
}
