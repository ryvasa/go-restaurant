package handler

type Handlers struct {
	MenuHandler   *MenuHandlerImpl
	UserHandler   *UserHandlerImpl
	ReviewHandler *ReviewHandlerImpl
	// Tambahkan handler lain di sini
	// OrderHandler *OrderHandler
	// UserHandler  *UserHandler
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	reviewHandler *ReviewHandlerImpl,
	// Tambahkan parameter handler lain
	// orderHandler *OrderHandler,
	// userHandler *UserHandler,
) *Handlers {
	return &Handlers{
		MenuHandler:   menuHandler,
		UserHandler:   userHandler,
		ReviewHandler: reviewHandler,
		// OrderHandler: orderHandler,
	}
}
