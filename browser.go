package warcbrowser

import "fmt"

type Tab struct {
	Id    int
	Url   string
	Title string
}

func (t *Tab) String() string {
	return fmt.Sprintf("[%d] %s ( %s )", t.Id, t.Url, t.Title)
}

type Browser interface {
	ListTabs(pattern string) []Tab
	ArchiveTab(tab int)     error
	ArchiveUrl(url string) error
}
