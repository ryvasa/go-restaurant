package handler

type Handlers struct {
	MenuHandler   *MenuHandlerImpl
	UserHandler   *UserHandlerImpl
	ReviewHandler *ReviewHandlerImpl
	AuthHandler   *AuthHandlerImpl
	// Tambahkan handler lain di sini
	// OrderHandler *OrderHandler
	// UserHandler  *UserHandler
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	reviewHandler *ReviewHandlerImpl,
	AuthHandler *AuthHandlerImpl,
) *Handlers {
	return &Handlers{
		MenuHandler:   menuHandler,
		UserHandler:   userHandler,
		ReviewHandler: reviewHandler,
		AuthHandler:   AuthHandler,
	}
}
