package models

import (
	"math/rand"
	"sync"
	"time"
)

var LaneStatus [NumLanes]bool
var LaneMutex sync.Mutex
const NumLanes  = 20

func WaitForVehiclePosition(id int, targetX float64) {
	for {
		vehiclePos := FindVehiclePosition(id)
		if vehiclePos.X >= targetX {
			break
		}
		time.Sleep(16 * time.Millisecond)
	}
}

func FindAvailableLane() (int, bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	rand.Seed(time.Now().UnixNano())
	lanes := rand.Perm(NumLanes)
	for _, l := range lanes {
		if !LaneStatus[l] {
			LaneStatus[l] = true
			return l, true
		}
	}
	return -1, false
}

func UpdateLaneStatus(lane int, status bool) {
	LaneMutex.Lock()
	defer LaneMutex.Unlock()
	if lane >= 0 && lane < NumLanes {
		LaneStatus[lane] = status
	}
}

func ManageLane(id int) {
	CreateVehicle(id)
	WaitForVehiclePosition(id, 100)
	lane, foundLane := FindAvailableLane()
	if !foundLane {
		ResetVehiclePosition(id)
		return
	}
	AssignLaneToVehicle(id, lane)
}
