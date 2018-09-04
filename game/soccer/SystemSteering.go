package soccer

func systemSteering(game *SoccerGame) {
	for _, entityresult := range game.steeringView.Get() {
		steeringAspect := entityresult.Components[game.steeringComponent].(*Steering)
		physicalAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)

		steers := steeringAspect.PopPendingSteers()
		if len(steers) == 0 {
			continue
		}

		targetPoint := steers[0]

		steering := targetPoint.Sub(physicalAspect.GetPosition())

		// velocity := physicalAspect.GetVelocity()
		// orientation := physicalAspect.GetOrientation()

		// prevmag := velocity.Mag()
		// diff := steering.Mag() - prevmag

		// maxSteeringForce := steeringAspect.GetMaxSteeringForce()
		// // maxAngularVelocity := physicalAspect.GetMaxAngularVelocity()
		// maxSpeed := physicalAspect.GetMaxSpeed()
		// if math.Abs(diff) > maxSteeringForce {
		// 	if diff > 0 {
		// 		steering = steering.SetMag(prevmag + maxSteeringForce)
		// 	} else {
		// 		steering = steering.SetMag(prevmag - maxSteeringForce)
		// 	}
		// }

		// abssteering := trigo.Limit(maxSpeed)
		// LocalAngleToAbsoluteAngleVec(orientation, steering, &maxAngularVelocity).

		// limitedSteering := steering.Limit(maxSpeed)

		physicalAspect.SetVelocity(steering) // absolute steering received, abs steering provided
	}
}
