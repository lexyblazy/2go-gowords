interface BasePayload {
  timestamp: number;
}

export interface RoomJoinOkEvent {
  type: "JOIN_ROOM_OK";
  payload: BasePayload & {
    systemMoniker: string;
    playerId: string;
    playerName: string;
    roomId: number;
  };
}

export interface RoomJoinErrorEvent {
  type: "JOIN_ROOM_ERROR";
  payload: BasePayload & {
    message: string;
  };
}

export interface GameRulesEvent {
  type: "GAME_RULES";
  payload: BasePayload & {
    message: string;
    rules: string[];
    systemMoniker: string;
  };
}

export interface RoundInfoEvent {
  type: "ROUND_INFO";
  payload: BasePayload & {
    systemMoniker: string;
    words: string[];
    validWordsCount: number;
    endsAt: number;
  };
}

export interface RoundOverEvent {
  type: "ROUND_OVER";
  payload: BasePayload & {
    systemMoniker: string;
  };
}

export interface RoundWinnerEvent {
  type: "ROUND_WINNER";
  payload: BasePayload & {
    systemMoniker: string;
    winnerPlayerName: string;
    winnerPlayerId: string;
    score: number;
  };
}

export interface PlayerWordSubmissionEvent {
  type: "PLAYER_WORD_SUBMISSION";
  payload: {
    playerId: string;
    word: string;
  };
}

export interface PlayerRoundScoresEvent {
  type: "PLAYER_ROUND_SCORES";
  payload: BasePayload & {
    systemMoniker: string;
    playerName: string;
    score: number;
  };
}

export interface PlayerWordAcceptedEvent {
  type: "PLAYER_WORD_ACCEPTED";
  payload: BasePayload & {
    systemMoniker: string;
    word: string;
    points: number;
  };
}

export interface PlayerWordRejectedEvent {
  type: "PLAYER_WORD_REJECTED";
  payload: BasePayload & {
    systemMoniker: string;
    word: string;
    message: string;
  };
}

export interface PlayerSubmissionBroadcastEvent {
  type: "PLAYER_SUBMISSION_BROADCAST";
  payload: BasePayload & {
    systemMoniker: string;
    playerName: string;
    word: string;
  };
}

export interface NextRoundCountdownEvent {
  type: "NEXT_ROUND_COUNTDOWN";
  payload: BasePayload & {
    endsAt: number;
    systemMoniker: string;
  };
}

export type ServerEvent =
  | GameRulesEvent
  | RoundInfoEvent
  | RoundOverEvent
  | RoundWinnerEvent
  | PlayerRoundScoresEvent
  | PlayerWordAcceptedEvent
  | PlayerWordRejectedEvent
  | PlayerSubmissionBroadcastEvent
  | NextRoundCountdownEvent
  | { type: "CONNECTED"; payload?: undefined }
  | { type: "DISCONNECTED"; payload?: undefined }
  | RoomJoinOkEvent
  | RoomJoinErrorEvent;

export interface FeedItem {
  id: string;
  timestamp: number;
  displayName: string;
  message: string;
  type?:
    | "system"
    | "wordAccepted"
    | "wordRejected"
    | "disconnected"
    | "otherPlayerSubmission"
    | "rules";
}

export interface RoundState {
  words: string[];
  validWordsCount: number;
  endsAt: number;
}
