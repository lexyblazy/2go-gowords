import { useState } from "react";
import type { FeedItem } from "../state/types";

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp);
  // format time to show AM/PM
  return date.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
};

function getItemClass(item: FeedItem) {
  switch (item.type) {
    case "system":
      return "text-center text-xs text-slate-400 dark:text-slate-500";

    case "wordAccepted":
      return "text-center text-green-600 dark:text-green-400";

    case "wordRejected":
      return "text-center text-red-500";

    case "disconnected":
      return "text-center text-red-500 bg-red-500/10 border border-red-500/20 p-2 rounded-lg";

    default:
      return "";
  }
}

export default function ActivityFeedItem({ item }: { item: FeedItem }) {
  const [showRules, setShowRules] = useState(false);

  if (item.type === "otherPlayerSubmission") {
    return (
      <div className="flex items-baseline gap-2 text-sm animate-fade-up">
        <span className="text-xs text-slate-400 w-12 shrink-0 text-[11px] opacity-40">
          {formatTime(item.timestamp)}
        </span>

        <span className="font-semibold text-slate-800 dark:text-slate-200">
          {item.displayName}:
        </span>

        <span className="text-slate-700 dark:text-slate-300 font-medium">
          {item.message}
        </span>
      </div>
    );
  }

  if (item.type === "rules") {
    return (
      <div className="px-2 py-2 animate-fade-up">
        <div
          className="
            max-w-2xl mx-auto text-sm rounded-xl
            bg-slate-100 dark:bg-zinc-800 text-slate-600 dark:text-slate-400
            leading-relaxed p-3"
        >
          <div className="font-semibold mb-2 text-slate-800 dark:text-slate-200">
            <button onClick={() => setShowRules((prev) => !prev)}>ℹ️ Game Rules {showRules ? "▲" : "▼"}</button>
          </div>
          {showRules && <div className="rounded-md bg-slate-100/60 dark:bg-zinc-800/60
          text-xs leading-snug text-slate-600 dark:text-slate-400">{item.message}</div>}
        </div>
      </div>
    );
  }

  return (
    <div className={`${getItemClass(item)} animate-fade-up`}>
      <div className="text-center">
        <span className="font-mono text-gray-500 text-[11px] opacity-40">
          {formatTime(item.timestamp)}
        </span>{" "}
        <span className="text-sm font-semibold">{item.displayName}:</span>{" "}
        <span className="text-sm text-center">{item.message}</span>
      </div>
    </div>
  );
}
