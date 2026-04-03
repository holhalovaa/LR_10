import asyncio
import websockets
import sys

async def send_messages(websocket):
    while True:
        msg = await asyncio.get_event_loop().run_in_executor(None, sys.stdin.readline)
        await websocket.send(msg.strip())

async def receive_messages(websocket):
    while True:
        msg = await websocket.recv()
        print(f"\n[Новое сообщение]: {msg}")

async def main():
    uri = "ws://localhost:8083/ws"
    async with websockets.connect(uri) as websocket:
        print("Подключено к чату. Введите сообщение:")
        await asyncio.gather(send_messages(websocket), receive_messages(websocket))

asyncio.run(main())