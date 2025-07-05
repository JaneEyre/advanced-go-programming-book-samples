
# Save this file as client.py
import socket
import json
import sys

def main():
    """
    A Python client to connect to the Go JSON-RPC server.
    """
    # Define the server address
    HOST = 'localhost'
    PORT = 1234

    # The JSON-RPC request payload.
    # This structure must match what the server expects.
    # 'id' should ideally be unique for each request.
    request_payload = {
        "method": "HelloService.Hello",
        "params": ["Python"], # The argument to the Hello method
        "id": 101 # Using a different ID to show it's a new request
    }

    try:
        # Create a socket object (AF_INET for IPv4, SOCK_STREAM for TCP)
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
            print(f"Connecting to {HOST}:{PORT}...")
            # Connect to the server
            s.connect((HOST, PORT))
            print("Connected.")

            # Encode the Python dictionary to a JSON formatted string, then to bytes
            # We add a newline character at the end because many server-side
            # buffered readers read until a newline.
            json_request = json.dumps(request_payload).encode('utf-8') + b'\n'
            
            print(f"Sending: {json_request.decode().strip()}")
            # Send the data
            s.sendall(json_request)

            # Receive the response from the server (up to 1024 bytes)
            data = s.recv(1024)
            print("Response received.")

            # Decode the received bytes back into a JSON string and then parse it
            if data:
                response_payload = json.loads(data.decode('utf-8'))
                print("\n--- Server Response ---")
                print(f"ID: {response_payload.get('id')}")
                print(f"Result: {response_payload.get('result')}")
                print(f"Error: {response_payload.get('error')}")
                print("-----------------------")
            else:
                print("No data received from server.")

    except ConnectionRefusedError:
        print(f"Connection failed. Is the Go server running on port {PORT}?")
        sys.exit(1)
    except Exception as e:
        print(f"An error occurred: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
