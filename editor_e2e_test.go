package select5_test

import (
	"github.com/g1eng/select5"
	"os"
	"testing"
	"time"
)

func TestEditor_Scenario_SimpleInput(t *testing.T) {
	asciiText := `Haruninari
Nosonoso
Kumasan
Okimashita
Aayokunetato
Akubishinagara`
	multibyteText := `春になり
のそのそ熊さん
起きました。
「ああよく寝た」と
あくびしながら`

	tests := []struct {
		name     string
		fixtures string
		want     string
	}{
		{
			name:     "Scenario_Ascii",
			fixtures: asciiText,
			want:     asciiText,
		}, {
			name:     "Scenario_Multibyte",
			fixtures: multibyteText,
			want:     multibyteText,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//prepare pipe for key input capturing
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()
			defer w.Close()

			oldStdin := os.Stdin
			os.Stdin = r
			defer func() { os.Stdin = oldStdin }()
			ed := select5.NewEditor()

			var got string
			go func() {
				got = ed.Edit()
			}()
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte(tt.fixtures))
			w.Write([]byte{select5.CtrlD})
			time.Sleep(100 * time.Millisecond)
			if got != tt.fixtures {
				t.Fatalf("ed.Edit(): got %q, want %q", got, tt.fixtures)
			}
		})
	}
}
