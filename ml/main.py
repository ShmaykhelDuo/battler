import json
import torch
import socketserver

device = torch.device(
    "cuda" if torch.cuda.is_available() else
    "mps" if torch.backends.mps.is_available() else
    "cpu"
)

module = torch.jit.load("ml/model-Ruby-vs-Milana-2.pt").to(device)

class Handler(socketserver.BaseRequestHandler):

    def handle(self):
        while True:
            msg = self.request.recv(1024)
            if len(msg) == 0:
                raise RuntimeError("socket connection broken")
            
            state = json.loads(msg)
            if type(state) is not dict:
                raise RuntimeError(f"invalid message: {msg.decode()}")
            
            state = torch.tensor(state["state"], dtype=torch.float32, device=device).unsqueeze(0)
            # action = module(state).max(1).indices.view(1, 1).item()
            actions = module(state).squeeze().tolist()

            print(actions)
            msg = json.dumps({"actions": actions}).encode()
            self.request.sendall(msg)


if __name__ == "__main__":
    HOST, PORT = "localhost", 9999

    with socketserver.TCPServer((HOST, PORT), Handler) as server:
        # Activate the server; this will keep running until you
        # interrupt the program with Ctrl-C
        server.serve_forever()
