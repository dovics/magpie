package scheduler

type Storage interface {
	Save(ID string, data []byte) error
}
