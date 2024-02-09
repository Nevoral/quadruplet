package Robot

import (
	"github.com/Nevoral/quadrupot/internals/Quaternions"
	"math"
)

func NewJoint(id int, currentAngle, maxAngle, motionRange Quaternions.Angle) *Joint {
	return &Joint{
		Id:            id,
		NormalVec:     nil,
		RelativePoint: nil,
		MotionRange:   motionRange,
		MaxAngle:      maxAngle,
		CurrentAngle:  currentAngle,
		AngleChange:   0,
		Arm:           0,
	}
}

type Joint struct {
	Id            int
	NormalVec     *Quaternions.Vector3D
	RelativePoint *Quaternions.Point3D
	MotionRange   Quaternions.Angle
	MaxAngle      Quaternions.Angle
	CurrentAngle  Quaternions.Angle
	AngleChange   Quaternions.Angle
	Arm           float64
}

func (j *Joint) SetAngle(angle Quaternions.Angle) *Joint {
	if angle > j.MaxAngle || angle < j.MaxAngle+j.MotionRange {
		j.CurrentAngle = angle
		return j
	}
	return j
}

func (j *Joint) SetPoint(point *Quaternions.Point3D) *Joint {
	j.RelativePoint = point
	return j
}

func (j *Joint) SetNormalVector(vec *Quaternions.Vector3D) *Joint {
	j.NormalVec = vec.Normalize()
	return j
}

func (j *Joint) SetArm(arm float64) *Joint {
	j.Arm = arm
	return j
}

func (j *Joint) RotateJoints(joint *Joint) *Joint {
	q := Quaternions.NewQuaternion(j.AngleChange, j.NormalVec.Normalize())
	joint.RelativePoint = j.RelativePoint.RotatePointAroundCenter(joint.RelativePoint, q)
	joint.NormalVec = q.ActiveRotation(joint.NormalVec).Normalize()
	return joint
}

func (j *Joint) CheckLenArm(point *Quaternions.Point3D) (bool, float64) {
	if distance := j.RelativePoint.Distance(point); math.Abs(distance-j.Arm) > .01 {
		//fmt.Printf("\n%.2f - %.2f\n", j.Arm, distance)
		return false, j.Arm - distance
	}
	return true, 0
}
