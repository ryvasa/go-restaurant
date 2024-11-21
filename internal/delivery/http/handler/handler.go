package handler

type Handlers struct {
    MenuHandler *MenuHandler
    // Tambahkan handler lain di sini
    // OrderHandler *OrderHandler
    // UserHandler  *UserHandler
}

func NewHandlers(
    menuHandler *MenuHandler,
    // Tambahkan parameter handler lain
    // orderHandler *OrderHandler,
    // userHandler *UserHandler,
) *Handlers {
    return &Handlers{
        MenuHandler: menuHandler,
        // OrderHandler: orderHandler,
        // UserHandler:  userHandler,
    }
}
