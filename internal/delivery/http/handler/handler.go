package handler

type Handlers struct {
	MenuHandler   *MenuHandlerImpl
	UserHandler   *UserHandlerImpl
	ReviewHandler *ReviewHandlerImpl
	AuthHandler   *AuthHandlerImpl
	OrderHandler  *OrderHandlerImpl
	// Tambahkan handler lain di sini
	// OrderHandler *OrderHandler
	// UserHandler  *UserHandler
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	reviewHandler *ReviewHandlerImpl,
	authHandler *AuthHandlerImpl,
	orderHandler *OrderHandlerImpl,
) *Handlers {
	return &Handlers{
		MenuHandler:   menuHandler,
		UserHandler:   userHandler,
		ReviewHandler: reviewHandler,
		AuthHandler:   authHandler,
		OrderHandler:  orderHandler,
	}
}
