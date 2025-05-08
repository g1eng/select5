package select5_test

import (
	"fmt"
	"github.com/g1eng/select5"
	"os"
	"strings"
	"testing"
)

func TestEditor_PutEnter_After_Down(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"middle",
			select5.Editor{
				select5.CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita
Aayok
unetato
Akubishinagara`,
		},
		{
			"base",
			select5.Editor{
				select5.CursorPosition{
					X: 6,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosono
so
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari

Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				os.Stderr,
				sl,
			},
			strings.Join(sl, "\n") + "\n",
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita

Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[3]),
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita
Aayokuneta
to
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.Down()
			tt.editor.PutEnter()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutEnter_After_Up(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"middle",
			select5.Editor{
				select5.CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumas
an
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"base",
			select5.Editor{
				select5.CursorPosition{
					X: 6,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruni
nari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`
Haruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato

Akubishinagara`,
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso

Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[3]),
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan

Okimashita
Aayokunetato
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.Up()
			tt.editor.PutEnter()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutEnter_After_Up_InputO(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"middle",
			select5.Editor{
				select5.CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumas
Oan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"base",
			select5.Editor{
				select5.CursorPosition{
					X: 6,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruni
Onari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				os.Stderr,
				sl,
			},
			`
OHaruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"document end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato
O
Akubishinagara`,
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso

OKumasan
Okimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[3]),
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
O
Okimashita
Aayokunetato
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.Up()
			tt.editor.PutEnter()
			tt.editor.PutS([]byte{'O'})
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutEnter_After_Right2_Left_Input_Dollar(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"middle",
			select5.Editor{
				select5.CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimas$Phita
AayokuneNeKotato
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.Right()
			tt.editor.Right()
			tt.editor.Left()
			tt.editor.PutS([]byte{'$'})
			tt.editor.PutS([]byte{'P'})
			tt.editor.Down()
			tt.editor.PutS([]byte{'N'})
			tt.editor.PutS([]byte{'e'})
			tt.editor.PutS([]byte{'K'})
			tt.editor.PutS([]byte{'o'})
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_CursorMove_Utf8Modification(t *testing.T) {
	sl := []string{
		"春になり",
		"のそのそ熊さん",
		"起きました。",
		"「ああよく寝た」と",
		"あくびしながら",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"core",
			select5.Editor{
				select5.CursorPosition{
					X: 12,
					Y: 3,
				},
				nil,
				os.Stderr,
				sl,
			},

			`春になり
のそのそ熊さん
起きまa〜~----
したよ。
「ああ寝た」と
あくびしながら`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutDelete()
			tt.editor.PutBackspace()
			tt.editor.Up()
			tt.editor.PutS([]byte{'a'})
			tt.editor.PutS([]byte("〜")[0:3])
			tt.editor.PutS([]byte{'~'})
			tt.editor.PutS([]byte{'-'})
			tt.editor.PutS([]byte{'-'})
			tt.editor.PutS([]byte{'-'})
			tt.editor.PutS([]byte{'-'})
			tt.editor.PutEnter()
			tt.editor.Right()
			tt.editor.Right()
			tt.editor.PutS([]byte("よ")[0:3])
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				var report string
				report += "[LINE1] \n"
				for _, b := range tt.editor.Line[1] {
					report += fmt.Sprintf("%c (%x)\n", b, b)
				}
				report += "[LINE2] \n"
				for _, b := range tt.editor.Line[2] {
					report += fmt.Sprintf("%c (%x)\n", b, b)
				}
				report += "[LINE3] \n"
				for _, b := range tt.editor.Line[3] {
					report += fmt.Sprintf("%c (%x)\n", b, b)
				}
				t.Errorf("select5.Editor.Line() = %v, want %v\n\nreport: %s", got, tt.want, report)
			}
		})
	}
}
