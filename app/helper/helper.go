package helper

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"

	"golang.org/x/text/encoding/charmap"
)

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func CalculateTotal(quantity uint, price float64) float64 {
	total := float64(quantity) * price
	total = RoundFloat(total, 2)
	return total
}

func NewDecoderXML(body io.ReadCloser) *xml.Decoder {

	d := xml.NewDecoder(body)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	return d
}

func BubbleSort(ar []int) []int {
	for i := 0; i < len(ar); i++ {
		for j := i; j < len(ar); j++ {
			if ar[i] > ar[j] {
				// swap(ar, i, j)
				ar[i], ar[j] = ar[j], ar[i]
			}
		}
	}
	return ar
}

func Unique(sl []int) []int {

	var unique = make(map[int]bool)
	for _, v := range sl {
		if _, ok := unique[v]; !ok {
			unique[v] = true
		}
	}
	var newSlice []int
	for k := range unique {
		newSlice = append(newSlice, k)
	}
	return newSlice
}
