import random
from env import GameEnv

env = GameEnv("/tmp/test.sock")

for i in range(10000):
    state = env.get_state()
    print("Got state:", state)

    if state["end"]:
        env.reset()
        continue

    action = random.randint(0, 3)
    print("Sending action:", action)
    env.send_action(action)
