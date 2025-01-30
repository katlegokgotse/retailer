// interface_adapters/repositories/user_repository.go
package repositories

import "main.go/entities"

type UserRepository interface {
	Save(user entities.User)
	FindByUsername(username string) (entities.User, bool)
}

type InMemoryUserRepository struct {
	users map[string]entities.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]entities.User),
	}
}

func (r *InMemoryUserRepository) Save(user entities.User) {
	r.users[user.Username] = user
}

func (r *InMemoryUserRepository) FindByUsername(username string) (entities.User, bool) {
	user, exists := r.users[username]
	return user, exists
}
