package configuration

import (
	"testing"
)

func TestConfigurationCutExtentions_map(t *testing.T) {
	ex := ".exe .json .txt .png file.jpg testtesttesttesttest.dll"
	expected := []string{
		".exe",
		".json",
		".txt",
		".png",
		".jpg",
		".dll",
	}

	//act
	m, err := cutExtensions(ex)
	if err != nil {
		t.Fatalf("Expected nil error: %v", err)
	}

	for _, item := range expected {
		if _, ok := m[item]; !ok {
			t.Errorf("Expeted data not found: %s", item)
		}
	}
}

func TestConfigurationCutExtentions_error(t *testing.T) {
	indalidEx := "s jdfd json dll file txt"
	dot := ". . . . . . ."

	// act
	_, err := cutExtensions(indalidEx)
	if err == nil {
		t.Errorf("Expected non-nil error: %v", err)
	}

	_, err = cutExtensions(dot)
	if err == nil {
		t.Errorf("Expected non-nil error: %v", err)
	}
}

func BenchmarkConfigurationCutExtentions(b *testing.B) {
	//arange
	extention := `.doc .exe .png .jpg .otd .gitignore .gif
	.svg .tif .pdf .mp3 .mp4 .wav .bat .go .css .json .php .zip`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := cutExtensions(extention)
		if err != nil {
			b.Fatalf("failed bench, wrong data: %v", err)
		}
	}
}
