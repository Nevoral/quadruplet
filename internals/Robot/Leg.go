package Robot

import (
	"encoding/json"
	"fmt"
	"github.com/Nevoral/quadrupot/internals/Quaternions"
	"github.com/Nevoral/quadrupot/internals/handlers"
	"github.com/Nevoral/quadrupot/web/templates"
	"github.com/gofiber/fiber/v3"
	"log"
	"strconv"
)

func NewLeg(id int, name string, position *Quaternions.Point3D) *Leg {
	return &Leg{
		Id:              id,
		Name:            name,
		Joints:          nil,
		AbsolutPosition: position,
		configFunc:      nil,
		TouchPosition:   nil,
	}
}

type Leg struct {
	Id              int                  `json:"Id"`
	Name            string               `json:"Name"`
	Joints          []*Joint             `json:"Joints"`
	AbsolutPosition *Quaternions.Point3D `json:"AbsolutPosition"`
	configFunc      func(*Leg)           `json:"-"`
	TouchPosition   *Quaternions.Point3D `json:"TouchPosition"`
}

func (l *Leg) CreateJoints(center []Quaternions.Point3D, normalVec []Quaternions.Vector3D, angleConfig [][]Quaternions.Angle, armLen []float64) {
	l.Joints = []*Joint{}
	for ind := 0; ind < len(center); ind++ {
		l.Joints = append(l.Joints, NewJoint(ind, angleConfig[ind][0], angleConfig[ind][1], angleConfig[ind][2]).
			SetPoint(&center[ind]).
			SetNormalVector(&normalVec[ind]).
			SetArm(armLen[ind]).
			SetAngle(angleConfig[ind][0]),
		)
	}
}

func (l *Leg) GetCurrentAngles() (a []Quaternions.Angle) {
	for _, val := range l.Joints {
		a = append(a, val.CurrentAngle)
	}
	return a
}

func (l *Leg) GetTouchPoint() *Quaternions.Point3D {
	return l.Joints[len(l.Joints)-1].RelativePoint.SumPoints(l.AbsolutPosition)
}

// PrintTouchPoint Print coords in better format
func (l *Leg) PrintTouchPoint() {
	last := l.GetTouchPoint()
	fmt.Printf("\n[X; Y; Z] = [%.3f; %.3f; %.3f]\n", last[0], last[1], last[2])
}

// ComputeTotalTransformation computes the total transformation matrix for a given set of joints
func (l *Leg) ComputeTotalTransformation(jointID int, a Quaternions.Angle) {
	change := l.SetAngle(jointID, a)
	if !change {
		return
	}
	if len(l.Joints)-jointID > 1 {
		for i := jointID + 1; i < len(l.Joints); i++ {
			l.Joints[i] = l.Joints[jointID].RotateJoints(l.Joints[i])
		}
	}
}

// ComputeTotalTransformation2 computes the total transformation matrix for a given set of joints
func (l *Leg) ComputeTotalTransformation2() {
	var q *Quaternions.Quaternion
	var j *Joint
	for _, joint := range l.Joints {
		if q != nil {
			joint.RelativePoint = j.RelativePoint.RotatePointAroundCenter(joint.RelativePoint, q)
			joint.NormalVec = q.ActiveRotation(joint.NormalVec).Normalize()
		}

		if joint.AngleChange != 0 && q == nil {
			q = Quaternions.NewQuaternion(joint.AngleChange, joint.NormalVec)
			j = joint
		} else if joint.AngleChange != 0 {
			q = Quaternions.NewQuaternion(joint.AngleChange, joint.NormalVec).Multiply(q)
		}
	}
}

// ComputeTotalTransformation3 computes the total transformation matrix for a given set of joints
func (l *Leg) ComputeTotalTransformation3() {
	var j []*Joint
	for _, joint := range l.Joints {
		if j != nil {
			for _, rot := range j {
				joint = rot.RotateJoints(joint)
			}
		}
		if joint.AngleChange != 0 {
			j = append(j, joint)
		}
	}
}

func (l *Leg) SetAngle(jointID int, a Quaternions.Angle) bool {
	l.Joints[jointID].AngleChange = l.Joints[jointID].CurrentAngle - a
	if l.Joints[jointID].AngleChange == 0 {
		return false
	}
	l.Joints[jointID].SetAngle(a)
	return true
}

func (l *Leg) CheckLengthOfArm() error {
	var j *Joint
	for _, joint := range l.Joints {
		if j == nil {
			j = joint
			continue
		}
		if ok, dist := j.CheckLenArm(joint.RelativePoint); !ok {
			return fmt.Errorf("\nJoint %d isn't at right distance is: %.3f\n", joint.Id, dist)
		}
		j = joint
	}
	return nil
}

func (l *Leg) GetLegPoints() *[]map[string][3]float64 {
	var ch []map[string][3]float64
	for _, joint := range l.Joints {
		ch = append(ch, map[string][3]float64{"value": joint.RelativePoint.SumPoints(l.AbsolutPosition).GetXYZ()})
	}
	return &ch
}

func (l *Leg) ChangeAngle(theta [3]Quaternions.Angle) {
	for ind, val := range theta {
		//l.ComputeTotalTransformation(ind, val)
		_ = l.SetAngle(ind, val)
	}
	l.ComputeTotalTransformation2()
	if err := l.CheckLengthOfArm(); err != nil {
		fmt.Println(err)
	}
}

func (l *Leg) SendGraphPageLeg() func(c fiber.Ctx) error {
	return handlers.GetHTML(templates.Graph(1))
}

func (l *Leg) SendGraphData() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) (err error) {
		var graph []*GraphData
		graph = append(graph, &GraphData{Name: "test",
			Data: &[]map[string][3]float64{{"value": {0, 0, 0}},
				{"value": {0, 55, 0}},
				{"value": {0, 55, -240}},
				{"value": {240, 55, -240}},
			},
		})

		err = l.ChangeAnglesJSON(c.Body())
		if err != nil {
			return err
		}
		graph = append(graph, &GraphData{Name: l.Name,
			Data: l.GetLegPoints()})

		c.Set("Content-Type", fiber.MIMEApplicationJSON)
		// Convert data to JSON
		jsonData, err := json.Marshal(graph)
		if err != nil {
			log.Fatal(err)
		}
		err = c.Send(jsonData)
		if err != nil {
			return err
		}
		return nil
	}
}

func (l *Leg) ChangeAnglesJSON(d []byte) (err error) {
	type sliderData struct {
		SliderValues []string `json:"sliderValues"`
	}

	//Parse the incoming JSON data
	data := new(sliderData)
	if err = json.Unmarshal(d, &data); err != nil {
		return err
	}

	var a [3]Quaternions.Angle
	for j := 0; j < 3; j++ {
		var p float64
		p, err = strconv.ParseFloat(data.SliderValues[j], 64)
		a[j] = Quaternions.Angle(p)
		if err != nil {
			return err
		}
	}
	l.ChangeAngle([3]Quaternions.Angle{a[0], a[1], a[2]})
	return nil
}
