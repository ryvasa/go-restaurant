package repository

type TransactionRepository interface {
	Transact(txFunc func(adapters Adapters) error) error
}
