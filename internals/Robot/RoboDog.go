package Robot

import (
	"encoding/json"
	"github.com/Nevoral/quadrupot/internals/Quaternions"
	"github.com/Nevoral/quadrupot/internals/handlers"
	"github.com/Nevoral/quadrupot/web/templates"
	"github.com/gofiber/fiber/v3"
	"log"
	"strconv"
)

func NewRobot(head, center, back *Quaternions.Point3D) *Robot {
	return &Robot{
		Legs:        nil,
		NormalVec:   &Quaternions.Point3D{0, 0, 1},
		HeadPoint:   head,
		CenterPoint: center,
		BackPoint:   back,
		Faze:        1,
	}
}

type Robot struct {
	Legs        []*Leg               `json:"Legs"`
	NormalVec   *Quaternions.Point3D `json:"NormalVec"`
	HeadPoint   *Quaternions.Point3D `json:"HeadPoint"`
	CenterPoint *Quaternions.Point3D `json:"CenterPoint"`
	BackPoint   *Quaternions.Point3D `json:"BackPoint"`
	Faze        int                  `json:"Faze"`
}

func (r *Robot) AddLeg(l *Leg, config func(leg *Leg)) {
	config(l)
	r.Legs = append(r.Legs, l)
}

type GraphData struct {
	Name string
	Data *[]map[string][3]float64
}

func (r *Robot) GetPoints() []byte {
	var m = []*GraphData{{
		Name: "body",
		Data: r.StringModelBody(40),
	}}

	var angles [][]Quaternions.Angle
	for _, l := range r.Legs {
		m = append(m, &GraphData{
			Name: l.Name,
			Data: l.GetLegPoints(),
		})
		angles = append(angles, l.GetCurrentAngles())
	}
	data := make(map[string]any)
	data["quadrubot"] = m
	data["angles"] = angles
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return jsonData
}

func (r *Robot) StringModelBody(lengthMotor float64) *[]map[string][3]float64 {
	var ch []map[string][3]float64
	var legPoint []*Quaternions.Point3D
	for _, leg := range r.Legs {
		legPoint = append(legPoint, leg.AbsolutPosition)
	}
	ch = append(ch, map[string][3]float64{"value": legPoint[0].DifPoints(&Quaternions.Point3D{lengthMotor, 0, 0}).GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[1].DifPoints(&Quaternions.Point3D{lengthMotor, 0, 0}).GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[1].GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[2].GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[2].SumPoints(&Quaternions.Point3D{lengthMotor, 0, 0}).GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[3].SumPoints(&Quaternions.Point3D{lengthMotor, 0, 0}).GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[3].GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[0].GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": legPoint[0].DifPoints(&Quaternions.Point3D{lengthMotor, 0, 0}).GetXYZ()})
	d := Quaternions.Point3D{legPoint[0].DifPoints(&Quaternions.Point3D{lengthMotor, 0, 0})[0], 0, legPoint[0][2]}
	ch = append(ch, map[string][3]float64{"value": d.GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": r.HeadPoint.GetXYZ()})
	ch = append(ch, map[string][3]float64{"value": r.BackPoint.GetXYZ()})
	return &ch
}

func (r *Robot) FindStandingPlane() {
	var tPoints []*Quaternions.Point3D
	for _, l := range r.Legs {
		tPoints = append(tPoints, l.TouchPosition)
	}

}

func (r *Robot) IsMassCenterInTriangle() bool {
	return true
}

func (r *Robot) TranslateRobot(offset *Quaternions.Point3D) {
	r.HeadPoint = r.HeadPoint.SumPoints(offset)
	r.CenterPoint = r.CenterPoint.SumPoints(offset)
	r.BackPoint = r.BackPoint.SumPoints(offset)
	for _, leg := range r.Legs {
		leg.AbsolutPosition = leg.AbsolutPosition.SumPoints(offset)
	}
}

func (r *Robot) SetTouchPoint() {
	for ind, leg := range r.Legs {
		if r.Faze-1 == ind {
			leg.TouchPosition = nil
		}
		leg.TouchPosition = leg.GetTouchPoint()
	}
}

func (r *Robot) ResetPosition() {
	r.Legs = []*Leg{}
	r.AddLeg(NewLeg(0, "Front Left Leg", &Quaternions.Point3D{160, 100, 0}), FrontLeftLeg)
	r.AddLeg(NewLeg(1, "Front Right Leg", &Quaternions.Point3D{160, -100, 0}), FrontRightLeg)
	r.AddLeg(NewLeg(2, "Back Left Leg", &Quaternions.Point3D{-150, 100, 0}), BackLeftLeg)
	r.AddLeg(NewLeg(3, "Back Right Leg", &Quaternions.Point3D{-150, -100, 0}), BackRightLeg)
}

func (r *Robot) SendGraphPage() func(c fiber.Ctx) error {
	r.ResetPosition()
	return handlers.GetHTML(templates.Graph(4))
}

func (r *Robot) SendDefaultConfig() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) (err error) {
		r.ResetPosition()
		c.Set("Content-Type", fiber.MIMEApplicationJSON)
		err = c.Send(r.GetPoints())
		if err != nil {
			return err
		}
		return nil
	}
}

func (r *Robot) SendGraphData() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) (err error) {
		r.SetTouchPoint()
		r.FindStandingPlane()
		//r.TranslateRobot(&quat.Point3D{0, 0, 240})

		err = r.ChangeAnglesJSON(c)
		if err != nil {
			return err
		}
		c.Set("Content-Type", fiber.MIMEApplicationJSON)
		err = c.Send(r.GetPoints())
		if err != nil {
			return err
		}
		return nil
	}
}

func (r *Robot) ChangeAnglesJSON(c fiber.Ctx) (err error) {
	type sliderData struct {
		SliderValues []string `json:"sliderValues"`
	}

	//Parse the incoming JSON data
	data := new(sliderData)
	if err = json.Unmarshal(c.Body(), &data); err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		var a [3]Quaternions.Angle
		for j := 0; j < 3; j++ {
			var p float64
			p, err = strconv.ParseFloat(data.SliderValues[i*3+j], 64)
			a[j] = Quaternions.Angle(p)
			if err != nil {
				return err
			}
		}
		r.Legs[i].ChangeAngle([3]Quaternions.Angle{a[0], a[1], a[2]})
	}
	return nil
}
