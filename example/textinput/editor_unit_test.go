package main

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
		editor Editor
		want   string
	}{
		{
			"base",
			Editor{
				CursorPosition{
					X: 12,
					Y: 3,
				},
				nil,
				r1,
				sl,
				select5.KeyEvent{},
			},
			sc,
		},
		{
			"document head",
			Editor{
				CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r2,
				sl2,
				select5.KeyEvent{},
			},
			sc,
		},
		{
			"document end",
			Editor{
				CursorPosition{
					X: len(sl[len(sl)-1]) - 1,
					Y: len(sl) - 1,
				},
				nil,
				r3,
				sl3,
				select5.KeyEvent{},
			},
			sc,
		},
		{
			"line end",
			Editor{
				CursorPosition{
					X: len(sl[3]) - 1,
					Y: 3,
				},
				nil,
				r4,
				sl4,
				select5.KeyEvent{},
			},
			sc,
		},
		{
			"line head",
			Editor{
				CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r5,
				sl5,
				select5.KeyEvent{},
			},
			sc,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutEnter()
			if got := strings.Join(tt.editor.Line, ""); got != tt.want {
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
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
		editor Editor
		want   string
	}{
		{
			"middle",
			Editor{
				CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				r1,
				sl,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 6,
					Y: 0,
				},
				nil,
				r2,
				sl,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r3,
				sl,
				select5.KeyEvent{},
			},
			"\n" + strings.Join(sl, "\n"),
		},
		{
			"document end",
			Editor{
				CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				r4,
				sl,
				select5.KeyEvent{},
			},
			strings.Join(sl, "\n") + "\n",
		},
		{
			"line head",
			Editor{
				CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r5,
				sl,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl[3]),
					Y: 3,
				},
				nil,
				r2,
				sl,
				select5.KeyEvent{},
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
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
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
		editor Editor
		want   string
	}{
		{
			"middle",
			Editor{
				CursorPosition{
					X: 5,
					Y: 3,
				},

				nil,
				r5,
				sl1,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r4,
				sl3,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl[len(sl)-1]),
					Y: len(sl) - 1,
				},
				nil,
				r3,
				sl4,
				select5.KeyEvent{},
			},
			strings.Join(sl, "\n") + "A",
		},
		{
			"line head",
			Editor{
				CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r2,
				sl5,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl[3]),
					Y: 3,
				},
				nil,
				r1,
				sl,
				select5.KeyEvent{},
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
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
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
		editor Editor
		want   string
	}{
		{
			"middle",
			Editor{
				CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				r1,
				sl,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 6,
					Y: 0,
				},
				nil,
				r5,
				sl1,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r4,
				sl2,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl3[len(sl3)-1]),
					Y: len(sl3) - 1,
				},
				nil,
				r3,
				sl3,
				select5.KeyEvent{},
			},
			strings.Join(sl3, "\n"),
		},
		{
			"line head",
			Editor{
				CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r2,
				sl4,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl5[3]),
					Y: 3,
				},
				nil,
				r1,
				sl5,
				select5.KeyEvent{},
			},
			strings.Join(sl5, "\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.editor.PutDelete()
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
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
		editor Editor
		want   string
	}{
		{
			"middle",
			Editor{
				CursorPosition{
					X: 5,
					Y: 3,
				},
				nil,
				r1,
				sl,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 0,
					Y: 0,
				},
				nil,
				r2,
				sl2,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: len(sl3[len(sl3)-1]),
					Y: len(sl3) - 1,
				},
				nil,
				r3,
				sl3,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 0,
					Y: 3,
				},
				nil,
				r4,
				sl4,
				select5.KeyEvent{},
			},
			`Haruninari
Nosonoso
KumasanOkimashita
Aayokunetato
Akubishinagara`,
		},
		{
			"line end",
			Editor{
				CursorPosition{
					X: len(sl5[3]),
					Y: 3,
				},
				nil,
				r5,
				sl5,
				select5.KeyEvent{},
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
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
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
		editor Editor
		key    select5.KeyEvent
		want   string
	}{
		{
			"core",
			Editor{
				CursorPosition{
					X: 12,
					Y: 2,
				},
				nil,
				os.Stderr,
				sl1,
				select5.KeyEvent{},
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
			Editor{
				CursorPosition{
					X: 12,
					Y: 2,
				},
				nil,
				os.Stderr,
				sl2,
				select5.KeyEvent{},
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
				t.Errorf("Editor.Utf8Char() error = %v", err)
			}
			tt.editor.PutS(k)
			if got := strings.Join(tt.editor.Line, "\n"); got != tt.want {
				t.Errorf("Editor.Line2,12() = %x%x%x, want %x%x%x", tt.editor.Line[2][12], tt.editor.Line[2][13], tt.editor.Line[2][14], k[0], k[1], k[2])
				t.Errorf("Editor.Line() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEditor_LineVisibleLength(t *testing.T) {
	sl := []string{
		"春になり",
		"のそのそ熊さん",
		"起きました.",
		"「ああよく寝た」と",
		"あくびしながら",
	}
	sl1, sl2 :=
		make([]string, len(sl)), make([]string, len(sl))
	copy(sl1, sl)
	copy(sl2, sl)

	tests := []struct {
		name   string
		editor Editor
		want   int
	}{
		{
			"core",
			Editor{
				CursorPosition{
					X: 10,
					Y: 2,
				},
				nil,
				nil,
				sl1,
				select5.KeyEvent{},
			},
			11,
		},
		{
			"core",
			Editor{
				CursorPosition{
					X: 3,
					Y: 3,
				},
				nil,
				nil,
				sl1,
				select5.KeyEvent{},
			},
			18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.editor.GetLineVisibleLength(); got != tt.want {
				t.Errorf("Editor.GetLineVisibleLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
