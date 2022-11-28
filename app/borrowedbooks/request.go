package borrowedbooks

type BorrowRequest struct {
	UserID uint
	BookID uint
	Period uint
}

type ExtendBorrowRequest struct {
	Id     uint
	Period uint
}

type UserIdRequest struct {
	Id uint
}
