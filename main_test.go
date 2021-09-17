package main

import (
	"strings"
	"testing"
)

func TestOnepifToCSV(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name: "",
			input: `{"uuid":"7890D22A12584A178484BF175E724E09","updatedAt":1624909782,"locationKey":"test.com","securityLevel":"SL5","contentsHash":"babf5bcb","title":"Happy","location":"https:\/\/sso.test.com/path/to/login?query","secureContents":{"URLs":[{"url":"https:\/\/sso.test.com"}],"fields":[{"id":"txtBoxEmail","value":"me@email.com","name":"emailAddr","type":"E","designation":"username"},{"id":"txtBoxPasswordConfirm","value":"pass","name":"newPasswordConfirm","type":"P","designation":"password"},{"id":"signUp-chk1","value":"✓","name":"textAgreementFlag","type":"C","designation":""},{"id":"signUp-chk2","value":"","name":"elecMarketingFlag","type":"C","designation":""}],"sections":[{"fields":[{"k":"string","n":"7fgasupcejqfx1pfwvgnipwdjy","v":"First","t":"Vorname*"},{"k":"string","n":"zvhjpv434jx6tbsaib2w4nnfdu","v":"LAst","t":"Nachname*"},{"k":"string","n":"12buh2ler6jctqucm3pld7s5fi","v":"19860801","t":"Geburtsdatum*"},{"k":"string","n":"x6ipyc6yhkmj2c4qhgh5ffgaas","v":"KP9C","t":"Bitte geben Sie den im Bild gezeigten Text ein oder hören Sie den Ton*"}],"title":"Saved on sso.test.com","name":"Section_tw6baovoe2bkzyh3dryjyfbgfa"}]},"txTimestamp":1624909782,"createdAt":1624909779,"typeName":"webforms.WebForm"}
***1BB0BAC0-DAA3-43BC-B825-C5A3526CA580***
{"uuid":"F72C569BEC624A91A58232EB31B8C112","updatedAt":1628374191,"securityLevel":"SL5","contentsHash":"40e33862","title":"noURL","secureContents":{"fields":[{"value":"me@email.com","name":"username","type":"T","designation":"username"},{"value":"USDFISDN2342§§","name":"password","type":"P","designation":"password"}],"passwordHistory":[{"value":"foo","time":1628374191}]},"txTimestamp":1628424597,"createdAt":1628373881,"typeName":"webforms.WebForm"}
***1BB0BAC0-DAA3-43BC-B825-C5A3526CA580***
{"uuid":"0D0091B7D4FB408BBD3263A631B49A86","updatedAt":1620074847,"locationKey":"foo.com","securityLevel":"SL5","contentsHash":"f32561ab","title":"WithOTP","location":"https:\/\/foo.com","secureContents":{"fields":[{"value":"pass","name":"password","type":"P"},{"value":"me","name":"username","type":"T","designation":"username"},{"id":"user_password_confirmation","value":"foo","name":"user[password_confirmation]","type":"P","designation":"password"}],"notesPlain":"note1\nnote2","sections":[{"fields":[{"k":"concealed","n":"TOTP_1234563E60774A679FAEB79CE942982B","v":"otpauth:\/\/totp\/foo.com:foo.com_me%40email.com?secret=123VS7BXMJBVPDMCUGYXDRJ3NSGXSE99&issuer=foo.com","t":"Einmal-Passwort"}],"name":"Section_921339DB42F549A8BD902EAD5D64075A"},{"title":"Verwandte Objekte","name":"linked items"}],"URLs":[{"url":"https:\/\/foo.com"}]},"txTimestamp":1620074847,"createdAt":1620074578,"typeName":"webforms.WebForm"}
***1BB0BAC0-DAA3-43BC-B825-C5A3526CA580***
`,
			want: `Title,Url,Username,Password,OTPAuth
Happy,https://sso.test.com/,me@email.com,pass,
WithOTP,https://foo.com/,me,foo,apple-otpauth://totp/foo.com:foo.com_me%40email.com?secret=123VS7BXMJBVPDMCUGYXDRJ3NSGXSE99&issuer=foo.com
`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			err := onepifToCSV(strings.NewReader(tc.input), &b)
			got := b.String()
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
			if (err != nil) != tc.wantErr {
				if tc.wantErr {
					t.Errorf("want error, but didn't get one")
				} else {
					t.Errorf("want no error, but got %v", err)
				}
			}
		})
	}
}
