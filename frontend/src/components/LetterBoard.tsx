import { useEffect, useState } from "react";
import type { RoundState } from "../state/types";

export default function LetterBoard({
  round,
}: {
  round: RoundState | undefined;
}) {
  const [remaining, setRemaining] = useState(0);

  useEffect(() => {
    if (!round?.endsAt) return;

    const interval = setInterval(() => {
      const seconds = Math.max(
        0,
        Math.floor((round.endsAt - Date.now()) / 1000),
      );
      setRemaining(seconds);
    }, 1000);

    return () => clearInterval(interval);
  }, [round?.endsAt]);

  if (!round) {
    return (
      <div className="flex flex-wrap justify-center gap-3">
        {[...Array(10)].map((_, i) => (
          <div
            key={i}
            className="w-12 h-12 bg-slate-100 dark:bg-zinc-800 rounded-xl animate-pulse"
          ></div>
        ))}
      </div>
    );
  }

  const letters = round.words.join("").split("");

  return (
    <div className="flex flex-col items-center gap-3">
      <div className="flex flex-wrap justify-center gap-3">
        {letters.map((letter, i) => (
          <div
            key={i}
            className="
              w-12 h-12
              md:w-14 md:h-14
              flex items-center justify-center
              rounded-xl
              bg-slate-100
              dark:bg-zinc-800
              text-xl md:text-2xl
              font-bold
              shadow-sm
            "
          >
            {letter.toUpperCase()}
          </div>
        ))}
      </div>

      <div className="text-sm text-slate-500 dark:text-slate-400">
        There are {round.validWordsCount} possible valid words
      </div>
      <div
        className={
          remaining <= 10
            ? "text-red-500 font-semibold animate-pulse"
            : "text-slate-500 dark:text-slate-400"
        }
      >
        Ends in {remaining} seconds
      </div>
    </div>
  );
}
