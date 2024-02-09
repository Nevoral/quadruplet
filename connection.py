import socket


def sendMessage(message):
    # Create a TCP/IP socket
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Connect the socket to the server
    server_address = ('localhost', 42069)
    print(f'Connecting to {server_address[0]} port {server_address[1]}')
    sock.connect(server_address)

    try:
        # Send data
        print(f'Sending: {message}')
        sock.sendall(message.encode())

        print(sock.recv(54).decode())
    finally:
        print('Closing socket')
        sock.close()


def main():
    sendMessage('This is the message. It will be echoed.')
    sendMessage('fuck this shit')
    sendMessage('ale jede to \n kurva \n')


if __name__ == '__main__':
    main()
