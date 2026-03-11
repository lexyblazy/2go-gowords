import type { FeedItem, RoundState, ServerEvent, BatchEvents } from "./types";
import { applyBatchEvents, getFeedItem, getNewState } from "./methods";
export interface AppState {
  connectionStatus: "connecting" | "connected" | "disconnected";
  playerName?: string;
  playerId?: string;
  joinError?: string;
  round?: RoundState;
  joinedRoom: boolean;
  feed: FeedItem[];
  isRoundActive: boolean;
  playerScore: number;
  nextRoundStartsAt?: number;
}

export const initialState: AppState = {
  connectionStatus: "connecting",
  joinedRoom: false,
  feed: [],
  isRoundActive: false,
  playerScore: 0,
};

export function reducer(
  state: AppState,
  event: ServerEvent | BatchEvents,
): AppState {
  switch (event.type) {
    case "CONNECTED":
      return { ...state, connectionStatus: "connected" };

    case "JOIN_ROOM_OK":
      return {
        ...state,
        joinedRoom: true,
        playerId: event.payload.playerId,
        playerName: event.payload.playerName,
        joinError: undefined,
        isRoundActive: true,
        feed: [...state.feed, getFeedItem(event)!],
      };

    case "JOIN_ROOM_ERROR":
      return {
        ...state,
        joinError: event.payload.message,
      };

    case "DISCONNECTED":
      return {
        ...state,
        isRoundActive: false,
        connectionStatus: "disconnected",
        feed: [...state.feed, getFeedItem(event)!],
      };
    // batch events are used to reduce the number of re-renders
    case "BATCH_EVENTS":
      return applyBatchEvents(state, event.payload);

    default:
      return getNewState(state, event);
  }
}
