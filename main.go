package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
)

const (
	width, height = 800, 520            // canvas size in pixels
	cells         = 125                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.2        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	surface := fmt.Sprintf(`<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Сиськи</title>
		</head>
		<body>
		<h1>Опа! Сиськи!</h1>
		<svg xmlns="http://www.w3.org/2000/svg" style="stroke: grey; fill: white; stroke-width: 0.7" width="%d" height="%d">`, width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)
			poligon := fmt.Sprintf(`<polygon points="%g,%g %g,%g %g,%g %g,%g"/>`,
				ax, ay, bx, by, cx, cy, dx, dy)
			surface = surface + poligon
		}
	}
	surface = surface + "</svg></body></html>"
	w.Write([]byte(surface))
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	fmt.Println("Open in your browser localhost:", port)
	http.ListenAndServe(":"+port, mux)

}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy
}

func f(x, y float64) float64 {
	posX := math.Pow(x+4, 2)
	negX := math.Pow(x-4, 2)
	posY := math.Pow(y+4, 2)
	negY := math.Pow(y-4, 2)
	posElem := (posX + posY) * (posX + posY)
	negElem := (negX + negY) * (negX + negY)
	return math.Exp(-posElem/1000) + math.Exp(-negElem/1000) + 0.1*(math.Exp(-posElem)+math.Exp(-negElem))
}

// func f(x, y float64) float64 {
// 	r := math.Hypot(x, y) // distance from (0,0)
// 	return math.Sin(r) / r
// }
