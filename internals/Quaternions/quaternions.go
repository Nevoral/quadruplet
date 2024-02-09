package Quaternions

import "math"

// NewQuaternion Creates new Quaternion
func NewQuaternion(angle Angle, normal *Vector3D) *Quaternion {
	halfAngle := (angle / 2).ToRadian()
	sinHalfAngle := math.Sin(halfAngle)
	return &Quaternion{
		W: math.Cos(halfAngle),
		X: normal[0] * sinHalfAngle,
		Y: normal[1] * sinHalfAngle,
		Z: normal[2] * sinHalfAngle,
	}
}

type Quaternion struct {
	W, X, Y, Z float64
}

// Multiply Quaternions multiplication
func (q *Quaternion) Multiply(r *Quaternion) *Quaternion {
	return &Quaternion{
		W: q.W*r.W - q.X*r.X - q.Y*r.Y - q.Z*r.Z,
		X: q.W*r.X + q.X*r.W + q.Y*r.Z - q.Z*r.Y,
		Y: q.W*r.Y - q.X*r.Z + q.Y*r.W + q.Z*r.X,
		Z: q.W*r.Z + q.X*r.Y - q.Y*r.X + q.Z*r.W,
	}
}

// Conjugate returns the conjugate of the quaternion is the same as Inverse of rotational quaternion
func (q *Quaternion) Conjugate() *Quaternion {
	return &Quaternion{W: q.W, X: -q.X, Y: -q.Y, Z: -q.Z}
}

func (q *Quaternion) Normalize() *Quaternion {
	mag := math.Sqrt(q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z)
	return &Quaternion{q.W / mag, q.X / mag, q.Y / mag, q.Z / mag}
}

func (q *Quaternion) Inverse() *Quaternion {
	conjugate := q.Conjugate()
	normSquared := q.W*q.W + q.X*q.X + q.Y*q.Y + q.Z*q.Z
	return &Quaternion{
		W: conjugate.W / normSquared,
		X: conjugate.X / normSquared,
		Y: conjugate.Y / normSquared,
		Z: conjugate.Z / normSquared,
	}
}

// ActiveRotation refers to rotating the object itself while the coordinate system remains fixed.
// In an active rotation, you're essentially asking, "How does the object look from the same viewpoint after it's been rotated?"
// The formula p' = q⁻¹pq is used for active rotations. Here, q is the rotation quaternion, p is the point represented as a quaternion, and q⁻¹ is the inverse of q.
func (q *Quaternion) ActiveRotation(vec *Vector3D) *Vector3D {
	// Convert the vec to a quaternion
	Qvec := &Quaternion{W: 0, X: vec[0], Y: vec[1], Z: vec[2]}

	// Rotate the vec: Q^(-1) * P * Q
	rotatedPoint := q.Conjugate().Multiply(Qvec).Multiply(q)

	// Extract the coordinates of the rotated vec
	return &Vector3D{rotatedPoint.X, rotatedPoint.Y, rotatedPoint.Z}
}

// PassiveRotation  is about rotating the coordinate system around the object. The object doesn't move; instead, the way we look at the object changes.
// In a passive rotation, you're essentially changing your viewpoint around a stationary object.
// The formula p' = qpq⁻¹ represents passive rotation.
func (q *Quaternion) PassiveRotation(vec *Vector3D) *Vector3D {
	// Convert the vec to a quaternion
	Qvec := &Quaternion{W: 0, X: vec[0], Y: vec[1], Z: vec[2]}

	// Rotate the vec: Q * P * Q^(-1)
	rotatedPoint := q.Multiply(Qvec).Multiply(q.Conjugate())

	// Extract the coordinates of the rotated vec
	return &Vector3D{rotatedPoint.X, rotatedPoint.Y, rotatedPoint.Z}
}
