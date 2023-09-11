package httpresp

const (
	DefaultLang = "en"
)

type Header struct {
	LanguageCode *string `header:"languageCode" json:"-"`
	OsType       *string `header:"osType" json:"-"`
	AppVersion   *string `header:"appVersion" json:"-"`
}

func (header *Header) GetLanguageCode() string {
	if header == nil {
		return DefaultLang
	}
	if header.LanguageCode == nil {
		return DefaultLang
	}

	return *header.LanguageCode
}
