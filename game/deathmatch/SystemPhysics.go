package deathmatch

func systemPhysics(deathmatch *DeathmatchGame, dt float64) {
	for _, entityresult := range deathmatch.physicalView.Get() {
		physicalAspect := entityresult.Components[deathmatch.physicalBodyComponent].(*PhysicalBody)

		if physicalAspect.static {
			continue
		}

		if physicalAspect.skipThisTurn {
			physicalAspect.skipThisTurn = false
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

	deathmatch.PhysicalWorld.Step(
		dt,
		4, // velocityIterations; higher improves stability; default 8 in testbed
		2, // positionIterations; higher improve overlap resolution; default 3 in testbed
	)

	//log.Println("Physical world step took ", float64(time.Now().UnixNano()-before.UnixNano())/1000000.0, "ms")
}
