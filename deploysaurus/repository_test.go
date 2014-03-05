package deploysaurus

import "testing"

func TestFormat(t *testing.T) {
	var repo Repository
	repo = Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	authenticatedUrl := "https://123:x-oauth-basic@api.github.com/repos/ys/rsss/zipball/abcdef"
	repoAuth := repo.AuthenticatedArchiveUrl("zipball", "abcdef", "123")
	if repoAuth != authenticatedUrl {
		t.Error("Expected ", authenticatedUrl, " got ", repoAuth)
	}
}

func TestDefaults(t *testing.T) {
	var repo Repository
	repo = Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}"}
	authenticatedUrl := "https://123:x-oauth-basic@api.github.com/repos/ys/rsss/tarball/master"
	repoAuth := repo.AuthenticatedArchiveUrl("", "", "123")
	if repoAuth != authenticatedUrl {
		t.Error("Expected ", authenticatedUrl, " got ", repoAuth)
	}
}
