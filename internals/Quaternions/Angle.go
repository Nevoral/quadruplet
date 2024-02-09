package Quaternions

import "math"

type Angle float64

func (a Angle) ToRadian() float64 {
	return float64(a * math.Pi / 180)
}
