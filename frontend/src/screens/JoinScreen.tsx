import { sendMessage } from "../socket/socket";
import type { AppState } from "../state/reducer";
import { useState } from "react";

export default function JoinScreen({ state }: { state: AppState }) {
  const [name, setName] = useState("");

  const handleJoin = () => {
    if (!name.trim()) return;

    sendMessage({
      type: "JOIN_ROOM_REQUEST",
      payload: { playerName: name.trim() },
    });
  };

  return (
    <div
      className="min-h-screen flex items-center justify-center px-4
      bg-slate-50 text-slate-900
      dark:bg-zinc-900 dark:text-slate-100"
    >
      <div className="w-full max-w-md space-y-6">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold">GoWords</h1>
          <p className="text-sm text-slate-500 dark:text-slate-400">
            Competitive real-time word battles
          </p>
        </div>

        <div className="bg-white dark:bg-zinc-800 rounded-2xl shadow-lg p-6 space-y-5">
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Shadow🔥Fox"
            autoFocus
            autoComplete="off"
            spellCheck={false}
            className="w-full px-4 py-3 rounded-xl border
              border-slate-300 bg-slate-100
              focus:ring-2 focus:ring-blue-500
              dark:bg-zinc-700 dark:border-zinc-600"
          />

          {state.joinError && (
            <div className="text-sm text-red-500">{state.joinError}</div>
          )}

          <button
            onClick={handleJoin}
            disabled={!name.trim()}
            className="w-full py-3 rounded-xl font-semibold text-white
              bg-blue-600 hover:bg-blue-700
              active:scale-[0.98] transition
              disabled:opacity-50 disabled:cursor-not-allowed
              dark:bg-blue-500 dark:hover:bg-blue-400"
          >
            Join Game
          </button>
        </div>
      </div>
    </div>
  );
}
