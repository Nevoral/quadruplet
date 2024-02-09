package Robot

import (
	"github.com/Nevoral/quadrupot/internals/Quaternions"
)

func FrontLeftLeg(r *Leg) {
	center := []Quaternions.Point3D{{0, 0, 0}, {0, 55, 0}, {0, 55, -240}, {240, 55, -240}}
	normalVec := []Quaternions.Vector3D{{1, 0, 0}, {0, 1, 0}, {0, 1, 0}, {0, 0, 0}}
	angleConfig := [][]Quaternions.Angle{{0, -48, 96}, {-90, -200, 290}, {-90, -156, 108}, {0, 0, 0}}
	armLen := []float64{55, 240, 240, 0}
	r.CreateJoints(center, normalVec, angleConfig, armLen)
	r.ConfigFunc = FrontLeftLeg
}

func FrontRightLeg(r *Leg) {
	center := []Quaternions.Point3D{{0, 0, 0}, {0, -55, 0}, {0, -55, -240}, {240, -55, -240}}
	normalVec := []Quaternions.Vector3D{{-1, 0, 0}, {0, 1, 0}, {0, 1, 0}, {0, 0, 0}}
	angleConfig := [][]Quaternions.Angle{{0, -48, 96}, {-90, -200, 290}, {-90, -156, 108}, {0, 0, 0}}
	armLen := []float64{55, 240, 240, 0}
	r.CreateJoints(center, normalVec, angleConfig, armLen)
	r.ConfigFunc = FrontRightLeg
}

func BackRightLeg(r *Leg) {
	center := []Quaternions.Point3D{{0, 0, 0}, {0, -55, 0}, {0, -55, -240}, {240, -55, -240}}
	normalVec := []Quaternions.Vector3D{{-1, 0, 0}, {0, 1, 0}, {0, 1, 0}, {0, 0, 0}}
	angleConfig := [][]Quaternions.Angle{{0, -48, 96}, {-90, -200, 290}, {-90, -156, 108}, {0, 0, 0}}
	armLen := []float64{55, 240, 240, 0}
	r.CreateJoints(center, normalVec, angleConfig, armLen)
	r.ConfigFunc = BackRightLeg
}

func BackLeftLeg(r *Leg) {
	center := []Quaternions.Point3D{{0, 0, 0}, {0, 55, 0}, {0, 55, -240}, {240, 55, -240}}
	normalVec := []Quaternions.Vector3D{{1, 0, 0}, {0, 1, 0}, {0, 1, 0}, {0, 0, 0}}
	angleConfig := [][]Quaternions.Angle{{0, -48, 96}, {-90, -200, 290}, {-90, -156, 108}, {0, 0, 0}}
	armLen := []float64{55, 240, 240, 0}
	r.CreateJoints(center, normalVec, angleConfig, armLen)
	r.ConfigFunc = BackLeftLeg
}
