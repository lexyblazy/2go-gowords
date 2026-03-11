import { soundManager } from "../lib/sound";
import type { AppState } from "./reducer";
import type {
  FeedItem,
  GameRulesEvent,
  RoundOverEvent,
  RoundWinnerEvent,
  ServerEvent,
  PlayerSubmissionBroadcastEvent,
  NextRoundCountdownEvent,
  RoomJoinOkEvent,
  PlayerWordAcceptedEvent,
  PlayerWordRejectedEvent,
  PlayerRoundScoresEvent,
} from "./types";

const MAX_FEED_ITEMS = 500;

const addDriftTime = (endsAt: number) => {
  const randomDriftTime = Math.random() * 1500;
  return endsAt + randomDriftTime;
};

export const playSound = (event: ServerEvent) => {
  switch (event.type) {
    case "PLAYER_WORD_ACCEPTED":
      soundManager.play("accepted");
      break;
    case "PLAYER_WORD_REJECTED":
      soundManager.play("rejected");
      break;
    case "ROUND_OVER":
      soundManager.play("over");
      break;
    case "GAME_RULES":
      soundManager.play("default");
      break;
    default:
      break;
  }
};

const makeGameRulesFeedItem = (event: GameRulesEvent): FeedItem => {
  const message = event.payload.rules.join("\n");

  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message,
    type: "rules",
  };
};

const makeRoundOverFeedItem = (event: RoundOverEvent): FeedItem => {
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker ?? "",
    message: `The round is over`,
    type: "system",
  };
};

const makeRoundWinnerFeedItem = (
  event: RoundWinnerEvent,
  playerId: string | undefined,
): FeedItem => {
  let message: string = `🏆 Kudos to ${event.payload.winnerPlayerName} for winning the round with ${event.payload.score} points! 🏆`;

  if (playerId === event.payload.winnerPlayerId) {
    message = `🥇🥇 You won the round with ${event.payload.score} points! 🥇🥇`;
    soundManager.stop("over");
    soundManager.play("winner");
  }
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message,
    type: "system",
  };
};

const makePlayerSubmissionBroadcastFeedItem = (
  event: PlayerSubmissionBroadcastEvent,
): FeedItem => {
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.playerName,
    message: event.payload.word,
    type: "otherPlayerSubmission",
  };
};

const makeNextRoundCountdownFeedItem = (
  event: NextRoundCountdownEvent,
): FeedItem => {
  const remaining = Math.max(
    0,
    Math.floor((addDriftTime(event.payload.endsAt) - Date.now()) / 1000),
  );

  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message: `Next round in ${remaining} seconds`,
    type: "system",
  };
};

const makeJoinRoomOkFeedItem = (event: RoomJoinOkEvent): FeedItem => {
  const message = `You have joined Room #${event.payload.roomId} as ${event.payload.playerName}`;

  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message,
    type: "system",
  };
};

const makePlayerWordAcceptedFeedItem = (
  event: PlayerWordAcceptedEvent,
): FeedItem => {
  const message = ` +${event.payload.points} "${event.payload.word}" 🎉`;
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message,
    type: "wordAccepted",
  };
};

const makePlayerWordRejectedFeedItem = (
  event: PlayerWordRejectedEvent,
): FeedItem => {
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message: `❌ "${event.payload.word}" rejected. ${event.payload.message}`,
    type: "wordRejected",
  };
};

const makePlayerRoundScoresFeedItem = (
  event: PlayerRoundScoresEvent,
): FeedItem => {
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message: `🎊 You scored ${event.payload.score} points for the round! 🎊`,
    type: "system",
  };
};

const makeDisconnectedFeedItem = (event: ServerEvent): FeedItem => {
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload?.timestamp ?? Date.now(),
    displayName: "Game Master",
    message: `‼️⛔️ You have been disconnected from the game. Reload the page to join a new game. ‼️⛔️`,
    type: "disconnected",
  };
};

export function getFeedItem(
  event: ServerEvent,
  playerId?: string,
): FeedItem | undefined {
  switch (event.type) {
    case "GAME_RULES":
      return makeGameRulesFeedItem(event);
    case "ROUND_OVER":
      return makeRoundOverFeedItem(event);
    case "ROUND_WINNER":
      return makeRoundWinnerFeedItem(event, playerId);
    case "PLAYER_SUBMISSION_BROADCAST":
      return makePlayerSubmissionBroadcastFeedItem(event);
    case "NEXT_ROUND_COUNTDOWN":
      return makeNextRoundCountdownFeedItem(event);
    case "JOIN_ROOM_OK":
      return makeJoinRoomOkFeedItem(event);
    case "PLAYER_WORD_ACCEPTED":
      return makePlayerWordAcceptedFeedItem(event);
    case "PLAYER_WORD_REJECTED":
      return makePlayerWordRejectedFeedItem(event);
    case "PLAYER_ROUND_SCORES":
      return makePlayerRoundScoresFeedItem(event);
    case "DISCONNECTED":
      return makeDisconnectedFeedItem(event);
    case "ROUND_INFO":
      return;
    default:
      return;
  }
}

export function getNewState(state: AppState, event: ServerEvent): AppState {
  playSound(event);

  const newFeedItem = getFeedItem(event, state.playerId);

  const newState: AppState = {
    ...state,
    feed: (newFeedItem ? [...state.feed, newFeedItem] : state.feed).slice(
      -MAX_FEED_ITEMS,
    ),
  };

  if (event.type === "ROUND_INFO") {
    newState.round = {
      words: event.payload.words,
      validWordsCount: event.payload.validWordsCount,
      endsAt: addDriftTime(event.payload.endsAt),
    };
  } else if (event.type === "ROUND_OVER") {
    newState.round = undefined;
  } else if (event.type === "NEXT_ROUND_COUNTDOWN") {
    newState.nextRoundStartsAt = addDriftTime(event.payload.endsAt);
    newState.round = undefined;
  }

  return newState;
}

export function applyBatchEvents(
  state: AppState,
  events: ServerEvent[],
): AppState {
  let playerId = state.playerId;
  let round = state.round;
  let nextRoundStartsAt = state.nextRoundStartsAt;
  let isRoundActive = state.isRoundActive;
  let connectionStatus = state.connectionStatus;
  let joinedRoom = state.joinedRoom;
  let playerName = state.playerName;
  let joinError = state.joinError;
  let playerScore = state.playerScore;

  const feed = state.feed.slice();

  for (const event of events) {
    playSound(event);

    const item = getFeedItem(event, playerId);
    if (item) {
      feed.push(item);
      if (feed.length > MAX_FEED_ITEMS) {
        feed.splice(0, feed.length - MAX_FEED_ITEMS);
      }
    }

    switch (event.type) {
      case "ROUND_INFO":
        round = {
          words: event.payload.words,
          validWordsCount: event.payload.validWordsCount,
          endsAt: addDriftTime(event.payload.endsAt),
        };
        break;
      case "ROUND_OVER":
        round = undefined;
        break;
      case "NEXT_ROUND_COUNTDOWN":
        nextRoundStartsAt = addDriftTime(event.payload.endsAt);
        round = undefined;
        break;
      default:
        break;
    }
  }

  return {
    ...state,
    connectionStatus,
    joinedRoom,
    playerId,
    playerName,
    joinError,
    isRoundActive,
    playerScore,
    round,
    nextRoundStartsAt,
    feed,
  };
}
