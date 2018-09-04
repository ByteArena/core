package soccer

func systemPhysics(game *SoccerGame, dt float64) {
	for _, entityresult := range game.physicalView.Get() {
		physicalAspect := entityresult.Components[game.physicalBodyComponent].(*PhysicalBody)

		if physicalAspect.static {
			continue
		}

		if physicalAspect.GetVelocity().Mag() > 0.01 {
			physicalAspect.SetOrientation(physicalAspect.GetVelocity().Angle())
		}
	}

	///////////////////////////////////////////////////////////////////////////
	// On simule le monde physique
	///////////////////////////////////////////////////////////////////////////

	//before := time.Now()

	game.PhysicalWorld.Step(
		dt,
		4, // velocityIterations; higher improves stability; default 8 in testbed
		2, // positionIterations; higher improve overlap resolution; default 3 in testbed
	)

	//log.Println("Physical world step took ", float64(time.Now().UnixNano()-before.UnixNano())/1000000.0, "ms")
}
