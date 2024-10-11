import socket
import json


class GameEnv:

    def __init__(self, socket_path: str) -> None:
        self.socket_path = socket_path
        self.client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.client.connect(socket_path)
    
    def reset(self) -> None:
        self.client.close()

        self.client = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.client.connect(self.socket_path)

    def get_state(self) -> dict:
        msg = self.client.recv(1024)
        res = json.loads(msg)
        if len(res) == 0:
            raise RuntimeError("socket connection broken")
        if type(res) is not dict:
            raise RuntimeError(f"invalid message: {msg.decode()}")
        return res

    def send_action(self, action: int) -> None:
        msg = json.dumps({"action": action}).encode()
        totalsent = 0
        while totalsent < len(msg):
            sent = self.client.send(msg[totalsent:])
            if sent == 0:
                raise RuntimeError("socket connection broken")
            totalsent += sent
