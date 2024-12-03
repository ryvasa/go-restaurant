package handler

type Handlers struct {
	MenuHandler        MenuHandler
	UserHandler        UserHandler
	ReviewHandler      ReviewHandler
	AuthHandler        AuthHandler
	OrderHandler       OrderHandler
	TableHandler       TableHandler
	ReservationHandler ReservationHandler
	RecipeHandler      RecipeHandler
	InventoryHandler   InventoryHandler
	IngredientHandler  IngredientHandler
}

func NewHandlers(
	menuHandler MenuHandler,
	userHandler UserHandler,
	reviewHandler ReviewHandler,
	authHandler AuthHandler,
	orderHandler OrderHandler,
	tableHandler TableHandler,
	reservationHandler ReservationHandler,
	recipeHandler RecipeHandler,
	inventoryHandler InventoryHandler,
	ingIngredientHandler IngredientHandler,

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
