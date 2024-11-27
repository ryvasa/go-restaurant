package handler

type Handlers struct {
	MenuHandler        *MenuHandlerImpl
	UserHandler        *UserHandlerImpl
	ReviewHandler      *ReviewHandlerImpl
	AuthHandler        *AuthHandlerImpl
	OrderHandler       *OrderHandlerImpl
	TableHandler       *TableHandlerImpl
	ReservationHandler *ReservationHandlerImpl
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	reviewHandler *ReviewHandlerImpl,
	authHandler *AuthHandlerImpl,
	orderHandler *OrderHandlerImpl,
	tableHandler *TableHandlerImpl,
	reservationHandler *ReservationHandlerImpl,
) *Handlers {
	return &Handlers{
		MenuHandler:        menuHandler,
		UserHandler:        userHandler,
		ReviewHandler:      reviewHandler,
		AuthHandler:        authHandler,
		OrderHandler:       orderHandler,
		TableHandler:       tableHandler,
		ReservationHandler: reservationHandler,
	}
}
