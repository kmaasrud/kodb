package bib

import "time"

type Name struct {
		Family		string				`json:"family"`
		Given		string				`json:"given"`
		Literal		string				`json:"literal"`
		NonDropPart string				`json:"non-dropping-particle"`
		DropPart	string				`json:"dropping-particle"`
	}

type BibEntry struct {
	Id				    string				`json:"id"`
	Type			    string				`json:"type"`

    Abstract            string              `json:"abstract"`
	Accessed		    time.Time			`json:"accessed"`
    Archive             string              `json:"archive"`
    ArchiveLocation     string              `json:"archive_location"`
	Author			    []Name              `json:"author"`
    Authority           string              `json:"authority"`
    CallNumber          string              `json:"call-number"`
    ChapterNumber       int                 `json:"chapter-number"`
    CollectionEditor    []Name              `json:"collection-editor"`
    CollectionNumber    int                 `json:"collection-number"`
    CollectionTitle     string              `json:"collection-title"`
    Composer            []Name              `json:"composer"`
    ContainerAuthor     []Name              `json:"container-author"`
    ContainerTitle      string              `json:"container-title"`
    Dimensions          string              `json:"dimensions"`
    Director            []Name              `json:"director"`
    DOI                 string              `json:"DOI"`
    Edition             int                 `json:"edition"`
    Editor              []Name              `json:"editor"`
    Event               string              `json:"event"`
    EventPlace          string              `json:"event-place"`
    Genre               string              `json:"genre"`
    Interviewer         []Name              `json:"interviewer"`
    ISBN                string              `json:"ISBN"`
    ISSN                string              `json:"ISSN"`
    Issue               int                 `json:"issue"`
	Issued			    time.Time			`json:"issued"`
    Language            string              `json:"language"`
    Medium              string              `json:"medium"`
    Note                string              `json:"note"`
    Number              int                 `json:"number"`
    NumberOfPages       int                 `json:"number-of-pages"`
    NumberOfVolumes     int                 `json:"number-of-volumes"`
    Page                string              `json:"page"`
    Publisher           string              `json:"publisher"`
    PublisherPlace      string              `json:"publisher-place"`
    Recipient           []Name              `json:"recipient"`
    References          string              `json:"references"`
    ReviewedAuthor      []Name              `json:"reviewed-author"`
    Scale               string              `json:"scale"`
    Section             string              `json:"section"`
    Source              string              `json:"source"`
    Status              string              `json:"status"`
    Submitted           time.Time           `json:"submitted"`
	Title			    string				`json:"title"`
    TitleShort          string              `json:"title-short"`
    Translator          []Name              `json:"translator"`
    URL                 string              `json:"URL"`
    Version             string              `json:"version"`
    Volume              int                 `json:"volume"`

	Fields			    map[string]string	`json:"-"`
}
