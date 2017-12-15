package deathmatch

func systemLifecycle(deathmatch *DeathmatchGame) {
	for _, entityresult := range deathmatch.lifecycleView.Get() {
		lifecycleAspect := entityresult.Components[deathmatch.lifecycleComponent].(*Lifecycle)
		if lifecycleAspect.maxAge > 0 && (deathmatch.ticknum-lifecycleAspect.tickBirth) > lifecycleAspect.maxAge {
			lifecycleAspect.SetDeath(deathmatch.ticknum)
		}
	}
}
