import random
from env import GameEnv

env = GameEnv("/tmp/test.sock")

for i in range(100):
    state = env.get_state()
    print("Got state:", state)

    if state["end"]:
        env.client.close()
        break

    action = random.randint(0, 3)
    print("Sending action:", action)
    env.send_action(action)
