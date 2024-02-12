package pythonAPI

import (
	"encoding/json"
	"fmt"
	"github.com/Nevoral/quadrupot/internals/Robot"
	"syscall"
)

type Message struct {
	Method   string
	Actions  []string
	Body     string
	Response []byte
}

func (m *Message) ActionsCall(r *Robot.Robot) {
	if m.Method == "POST" {
		for _, val := range m.Actions {
			switch val {
			case "setrobot":
				data, _ := syscall.ByteSliceFromString(m.Body)
				data = data[:len(data)-1]
				var rob Robot.Robot
				err := json.Unmarshal(data, &rob)
				if err != nil {
					fmt.Println(err)
				}
				r.Legs = rob.Legs
				r.HeadPoint = rob.HeadPoint
				r.BackPoint = rob.BackPoint
				r.CenterPoint = rob.CenterPoint
				r.NormalVec = rob.NormalVec
				r.Faze = rob.Faze
				return
			}
		}
	} else if m.Method == "GET" {
		for _, val := range m.Actions {
			switch val {
			case "getallpoints":
				// Marshal the config to TOML
				data, err := json.Marshal(&r)

				if err != nil {
					fmt.Println("Error encoding to JSON:", err)
					return
				}
				m.Response = append(m.Response, data...)
			}
		}
	} else {
		r.ResetPosition()
		m.Response = append(m.Response, byte('O'))
		m.Response = append(m.Response, byte('K'))
	}
	m.Response = append(m.Response, byte('>'))
	m.Response = append(m.Response, byte('>'))
	m.Response = append(m.Response, 10)
}
