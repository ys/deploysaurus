package deploysaurus

import "testing"
import "github.com/joho/godotenv"

func TestTarball(t *testing.T) {
	_ = godotenv.Load(".env.test")
	repo := Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	event := Event{Sha: "sha", Repository: &repo}
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

func TestWho(t *testing.T) {
	user := User{Login: "login"}
	event := Event{Sender: &user}
	who := "login"
	eventWho := event.Who()
	if who != eventWho {
		t.Error("Expected ", who, " got ", eventWho)
	}
}

func TestWhoWithoutRepo(t *testing.T) {
	event := Event{}
	who := ""
	eventWho := event.Who()
	if who != eventWho {
		t.Error("Expected ", who, " got ", eventWho)
	}
}
