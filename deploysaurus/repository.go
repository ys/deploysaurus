package deploysaurus

import (
	"fmt"
	"strings"
)

type Repository struct {
	FullName   string `json:"full_name"`
	ArchiveUrl string `json:"archive_url"`
}

func (repo *Repository) AuthenticatedArchiveUrl(format string, ref string, authToken string) string {
	return replaceFormatAndRef(addAuthentication(repo.ArchiveUrl, authToken), format, ref)
}

func replaceFormatAndRef(url string, format string, ref string) string {
	if format == "" {
		format = "tarball"
	}
	if ref == "" {
		ref = "master"
	}
	urlWithFormat := strings.Replace(url, "{archive_format}", format, -1)
	return strings.Replace(urlWithFormat, "{/ref}", fmt.Sprintf("/%s", ref), -1)
}

func addAuthentication(url string, auth string) string {
	if auth != "" {
		return strings.Replace(url, "://", fmt.Sprintf("://%s:x-oauth-basic@", auth), -1)
	} else {
		return url
	}
}
