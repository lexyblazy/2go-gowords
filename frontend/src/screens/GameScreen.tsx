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
    <div className="flex flex-col min-h-dvh w-full max-w-full overflow-x-hidden bg-white dark:bg-zinc-900 text-slate-900 dark:text-slate-100">
      {/* 🖥 Desktop LetterBoard (Top) */}
      <div className="hidden md:block shrink-0 border-b border-slate-200 dark:border-zinc-700 p-4">
        <LetterBoard round={round} nextRoundStartsAt={nextRoundStartsAt} />
      </div>

      {/* 🧾 Scrollable Feed */}
      <div className="flex-1 overflow-y-auto px-4 py-2">
        <ActivityFeed feedItems={feedItems} />
      </div>

      {/* 📱 Mobile LetterBoard (Above Input) */}
      <div className="block md:hidden shrink-0 border-t border-slate-200 dark:border-zinc-700 p-3 bg-white dark:bg-zinc-900">
        <LetterBoard round={round} nextRoundStartsAt={nextRoundStartsAt} />
      </div>

      {/* ⌨️ Input */}
      <div className="shrink-0">
        <WordInputBar isRoundActive={!!round} sendWord={handleSendWord} />
      </div>
    </div>
  );
}
