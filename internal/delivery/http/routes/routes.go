package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func NewRoutes(public, protected *mux.Router, handlers *handler.Handlers) {

	MenuRoutes(public, protected, handlers.MenuHandler)
	UserRoutes(public, protected, handlers.UserHandler)
	ReviewRoutes(public, protected, handlers.ReviewHandler)
	AuthRoutes(public, handlers.AuthHandler)
	OrderRoutes(protected, handlers.OrderHandler)
	TableRoutes(public, protected, handlers.TableHandler)
	ReservationRoutes(public, protected, handlers.ReservationHandler)
	RecipeRoutes(protected, handlers.RecipeHandler)
	InventoryRoutes(protected, handlers.InventoryHandler)
	IngredientRoutes(protected, handlers.IngredientHandler)

}
