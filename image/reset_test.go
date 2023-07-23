package main

import (
	"testing"

	"github.com/joker-circus/gopkg/internal"
)

const _photo = `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACgAAAAVCAYAAAA0GqweAAAB8klEQVR42mNgGAWjYBCA0CMi6kFHhLYA8Rcgfh10WGh50CERSWLlsQG/IyJSQLXHgNgCxa6D/IpBRwXXA8XfA/GToMPCnQ77GVjwOhBs+VGh+sCjgjqBR4WdgPwrQHyEWHkUB5wR5Ad6oDj4iNA7oJr/6A4EmtMTfFgoKfigqCrQvAig/DcgrsXrQL8TYuLI/JCjQm4gw33384oQI4/uAKDcw8AjQiXYHJiwX4EDNXCElwDVHiIpyoOOCtiBDA89xCNKSD5gv4AAyEFAnAKWOyEk0/CfgQlEY3MgOgA6biUwtNeS5sAjwkuBBl8gRj7tDANr4FGhuUBsiaKGgAO9D/MLBh0VTgClxeDDwubEZZZVDMxADf1A/B2bJkLyxDoQmtlAcl8DjwhmEOU4UDoLOiK4D6jpGTbLCcmT4kBQDg89JKQFzMGxwAz1GBQDeA3zPygkCzToHhBvx5buCMmT6kBkEHxUKBikDj0joibUo0J7gOlqEyhxkyNPigM9tzGwo6g7LBwIznD7RSVwOE5cDGzQYcGsgKMCBsgYJEdIHpxJgDkRlLMJORBUIAOj9GTwMSH3gKP8ysAy1R8ofx+YdDbiDuJDgnpQgzBw8BHBCkLysGIm+IhwPDEhCCkfBW9DC+iboArA94wU12h9PgpGwXAFAFtCUmgSKzHlAAAAAElFTkSuQmCC`

var _img, _ = internal.Base64ToImage(internal.S2b(internal.Base64TrimData(_photo)))

func BenchmarkResetImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ResetImage(_img, _img.Bounds().Dx()*3, _img.Bounds().Dy()*3)
	}
}

func BenchmarkResetImage2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ResetImage2(_img, _img.Bounds().Dx()*3, _img.Bounds().Dy()*3)
	}
}
