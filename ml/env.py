import threading
import json
import numpy as np
from tf_agents.typing import types
from tf_agents.environments import py_environment
from tf_agents.specs import array_spec
from tf_agents.trajectories import time_step as ts
from go.game import game, match, bot, ml, context, json as gojson

class BattlerGameEnv(py_environment.PyEnvironment):

    def __init__(self, c: game.CharacterData, opp: game.CharacterData, bot: match.Player, format: ml.Format, state_format: types.NestedArraySpec):
        super().__init__()
        
        self._action_spec = array_spec.BoundedArraySpec(
            shape=(), dtype=np.int64, minimum=0, maximum=3, name='action'
        )

        self._format = format
        self._observation_spec = state_format
        
        self._c_data = c
        self._opp_data = opp
        self._bot = bot

        self._episode_ended = True
    
    def action_spec(self) -> types.NestedArraySpec:
        return self._action_spec

    def observation_spec(self) -> types.NestedArraySpec:
        return self._observation_spec
    
    def _reset(self) -> ts.TimeStep:
        self._init_game()
        return ts.restart(self._state)
    
    def _init_game(self):
        self._adapter = bot.NewAdapter()

        self._c = game.NewCharacter(self._c_data)
        self._opp = game.NewCharacter(self._opp_data)

        p1 = match.CharacterPlayer(self._c, self._adapter)
        p2 = match.CharacterPlayer(self._opp, self._bot)
        self._m = match.New(p1, p2, True)

        threading.Thread(target=self._run_match, args=()).start()

        state = self._adapter.GetStateInit()
        res = self._format.Row(state)
        self._state = json.loads(bytes(gojson.Marshal(res)))
        self._state = {key: np.array([val], dtype=self._observation_spec[key].dtype) for key, val in self._state.items()}

        self._reward = 0
        self._episode_ended = False
    
    def _run_match(self):
        self._m.Run(context.Background())
    
    def _step(self, action: types.NestedArray) -> ts.TimeStep:
        if self._episode_ended:
            self.reset()

        try:
            state = self._adapter.GetState(int(action))
            res = self._format.Row(state)
            self._state = json.loads(bytes(gojson.Marshal(res)))
            self._state = {key: np.array([val], dtype=self._observation_spec[key].dtype) for key, val in self._state.items()}

            reward = self._c.HP() - self._opp.HP()
            reward -= self._reward
            self._reward += reward

            if state.IsEnd():
                diff = self._c.HP() - self._opp.HP()
                if diff > 0:
                    reward += 100
                elif diff < 0:
                    reward -= 100

                # print("END!!!!", self._state, reward)
                self._episode_ended = True
                return ts.termination(self._state, reward)
            else:
                # print("not end :(", self._state, reward)
                return ts.transition(self._state, reward)

        except:
            self._reward -= 10
            return ts.termination(self._state, -10)

        