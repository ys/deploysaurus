package deploysaurus

import "testing"

func TestFormat(t *testing.T) {
	var repo Repository
	repo = Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}",
		AuthToken: "123"}
	authenticatedUrl := "https://123@api.github.com/repos/ys/rsss/zipball/abcdef"
	repoAuth := repo.AuthenticatedArchiveUrl("zipball", "abcdef")
	if repoAuth != authenticatedUrl {
		t.Error("Expected ", authenticatedUrl, " got ", repoAuth)
	}
}

func TestDefaults(t *testing.T) {
	var repo Repository
	repo = Repository{ArchiveUrl: "https://api.github.com/repos/ys/rsss/{archive_format}{/ref}",
		AuthToken: "123"}
	authenticatedUrl := "https://123@api.github.com/repos/ys/rsss/tarball/master"
	repoAuth := repo.AuthenticatedArchiveUrl("", "")
	if repoAuth != authenticatedUrl {
		t.Error("Expected ", authenticatedUrl, " got ", repoAuth)
	}
}
