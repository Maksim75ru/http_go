package storages

type Employee struct {
	Id         string
	Name       string
	Sex        string
	Age        int
	Salary     int
	Department string
}

type Storage interface {
	Create(e Employee) Employee
	Get(id string) (Employee, error)
	GetAll() []Employee
	Update(id string, e Employee) (bool, error)
	Delete(id string) (bool, error)
}
