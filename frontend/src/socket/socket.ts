import type { ServerEvent } from "../state/types";

type Dispatch = (event: ServerEvent) => void;

let socket: WebSocket | null = null;

export function initSocket(dispatch: Dispatch) {
  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    dispatch({ type: "CONNECTED" });
  };

  socket.onmessage = (event) => {
    const parsed: ServerEvent = JSON.parse(event.data);
    dispatch(parsed);
  };

  socket.onclose = () => {
    dispatch({ type: "DISCONNECTED" });
  };
}

export function sendMessage(message: unknown) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    console.error("Socket is not open");
    return;
  }
  socket.send(JSON.stringify(message));
}
