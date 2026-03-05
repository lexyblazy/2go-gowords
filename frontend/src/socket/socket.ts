import type { BatchEvents, ServerEvent } from "../state/types";
const WS_BASE = import.meta.env.VITE_WS_BASE_URL;

const eventQueue: ServerEvent[] = [];
let scheduled = false;

type Dispatch = (event: ServerEvent | BatchEvents) => void;

let socket: WebSocket | null = null;

function flush(dispatch: Dispatch) {
  if (!scheduled) {
    return;
  }

  scheduled = false;

  if (eventQueue.length > 0) {
    dispatch({
      type: "BATCH_EVENTS",
      payload: eventQueue.splice(0),
    });
  }
}

function scheduleFlush(dispatch: Dispatch) {
  if (scheduled) {
    return;
  }
  scheduled = true;

  requestAnimationFrame(() => flush(dispatch));
  setTimeout(() => flush(dispatch), 100); // fallback when tab inactive
}

export function initSocket(dispatch: Dispatch) {
  const WS_URL = WS_BASE
    ? `${WS_BASE}/ws`
    : `${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/ws`;

  socket = new WebSocket(WS_URL);
  // socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    dispatch({ type: "CONNECTED" });
  };

  socket.onmessage = (event) => {
    const parsed: ServerEvent = JSON.parse(event.data);

    // some events are not batchable, so we dispatch them immediately
    const UNBATCHABLE_EVENTS = [
      "JOIN_ROOM_OK",
      "JOIN_ROOM_ERROR",
      "DISCONNECTED",
    ];

    if (UNBATCHABLE_EVENTS.includes(parsed.type)) {
      dispatch(parsed);
      return;
    }

    eventQueue.push(parsed);
    scheduleFlush(dispatch);
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
