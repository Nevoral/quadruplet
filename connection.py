import socket
import time
import json

GET = "GET"
POST = "POST"
RESET = "RESET"


class BSDSocket:
    def __init__(self):
        self.address = 'localhost'
        self.port = 42069
        self.retry_interval = 5
        self.max_retries = 10

    def OpenSocket(self):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        attempts = 0

        while attempts < self.max_retries:
            try:
                print(f'Attempt {attempts + 1}: Connecting to {self.address} port {self.port}')
                sock.connect((self.address, self.port))
                return sock
            except ConnectionRefusedError:
                print(
                    f'Connection to {self.address} port {self.port} failed. Retrying in {self.retry_interval} seconds.')
                time.sleep(self.retry_interval)
                attempts += 1
                if attempts == self.max_retries:
                    print('Maximum retry attempts reached. Exiting.')
                    raise ConnectionError("Maximum retry attempts reached.")


class Message:
    def __init__(self, s: BSDSocket, method: str = GET):
        self.connection = s
        self.method = method
        self.actions = []
        self.body = object
        self.response = ""

    def CreateMessage(self) -> str:
        if self.method == POST:
            return f"<<<method:{self.method};actions:{self.GetActions()}><{json.dumps(self.body)}>>>\n"
        return f"<<<method:{self.method};actions:{self.GetActions()}>>>\n"

    def Action(self, *action: str) -> None:
        for a in action:
            if not a.lower() in self.actions:
                self.actions.append(a.lower())

    def GetActions(self) -> str:
        actions = ""
        for i in self.actions:
            actions += f"{i};"
        return actions

    def SendMessage(self):
        sock = self.connection.OpenSocket()
        try:
            # Send data
            m = self.CreateMessage()
            print(f'Sending:\n\t{m}')
            sock.sendall(m.encode())

            # Wait for messages from the server and print them
            while True:
                data = sock.recv(1024).decode()
                if data[0] == "<" and data[1] == "<":
                    self.response += data[2:]
                elif data[-3:-1] == ">>":
                    self.response += data[:-3]
                    self.response = json.loads(self.response)
                    print(f'Received:\n\t{self.response}')
                    break
                else:
                    self.response += data
        finally:
            print('Closing socket')
            sock.close()


def main():
    s: BSDSocket = BSDSocket()
    l = Message(s, GET)

    # You will get all your points
    l.Action("getallpoints")
    l.SendMessage()
    k = Message(s, POST)
    # you can manipulate with those points by yourself just uplouad same model that you have get but with different
    # coords.
    k.Action("setrobot")
    k.body = l.response
    k.body["Legs"][0]["Name"] = "Frog Left Leg"
    k.SendMessage()


if __name__ == '__main__':
    main()
