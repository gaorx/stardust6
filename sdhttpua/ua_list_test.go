package sdhttpua

import (
	"github.com/gaorx/stardust6/sdstrings"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserAgents(t *testing.T) {
	is := assert.New(t)

	// 两个UA，mac的比例是20%，win的比例是80%
	csv := sdstrings.TrimMargin(`
	|20,"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36"
	|80,"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.19582"
	`, "|")
	ual, err := ParseCSV(csv, false)
	is.NoError(err)

	// IsEmpty
	is.False(ual.IsEmpty())

	// Len
	is.Equal(2, ual.Len())

	// All
	is.EqualValues(ual.list, ual.All())

	// Sub
	is.Equal(1, ual.Sub(PlatformIsWindows()).Len())
	is.Equal(1, ual.Sub(PlatformIsMacintosh()).Len())

	// Find
	is.Len(ual.Find(PlatformIsWindows()), 1)
	is.Len(ual.Find(PlatformIsMacintosh()), 1)

	// FindRaw
	is.Len(ual.FindRaw(PlatformIsWindows()), 1)
	is.Len(ual.FindRaw(PlatformIsMacintosh()), 1)

	// Platforms
	is.EqualValues([]string{"Macintosh", "Windows"}, ual.Platforms())

	// OSes
	is.EqualValues([]string{"Intel Mac OS X 10_15_7", "Windows 10"}, ual.OSes())

	// BrowserEngines
	is.EqualValues([]string{"AppleWebKit", "EdgeHTML"}, ual.BrowserEngines())

	// BrowserNames
	is.EqualValues([]string{"Chrome", "Edge"}, ual.BrowserNames())

	// Samples
	samples := ual.Samples(100000)
	var numMac, numWin int
	for _, ua := range samples {
		if ua.Platform == "Macintosh" {
			numMac++
		} else if ua.Platform == "Windows" {
			numWin++
		}
	}
	macRatio, winRatio := float64(numMac)/100000, float64(numWin)/100000
	is.True(macRatio >= 0.15 && macRatio <= 0.25)
	is.True(winRatio >= 0.75 && winRatio <= 0.85)
}
