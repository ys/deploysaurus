package deploysaurus

import "testing"
import "github.com/joho/godotenv"

func TestTarball(t *testing.T) {
	_ = godotenv.Load(".env.test")
	repo := Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	sender := Sender{Id: 123, Login: "ys", DbUser: &DbUser{GitHubToken: "github_deploy_key"}}
	event := Event{Sha: "sha", Repository: &repo, Sender: &sender}
	tarball := "https://github_deploy_key:x-oauth-basic@api.github.com/repos/ys/rsss/tarball/sha"
	eventTarball := event.Tarball()
	if tarball != eventTarball {
		t.Error("Expected ", tarball, " got ", eventTarball)
	}
}

func TestWhat(t *testing.T) {
	repo := Repository{FullName: "full_name"}
	event := Event{Repository: &repo}
	what := "full_name"
	eventWhat := event.What()
	if what != eventWhat {
		t.Error("Expected ", what, " got ", eventWhat)
	}
}

func TestWhatWithoutRepo(t *testing.T) {
	event := Event{}
	what := ""
	eventWhat := event.What()
	if what != eventWhat {
		t.Error("Expected ", what, " got ", eventWhat)
	}
}

func TestProcessable(t *testing.T) {
	_ = godotenv.Load(".env.test")
	repo := Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	sender := Sender{Id: 123, Login: "ys", DbUser: &DbUser{HerokuId: "AZE"}}
	event := Event{Sha: "sha", Repository: &repo, Sender: &sender}
	if _, err := event.Processable(); err != nil {
		t.Error("Expected event to be processable")
	}
}

func TestProcessableWithoutHeroku(t *testing.T) {
	_ = godotenv.Load(".env.test")
	repo := Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	sender := Sender{Id: 123, Login: "ys", DbUser: &DbUser{}}
	event := Event{Sha: "sha", Repository: &repo, Sender: &sender}
	if _, err := event.Processable(); err == nil && err.Error() != "No Heroku" {
		t.Error("Expected event to be unprocessable")
	}
}

func TestProcessableWithoutSender(t *testing.T) {
	_ = godotenv.Load(".env.test")
	repo := Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	event := Event{Sha: "sha", Repository: &repo}
	if _, err := event.Processable(); err == nil && err.Error() != "No GitHub" {
		t.Error("Expected event to be unprocessable")
	}
}

func TestSenderLogin(t *testing.T) {
	event := Event{Sender: &Sender{Login: "Bob"}}
	if event.SenderLogin() != "Bob" {
		t.Error("Expected 'Bob' got", event.SenderLogin())
	}
}

func TestSenderLoginWithoutSender(t *testing.T) {
	event := Event{}
	if event.SenderLogin() != "Somebody" {
		t.Error("Expected 'Somebody' got", event.SenderLogin())
	}
}
