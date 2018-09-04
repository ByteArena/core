package soccer

import (
	json "encoding/json"
	"errors"

	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/types"
	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/utils/vector"
)

func systemMutations(game *SoccerGame, mutations []types.AgentMutationBatch) {

	for _, batch := range mutations {

		// Ordering actions
		// This is important because operations like shooting are taken from the previous position of the agent
		// 1. Non-movement actions (shoot, etc.)
		// 2. Movement actions

		// 1. Non-movement actions
		//		for _, mutation := range batch.Mutations {
		//			switch mutation.GetMethod() {
		//			case "shoot":
		//				{
		//					if err := handleShootMutationMessage(deathmatch, batch.AgentEntityId, mutation); err != nil {
		//						utils.Debug("arenaserver-mutation", err.Error()+"; coming from agent "+batch.AgentProxyUUID.String())
		//					}
		//				}
		//			}
		//		}

		// 2. Movement actions
		for _, mutation := range batch.Mutations {
			switch mutation.GetMethod() {
			case "steer":
				{
					if err := handleSteerMutationMessage(game, batch.AgentEntityId, mutation); err != nil {
						utils.Debug("arenaserver-mutation", err.Error()+"; coming from agent "+batch.AgentProxyUUID.String())
					}
				}
			}
		}

		// 3. if any debug, communicate them
		// for _, mutation := range batch.Mutations {
		// 	switch mutation.GetMethod() {
		// 	case "debugpoint":
		// 		{
		// 			handleDebugPointMutationMessage(deathmatch, batch.AgentEntityId, mutation)
		// 		}
		// 	}
		// }

	}
}

// func handleShootMutationMessage(deathmatch *SoccerGame, entityID ecs.EntityID, mutation types.AgentMessagePayloadActions) error {

// 	var aimingFloats []float64
// 	err := json.Unmarshal(mutation.GetArguments(), &aimingFloats)
// 	if err != nil {
// 		return errors.New("Failed to unmarshal JSON arguments for shoot mutation")
// 	}

// 	entityresult := deathmatch.getEntity(entityID, deathmatch.shootingComponent)
// 	if entityresult == nil {
// 		return errors.New("Failed to find entity associated to shoot mutation")
// 	}

// 	aiming := vector.MakeVector2(aimingFloats[0], aimingFloats[1]) //.
// 	//Transform(deathmatch.physicalToAgentSpaceInverseTransform)

// 	shootingAspect := entityresult.Components[deathmatch.shootingComponent].(*Shooting)
// 	shootingAspect.PushShot(aiming)

// 	return nil
// }

func handleSteerMutationMessage(game *SoccerGame, entityID ecs.EntityID, mutation types.AgentMessagePayloadActions) error {
	var steeringFloats []float64
	err := json.Unmarshal(mutation.GetArguments(), &steeringFloats)
	if err != nil {
		return errors.New("Failed to unmarshal JSON arguments for steer mutation")
	}

	entityresult := game.getEntity(entityID, game.steeringComponent)
	if entityresult == nil {
		return errors.New("Failed to find entity associated to steer mutation")
	}

	steering := vector.MakeVector2(steeringFloats[0], steeringFloats[1])

	steeringAspect := entityresult.Components[game.steeringComponent].(*Steering)
	steeringAspect.PushSteer(steering)

	return nil
}
