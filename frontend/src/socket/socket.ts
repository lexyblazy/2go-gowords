import type { BatchEvents, ServerEvent } from "../state/types";

const WS_BASE = import.meta.env.VITE_WS_BASE_URL;

type Dispatch = (event: ServerEvent | BatchEvents) => void;

const UNBATCHABLE_EVENTS = new Set([
  "JOIN_ROOM_OK",
  "JOIN_ROOM_ERROR",
  "DISCONNECTED",
]);

const MAX_BATCHED_EVENTS = 2000;

let socket: WebSocket | null = null;
let eventQueue: ServerEvent[] = [];
let scheduled = false;
let rafId: number | null = null;
let timeoutId: number | null = null;

function clearScheduledWork() {
  if (rafId !== null) {
    cancelAnimationFrame(rafId);
    rafId = null;
  }

  if (timeoutId !== null) {
    clearTimeout(timeoutId);
    timeoutId = null;
  }

  scheduled = false;
}

function flush(dispatch: Dispatch) {
  if (!scheduled) {
    return;
  }

  clearScheduledWork();

  if (eventQueue.length === 0) {
    return;
  }

  const batch = eventQueue;
  eventQueue = [];

  dispatch({
    type: "BATCH_EVENTS",
    payload: batch,
  });
}

function scheduleFlush(dispatch: Dispatch) {
  if (scheduled) {
    return;
  }

  scheduled = true;

  rafId = requestAnimationFrame(() => flush(dispatch));
  timeoutId = window.setTimeout(() => flush(dispatch), 16);
}

function cleanupSocketInstance(ws?: WebSocket | null) {
  const target = ws ?? socket;

  clearScheduledWork();
  eventQueue = [];

  if (!target) {
    return;
  }

  target.onopen = null;
  target.onmessage = null;
  target.onclose = null;
  target.onerror = null;

  if (
    target.readyState === WebSocket.OPEN ||
    target.readyState === WebSocket.CONNECTING
  ) {
    target.close();
  }

  if (target === socket) {
    socket = null;
  }
}

export function initSocket(dispatch: Dispatch) {
  cleanupSocketInstance();

  const WS_URL = WS_BASE
    ? `${WS_BASE}/ws`
    : `${window.location.protocol === "https:" ? "wss" : "ws"}://${window.location.host}/ws`;

  const ws = new WebSocket(WS_URL);

  socket = ws;

  ws.onopen = () => {
    if (socket !== ws) {
      return;
    }

    dispatch({ type: "CONNECTED" });
  };

  ws.onmessage = (event) => {
    if (socket !== ws) {
      return;
    }

    const parsed: ServerEvent = JSON.parse(event.data);

    if (UNBATCHABLE_EVENTS.has(parsed.type)) {
      dispatch(parsed);
      return;
    }

    if (eventQueue.length >= MAX_BATCHED_EVENTS) {
      eventQueue.shift();
    }

    eventQueue.push(parsed);
    scheduleFlush(dispatch);
  };

  ws.onclose = () => {
    if (socket !== ws) {
      return;
    }

    clearScheduledWork();
    socket = null;
    dispatch({ type: "DISCONNECTED" });
  };

  return () => {
    if (socket === ws) {
      cleanupSocketInstance(ws);
    } else {
      cleanupSocketInstance(ws);
    }
  };
}

export function closeSocket() {
  cleanupSocketInstance();
}

export function sendMessage(message: unknown) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    console.error("Socket is not open");
    return;
  }

  socket.send(JSON.stringify(message));
}
