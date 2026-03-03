import LetterBoard from "../components/LetterBoard";
import ActivityFeed from "../components/ActivityFeed";
import WordInputBar from "../components/WordInputBar";
import { sendMessage } from "../socket/socket";
import type { FeedItem, RoundState } from "../state/types";
import { useCallback } from "react";

interface Props {
  round?: RoundState;
  feedItems: FeedItem[];
  nextRoundStartsAt?: number;
  playerId?: string;
}

export default function GameScreen({
  round,
  feedItems,
  nextRoundStartsAt,
  playerId,
}: Props) {
  const handleSendWord = useCallback(
    (word: string) => {
      sendMessage({
        type: "PLAYER_WORD_SUBMISSION",
        payload: {
          playerId,
          word,
        },
      });
    },
    [playerId],
  );

  return (
    <div className="flex flex-col h-screen bg-white dark:bg-zinc-900 text-slate-900 dark:text-slate-100">
      {/* 🔒 Static Top Section */}
      <div className="border-b border-slate-200 dark:border-zinc-700 p-4 space-y-3">
        <LetterBoard round={round} nextRoundStartsAt={nextRoundStartsAt} />
      </div>

      {/* 🧾 Scrollable Feed */}
      <div className="flex-1 overflow-y-auto p-4">
        <ActivityFeed feedItems={feedItems} />
      </div>

      {/* ⌨️ Input */}
      <WordInputBar isRoundActive={!!round} sendWord={handleSendWord} />
    </div>
  );
}
