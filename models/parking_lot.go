package models

import "time"

func MoveVehiclesLogic() {
	for i := len(Vehicles) - 1; i >= 0; i-- {
		vehicle := &Vehicles[i]
		if vehicle.Lane == -1 && !vehicle.IsEntering {
			vehicle.Position.X += 10
			if vehicle.Position.X > 100 {
				vehicle.Position.X = 100
			}
		} else if vehicle.Lane != -1 && !vehicle.Parked {
			var targetX, targetY float64
			laneWidth := 600.0 / 10
			if vehicle.Lane < 10 {
				targetX = 100.0 + float64(vehicle.Lane)*laneWidth + laneWidth/2
				targetY = 400 + (500-350)/2
			} else {
				targetX = 100.0 + float64(vehicle.Lane-10)*laneWidth + laneWidth/2
				targetY = 100 + (250-100)/2
			}
			moveTowardsTarget(vehicle, targetX, targetY)
			if nearTarget(vehicle.Position.X, vehicle.Position.Y, targetX, targetY) {
				ParkVehicle(vehicle, targetX, targetY)
			}
		}
	}
	ExitVehicleLogic()
}

func ExitVehicleLogic() {
	for i := len(Vehicles) - 1; i >= 0; i-- {
		vehicle := &Vehicles[i]
		if vehicle.Parked && time.Now().After(vehicle.ExitTime) && !vehicle.IsEntering {
			if !vehicle.Teleporting {
				vehicle.Teleporting = true
				vehicle.Position.X = 100 + float64(vehicle.Lane%10)*60
				vehicle.Position.Y = 400
			} else {
				if !CarEnteringOrExiting {
					moveTowardsTarget(vehicle, 50, 400)
					if nearTarget(vehicle.Position.X, vehicle.Position.Y, 50, 400) {
						UpdateLaneStatus(vehicle.Lane, false)
						removeVehicle(i)
					}
				}
			}
		}
	}
}

func moveTowardsTarget(vehicle *Vehicle, targetX, targetY float64) {
	speed := 5.0
	if vehicle.Position.X < targetX {
		vehicle.Position.X += speed
		if vehicle.Position.X > targetX {
			vehicle.Position.X = targetX
		}
	} else if vehicle.Position.X > targetX {
		vehicle.Position.X -= speed
		if vehicle.Position.X < targetX {
			vehicle.Position.X = targetX
		}
	}

	if vehicle.Position.Y < targetY {
		vehicle.Position.Y += speed
		if vehicle.Position.Y > targetY {
			vehicle.Position.Y = targetY
		}
	} else if vehicle.Position.Y > targetY {
		vehicle.Position.Y -= speed
		if vehicle.Position.Y < targetY {
			vehicle.Position.Y = targetY
		}
	}
}

func nearTarget(posX, posY, targetX, targetY float64) bool {
	const tolerance = 5.0
	return abs(posX-targetX) < tolerance && abs(posY-targetY) < tolerance
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
