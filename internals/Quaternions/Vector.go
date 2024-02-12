package Quaternions

import (
	"fmt"
	"math"
)

type Vector3D [3]float64

func (v *Vector3D) GetXYZ() [3]float64 {
	return [3]float64{v[0], v[1], v[2]}
}

func (v *Vector3D) PrintVector() {
	fmt.Printf("(X,Y,Z) = (%.2f, %.2f, %.2f)\n", v[0], v[1], v[2])
}

func (v *Vector3D) Add(vec *Vector3D) *Vector3D {
	return &Vector3D{v[0] + vec[0], v[1] + vec[1], v[2] + vec[2]}
}

func (v *Vector3D) Subtract(vec *Vector3D) *Vector3D {
	return &Vector3D{v[0] - vec[0], v[1] - vec[1], v[2] - vec[2]}
}

func (v *Vector3D) Scale(scalar float64) *Vector3D {
	return &Vector3D{v[0] * scalar, v[1] * scalar, v[2] * scalar}
}

func (v *Vector3D) VecLength() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v *Vector3D) DotProduct(vec *Vector3D) float64 {
	return v[0]*vec[0] + v[1]*vec[1] + v[2]*vec[2]
}

func (v *Vector3D) Normalize() *Vector3D {
	length := v.VecLength()
	if length == 0 {
		return v
	}
	v[0], v[1], v[2] = v[0]/length, v[1]/length, v[2]/length
	return v
}

// CrossProduct calculates the cross product of two vectors in 3D space.
func (v *Vector3D) CrossProduct(vec *Vector3D) *Vector3D {
	return &Vector3D{
		v[1]*vec[2] - v[2]*vec[1],
		v[2]*vec[0] - v[0]*vec[2],
		v[0]*vec[1] - v[1]*vec[0],
	}
}

func (v *Vector3D) NormalToPlaneBy2Vec(vec *Vector3D) *Vector3D {
	return v.CrossProduct(vec).Normalize()
}

func (v *Vector3D) AngleBetweenVectors(vec *Vector3D) Angle {
	return Angle(math.Acos(v.DotProduct(vec) / (v.VecLength() * vec.VecLength())))
}
