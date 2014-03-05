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
