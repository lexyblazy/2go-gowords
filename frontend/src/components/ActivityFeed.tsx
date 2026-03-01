import { useEffect, useRef } from "react";
import type { FeedItem } from "../state/types";

interface Props {
  feed: FeedItem[];
}

const formatTime = (timestamp: number) => {
  const date = new Date(timestamp);
  // format time to show AM/PM
  return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" });
};

export default function ActivityFeed({ feed }: Props) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "auto" });
  }, [feed]);

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

  return (
    <div
      ref={containerRef}
      className="flex-1 overflow-y-auto p-4 space-y-1.5 h-full"
    >
      {feed.map((item) => {
        if (item.type === "otherPlayerSubmission") {
          return (
            <div className="flex items-baseline gap-2 text-sm" key={item.id}>
              <span className="text-xs text-slate-400 w-12 shrink-0">
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
            <div className="px-4 py-3" key={item.id}>
              <div
                className="
    max-w-2xl mx-auto
    rounded-xl
    bg-slate-100 dark:bg-zinc-800
    p-4
    text-sm
    text-slate-600 dark:text-slate-400
    leading-relaxed
  "
              >
                <div className="font-semibold mb-2 text-slate-800 dark:text-slate-200">
                  ℹ️ Game Rules
                </div>
                {item.message}
              </div>
            </div>
          );
        }

        return (
          <div key={item.id} className={getItemClass(item)}>
            <div className="text-center">
              <span className="font-mono text-xs text-gray-500">
                {formatTime(item.timestamp)}
              </span>{" "}
              <span className="font-semibold">{item.displayName}</span>
            </div>

            <div className="text-sm text-center">{item.message}</div>
          </div>
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
}
