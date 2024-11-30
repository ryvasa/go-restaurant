package handler

type Handlers struct {
	MenuHandler        *MenuHandlerImpl
	UserHandler        *UserHandlerImpl
	ReviewHandler      *ReviewHandlerImpl
	AuthHandler        *AuthHandlerImpl
	OrderHandler       *OrderHandlerImpl
	TableHandler       *TableHandlerImpl
	ReservationHandler *ReservationHandlerImpl
	RecipeHandler      *RecipeHandlerImpl
	InventoryHandler   *InventoryHandlerImpl
	IngredientHandler  *IngredientHandlerImpl
}

func NewHandlers(
	menuHandler *MenuHandlerImpl,
	userHandler *UserHandlerImpl,
	reviewHandler *ReviewHandlerImpl,
	authHandler *AuthHandlerImpl,
	orderHandler *OrderHandlerImpl,
	tableHandler *TableHandlerImpl,
	reservationHandler *ReservationHandlerImpl,
	recipeHandler *RecipeHandlerImpl,
	inventoryHandler *InventoryHandlerImpl,
	ingIngredientHandler *IngredientHandlerImpl,

) *Handlers {
	return &Handlers{
		MenuHandler:        menuHandler,
		UserHandler:        userHandler,
		ReviewHandler:      reviewHandler,
		AuthHandler:        authHandler,
		OrderHandler:       orderHandler,
		TableHandler:       tableHandler,
		ReservationHandler: reservationHandler,
		RecipeHandler:      recipeHandler,
		InventoryHandler:   inventoryHandler,
		IngredientHandler:  ingIngredientHandler,
	}
}
