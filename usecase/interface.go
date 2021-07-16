package usecase

type UseCaseErr interface {
	error
	HttpStatus() int
}
