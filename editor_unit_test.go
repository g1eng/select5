package select5_test

import (
	"github.com/g1eng/select5"
	"os"
	"strings"
	"testing"
)

func TestEditor_PutEnter_NoDataCorruption(t *testing.T) {
	sl := []string{
		"\x1b[01;07mHaruninari\x1b[01;00m",
		"\x1b[01;31mNosonoso\x1b[01;00m",
		"\x1b[01;32mKumasan\x1b[01;00m",
		"\x1b[01;33mOkimashita\x1b[01;00m",
		"\x1b[01;34mAayokunetato\x1b[01;00m",
		"\x1b[01;35mAkubishinagara\x1b[01;00m",
	}
	sl1, sl2, sl3, sl4, sl5 :=
		make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6)
	copy(sl1, sl)
	copy(sl2, sl)
	copy(sl3, sl)
	copy(sl4, sl)
	copy(sl5, sl)
	r1, _, err1 := os.Pipe()
	r2, _, err2 := os.Pipe()
	r3, _, err3 := os.Pipe()
	r4, _, err4 := os.Pipe()
	r5, _, err5 := os.Pipe()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("Error creating pipe")
	}

	sc := strings.Join(sl1, "")
	tests := []struct {
		name   string
		editor select5.Editor
		want   string
	}{
		{
			"base",
			select5.Editor{
				select5.CursorPosition{
					X: 12,
					Y: 3,
				},
				nil,
				r1,
				sl,
			},
			sc,
		},
		{
			"document head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r2,
				sl2,
			},
			sc,
		},
		{
			"document end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[len(sl)-1]) - 1,
					Y: len(sl) - 1,
				},
				nil,
				r3,
				sl3,
			},
			sc,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[3]) - 1,
					Y: 3,
				},
				nil,
				r4,
				sl4,
			},
			sc,
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r5,
				sl5,
			},
			sc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutEnter()
			if got := strings.Join(tt.editor.Line, ""); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutEnter_ValidSplitting(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}
	r1, _, err1 := os.Pipe()
	r2, _, err2 := os.Pipe()
	r3, _, err3 := os.Pipe()
	r4, _, err4 := os.Pipe()
	r5, _, err5 := os.Pipe()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("Error creating pipe")
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
				r1,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okima
shita
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
				r2,
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
				r3,
				sl,
			},
			"\n" + strings.Join(sl, "\n"),
		},
		{
			"document end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				r4,
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
				r5,
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
				r2,
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
			tt.editor.PutEnter()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutCharA(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}
	sl1, sl2, sl3, sl4, sl5 :=
		make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6)
	copy(sl1, sl)
	copy(sl2, sl)
	copy(sl3, sl)
	copy(sl4, sl)
	copy(sl5, sl)
	r1, _, err1 := os.Pipe()
	r2, _, err2 := os.Pipe()
	r3, _, err3 := os.Pipe()
	r4, _, err4 := os.Pipe()
	r5, _, err5 := os.Pipe()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("Error creating pipe")
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
				r5,
				sl1,
			},
			`Haruninari
Nosonoso
Kumasan
OkimaAshita
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
				r4,
				sl3,
			},
			`AHaruninari
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
				r3,
				sl4,
			},
			strings.Join(sl, "\n") + "A",
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r2,
				sl5,
			},
			`Haruninari
Nosonoso
Kumasan
AOkimashita
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
				r1,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
OkimashitaA
Aayokunetato
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutS([]byte{'A'})
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutDelete(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}
	sl1, sl2, sl3, sl4, sl5 :=
		make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6)
	copy(sl1, sl)
	copy(sl2, sl)
	copy(sl3, sl)
	copy(sl4, sl)
	copy(sl5, sl)
	r1, _, err1 := os.Pipe()
	r2, _, err2 := os.Pipe()
	r3, _, err3 := os.Pipe()
	r4, _, err4 := os.Pipe()
	r5, _, err5 := os.Pipe()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("Error creating pipe")
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
				r1,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimahita
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
				r5,
				sl1,
			},
			`Haruniari
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
				r4,
				sl2,
			},
			`aruninari
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
					X: len(sl3[len(sl3)-1]),
					Y: len(sl3) - 1,
				},
				nil,
				r3,
				sl3,
			},
			strings.Join(sl3, "\n"),
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r2,
				sl4,
			},
			`Haruninari
Nosonoso
Kumasan
kimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl5[3]),
					Y: 3,
				},
				nil,
				r1,
				sl5,
			},
			strings.Join(sl5, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutDelete()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_PutBackspace(t *testing.T) {
	sl := []string{
		"Haruninari",
		"Nosonoso",
		"Kumasan",
		"Okimashita",
		"Aayokunetato",
		"Akubishinagara",
	}
	sl1, sl2, sl3, sl4, sl5 :=
		make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6), make([]string, 6)
	copy(sl1, sl)
	copy(sl2, sl)
	copy(sl3, sl)
	copy(sl4, sl)
	copy(sl5, sl)

	r1, _, err1 := os.Pipe()
	r2, _, err2 := os.Pipe()
	r3, _, err3 := os.Pipe()
	r4, _, err4 := os.Pipe()
	r5, _, err5 := os.Pipe()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		t.Errorf("Error creating pipe")
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
				r1,
				sl,
			},
			`Haruninari
Nosonoso
Kumasan
Okimshita
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
				r2,
				sl2,
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
					X: len(sl3[len(sl3)-1]),
					Y: len(sl3) - 1,
				},
				nil,
				r3,
				sl3,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagar`,
		},
		{
			"line head",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r4,
				sl4,
			},
			`Haruninari
Nosonoso
KumasanOkimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			select5.Editor{
				select5.CursorPosition{
					X: len(sl5[3]),
					Y: 3,
				},
				nil,
				r5,
				sl5,
			},
			`Haruninari
Nosonoso
Kumasan
Okimashit
Aayokunetato
Akubishinagara`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutBackspace()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_Utf8Input(t *testing.T) {
	sl := []string{
		"春になり",
		"のそのそ熊さん",
		"起きました。",
		"「ああよく寝た」と",
		"あくびしながら",
	}
	sl1, sl2 :=
		make([]string, len(sl)), make([]string, len(sl))
	copy(sl1, sl)
	copy(sl2, sl)

	tests := []struct {
		name   string
		editor select5.Editor
		key    select5.KeyEvent
		want   string
	}{
		{
			"core",
			select5.Editor{
				select5.CursorPosition{
					X: 12,
					Y: 2,
				},
				nil,
				os.Stderr,
				sl1,
			},
			select5.KeyEvent{
				Key:         'a',
				Code:        'a',
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: false,
				Runes:       []byte{'a'},
			},

			`春になり
のそのそ熊さん
起きましaた。
「ああよく寝た」と
あくびしながら`,
		},
		{
			"neko",
			select5.Editor{
				select5.CursorPosition{
					X: 12,
					Y: 2,
				},
				nil,
				os.Stderr,
				sl2,
			},
			select5.KeyEvent{
				Key:         0xe7,
				Code:        0xe7,
				Ctrl:        false,
				Alt:         false,
				Shift:       false,
				Special:     0,
				IsRuneStart: true,
				Runes:       []byte{0xe7, 0x8c, 0xab, 0, 0, 0},
			},

			`春になり
のそのそ熊さん
起きまし猫た。
「ああよく寝た」と
あくびしながら`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := tt.key.Utf8Char()
			if err != nil {
				t.Errorf("select5.Editor.Utf8Char() error = %v", err)
			}
			tt.editor.PutS(k)
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("select5.Editor.Line2,12() = %x%x%x, want %x%x%x", tt.editor.Line[2][12], tt.editor.Line[2][13], tt.editor.Line[2][14], k[0], k[1], k[2])
				t.Errorf("select5.Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_LineVisibleLength(t *testing.T) {
	sl := []string{
		"熊",
		"",
		"春になり",
		"のそのそ熊さん",
		"起きました.",
		"「ああよく寝た」と",
		"あくびしながら",
	}

	tests := []struct {
		name   string
		editor select5.Editor
		want   int
	}{
		{
			"core",
			select5.Editor{
				select5.CursorPosition{
					X: 10,
					Y: 4,
				},
				nil,
				nil,
				sl,
			},
			2*5 + 1,
		},
		{
			"core",
			select5.Editor{
				select5.CursorPosition{
					X: 3,
					Y: 5,
				},
				nil,
				nil,
				sl,
			},
			18,
		},
		{
			"one mb character",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				nil,
				sl,
			},
			2,
		},
		{
			"blank line",
			select5.Editor{
				select5.CursorPosition{
					X: 0,
					Y: 1,
				},
				nil,
				nil,
				sl,
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.editor.GetLineVisibleLength(); got != tt.want {
				t.Errorf("select5.Editor.GetLineVisibleLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
