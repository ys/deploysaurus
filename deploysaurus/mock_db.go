package deploysaurus

type MockDB struct {
}

func (db MockDB) SaveUser(user DbUser) (string, error) {
	return "123", nil
}

func (db MockDB) CreateUser(user DbUser) (string, error) {
	return "123", nil
}

func (db MockDB) UpdateUser(user DbUser) (string, error) {
	return "123", nil
}

func (db MockDB) GetUser(id string) (DbUser, error) {
	return DbUser{GitHubToken: "github_deploy_key"}, nil
}

func (db MockDB) GetUserFromProvider(provider string, userId string) (DbUser, error) {
	return DbUser{GitHubToken: "github_deploy_key"}, nil
}

func (db MockDB) GetUsersCount() (int, error) {
	return 1, nil
}
