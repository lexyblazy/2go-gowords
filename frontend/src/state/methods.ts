import type { AppState } from "./reducer";
import type {
  FeedItem,
  GameRulesEvent,
  RoundInfoEvent,
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

// const makeRoundInfoFeedItem = (event: RoundInfoEvent): FeedItem => {
//   const message = `The words are: ${event.payload.words.join(" ")}. There are ${event.payload.validWordsCount} possible valid words.`;
//   return {
//     id: crypto.randomUUID(),
//     timestamp: event.payload.timestamp,
//     displayName: event.payload.systemMoniker,
//     message,
//     type: "system",
//   };
// };

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

const makeRoundWinnerFeedItem = (event: RoundWinnerEvent): FeedItem => {
  const message = `🏆🥇 Kudos to ${event.payload.winnerPlayerName} for winning the round with ${event.payload.score} points! 🏆🥇`;
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
  return {
    id: crypto.randomUUID(),
    timestamp: event.payload.timestamp,
    displayName: event.payload.systemMoniker,
    message: `Next round in ${event.payload.roundIntervalSeconds} seconds`,
    type: "system",
  };
};

const makeJoinRoomOkFeedItem = (event: RoomJoinOkEvent): FeedItem => {
  const message = `You have joined Room ${event.payload.roomId} as ${event.payload.playerName}`;

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
  const message = `🎉 You scored ${event.payload.points} points for the word "${event.payload.word}"! 🎉`;
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
    message: `❌😞 "${event.payload.word}" was rejected. ${event.payload.message} 😞❌`,
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

export function getFeedItem(event: ServerEvent): FeedItem | undefined {
  switch (event.type) {
    case "GAME_RULES":
      return makeGameRulesFeedItem(event);
    case "ROUND_OVER":
      return makeRoundOverFeedItem(event);
    case "ROUND_WINNER":
      return makeRoundWinnerFeedItem(event);
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
    default:
      console.log("Unknown event: ", event);
      return undefined;
  }
}

export function getNewState(state: AppState, event: ServerEvent): AppState {
  const newFeedItem = getFeedItem(event);

  const newState: AppState = {
    ...state,
    feed: newFeedItem ? [...state.feed, newFeedItem] : state.feed,
  };

  if (event.type === "ROUND_INFO") {
    newState.round = {
      words: event.payload.words,
      validWordsCount: event.payload.validWordsCount,
    };
  } else if (event.type === "ROUND_OVER") {
    newState.round = undefined;
  }

  return newState;
}
