package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/llgcode/draw2d/draw2dimg"
)

func main() {
	// Initialize the graphic context on an RGBA image
	chart := NewChart(15, dataframe.LoadRecords(
		[][]string{
			{Open, Low, High, Close},
			{"4.02", "4", "4.1", "4.08"},
			{"4.05", "4", "4.2", "4.15"},
			{"4.12", "4.1", "4.3", "4.2"},
			{"4.98", "4.8", "5", "4.9"},
			{"5.06", "5", "5.1", "5.05"},
		},
	),
	)

	gCtx := draw2dimg.NewGraphicContext(chart.canvas)

	for i := 1; i < len(chart.data); i++ {
		chart.NewCandle(i, gCtx)
	}

	// Save to file
	draw2dimg.SaveToPngFile("hello.png", chart.canvas)
}

type Chart struct {
	data        [][]string
	columnIndex map[string]int
	lineWidth   float64
	candleWidth uint
	canvas      *image.RGBA
	xMapping    func(float64) float64
	yMapping    func(float64) float64
}

const (
	High  = "high"
	Low   = "low"
	Open  = "open"
	Close = "close"
	Time  = "time"
	Max   = "max"
	Min   = "min"
)

func NewChart(candleWidth uint, df dataframe.DataFrame) Chart {
	if candleWidth == 0 {
		candleWidth = 1
	}
	records := df.Records()
	columnIndex := make(map[string]int)
	for i, value := range records[0] {
		columnIndex[value] = i
	}
	description := df.Describe().Records()
	var maxIdx, minIdx, highIdx, lowIdx int
	for row := range description {
		for col := range description[row] {
			switch description[row][col] {
			case Max:
				maxIdx = row
			case Min:
				minIdx = row
			case High:
				highIdx = col
			case Low:
				lowIdx = col
			}
		}
	}

	highest, err := strconv.ParseFloat(description[maxIdx][highIdx], 64)
	if err != nil {
		panic(err)
	}
	lowest, err := strconv.ParseFloat(description[minIdx][lowIdx], 64)
	if err != nil {
		panic(err)
	}
	canvasHight := 1080.0
	margin := 10.0
	num := canvasHight - 2*margin
	den := lowest - highest
	a := num / den

	return Chart{
		data:        records,
		columnIndex: columnIndex,
		lineWidth:   1.0,
		candleWidth: candleWidth,
		xMapping: func(x float64) float64 {
			return float64(candleWidth)*(2*x-0.5) + 0.5 + margin
		},
		yMapping: func(y float64) float64 {
			return a*(y-highest) + margin
		},
		canvas: image.NewRGBA(
			image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{int(candleWidth)*(2*df.Nrow()+1) + int(margin), int(canvasHight)},
			},
		),
	}
}

type MarketTrend string

var (
	Bullish MarketTrend = "bullish"
	Bear    MarketTrend = "bear"
)

func (c *Chart) NewCandle(idx int, gCtx *draw2dimg.GraphicContext) {
	const centeringFactor = 0.4
	high, err := strconv.ParseFloat(c.data[idx][c.columnIndex[High]], 32)
	if err != nil {
		panic(fmt.Errorf("parsing high: %s", err))
	}

	low, err := strconv.ParseFloat(c.data[idx][c.columnIndex[Low]], 32)
	if err != nil {
		panic(fmt.Errorf("parsing low: %s", err))
	}

	open, err := strconv.ParseFloat(c.data[idx][c.columnIndex[Open]], 32)
	if err != nil {
		panic(fmt.Errorf("parsing open: %s", err))
	}

	close, err := strconv.ParseFloat(c.data[idx][c.columnIndex[Close]], 32)
	if err != nil {
		panic(fmt.Errorf("parsing close: %s", err))
	}

	middle := (high + low) / 2

	gCtx.SetStrokeColor(color.RGBA{0x00, 0xff, 0x00, 0xff})
	gCtx.SetLineWidth(c.lineWidth)
	gCtx.MoveTo(c.xMapping(float64(idx)), c.yMapping(middle))
	gCtx.LineTo(c.xMapping(float64(idx)), c.yMapping(low))
	gCtx.MoveTo(c.xMapping(float64(idx)), c.yMapping(middle))
	gCtx.LineTo(c.xMapping(float64(idx)), c.yMapping(high))
	gCtx.FillStroke()

	if c.candleWidth > 1 {
		if close > open {
			gCtx.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
		} else {
			gCtx.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
		}
		gCtx.MoveTo(c.xMapping(float64(idx))-centeringFactor*float64(c.candleWidth), c.yMapping(open))
		gCtx.LineTo(c.xMapping(float64(idx))+centeringFactor*float64(c.candleWidth), c.yMapping(open))
		gCtx.LineTo(c.xMapping(float64(idx))+centeringFactor*float64(c.candleWidth), c.yMapping(close))
		gCtx.LineTo(c.xMapping(float64(idx))-centeringFactor*float64(c.candleWidth), c.yMapping(close))
		gCtx.Close()
		gCtx.FillStroke()
	} else {
		middle = (open + close) / 2
		gCtx.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
		gCtx.MoveTo(c.xMapping(float64(idx)), c.yMapping(middle))
		gCtx.LineTo(c.xMapping(float64(idx)), c.yMapping(open))
		gCtx.MoveTo(c.xMapping(float64(idx)), c.yMapping(middle))
		gCtx.LineTo(c.xMapping(float64(idx)), c.yMapping(close))
		gCtx.FillStroke()
	}
}
