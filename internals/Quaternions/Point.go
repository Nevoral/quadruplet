package Quaternions

import (
	"fmt"
	"math"
)

type Point3D [3]float64

func (p *Point3D) GetXYZ() [3]float64 {
	return [3]float64{p[0], p[1], p[2]}
}

func (p *Point3D) PrintPoint() {
	fmt.Printf("[X,Y,Z] = [%.2f, %.2f, %.2f]\n", p[0], p[1], p[2])
}

func (p *Point3D) RotatePointAroundCenter(point *Point3D, rot *Quaternion) *Point3D {
	vec := p.CreateVector3D(point)
	vec = rot.ActiveRotation(vec)
	point = p.TranslatePointOnVec(vec)
	return point
}

// AngleWith Computes the angle between two points with respect to the origin.
func (p *Point3D) AngleWith(point *Point3D) Angle {
	dot := p[0]*point[0] + p[1]*point[1] + p[2]*point[2]
	mag1 := math.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
	mag2 := math.Sqrt(point[0]*point[0] + point[1]*point[1] + point[2]*point[2])
	return Angle(math.Acos(dot / (mag1 * mag2)))
}

func (p *Point3D) Reflect(across *Point3D) *Point3D {
	return &Point3D{
		2*across[0] - p[0],
		2*across[1] - p[1],
		2*across[2] - p[2],
	}
}

func (p *Point3D) Midpoint(target *Point3D) *Point3D {
	return p.Lerp(target, 0.5)
}

func (p *Point3D) Lerp(target *Point3D, t float64) *Point3D {
	return &Point3D{
		p[0] + t*(target[0]-p[0]),
		p[1] + t*(target[1]-p[1]),
		p[2] + t*(target[2]-p[2]),
	}
}

func (p *Point3D) MultiplyScalar(scalar float64) *Point3D {
	return &Point3D{p[0] * scalar, p[1] * scalar, p[2] * scalar}
}

func (p *Point3D) DivideScalar(scalar float64) *Point3D {
	return &Point3D{p[0] / scalar, p[1] / scalar, p[2] / scalar}
}

func (p *Point3D) SumPoints(point *Point3D) *Point3D {
	return &Point3D{p[0] + point[0], p[1] + point[1], p[2] + point[2]}
}

func (p *Point3D) DifPoints(point *Point3D) *Point3D {
	return &Point3D{p[0] - point[0], p[1] - point[1], p[2] - point[2]}
}

func (p *Point3D) Distance(point *Point3D) float64 {
	return math.Sqrt(math.Pow(p[0]-point[0], 2) + math.Pow(p[1]-point[1], 2) + math.Pow(p[2]-point[2], 2))
}

func (p *Point3D) ToVector3D() *Vector3D {
	return &Vector3D{p[0], p[1], p[2]}
}

func (p *Point3D) CreateVector3D(point *Point3D) *Vector3D {
	return &Vector3D{point[0] - p[0], point[1] - p[1], point[2] - p[2]}
}

func (p *Point3D) NormalToPlaneBy3Points(b, c *Point3D) *Vector3D {
	ab := b.CreateVector3D(p)
	ac := c.CreateVector3D(p)
	return ab.CrossProduct(ac).Normalize()
}

func (p *Point3D) TranslatePointOnVec(vec *Vector3D) *Point3D {
	return &Point3D{p[0] + vec[0], p[1] + vec[1], p[2] + vec[2]}
}
