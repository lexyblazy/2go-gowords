import type { ServerEvent } from "../state/types";
const WS_BASE = import.meta.env.VITE_WS_BASE_URL;

type Dispatch = (event: ServerEvent) => void;

let socket: WebSocket | null = null;

export function initSocket(dispatch: Dispatch) {
  const WS_URL = WS_BASE
    ? `${WS_BASE}/ws`
    : `${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/ws`;

  socket = new WebSocket(WS_URL);

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
