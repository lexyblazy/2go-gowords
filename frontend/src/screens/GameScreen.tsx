import LetterBoard from "../components/LetterBoard";
import ScoreBar from "../components/ScoreBar";
import ActivityFeed from "../components/ActivityFeed";
import WordInputBar from "../components/WordInputBar";
import type { AppState } from "../state/reducer";
import { sendMessage } from "../socket/socket";

export default function GameScreen({ state }: { state: AppState }) {
  return (
    <div className="flex flex-col h-screen bg-white dark:bg-zinc-900 text-slate-900 dark:text-slate-100">
      {/* 🔒 Static Top Section */}
      <div className="border-b border-slate-200 dark:border-zinc-700 p-4 space-y-3">
        <LetterBoard round={state.round} />
        <ScoreBar score={state.playerScore} />
      </div>

      {/* 🧾 Scrollable Feed */}
      <div className="flex-1 overflow-y-auto p-4">
        <ActivityFeed feed={state.feed} />
      </div>

      {/* ⌨️ Input */}
      <WordInputBar
        isRoundActive={state.isRoundActive}
        sendWord={(word) =>
          sendMessage({
            type: "PLAYER_WORD_SUBMISSION",
            payload: {
              playerId: state.playerId ?? "",
              word,
            },
          })
        }
      />
    </div>
  );
}
