import type { RoundState } from "../state/types";
import CountdownTimer from "./CountdownTimer";

export default function LetterBoard({
  round,
  nextRoundStartsAt,
}: {
  round: RoundState | undefined;
  nextRoundStartsAt: number | undefined;
}) {
  if (!round) {
    return (
      <div className="w-full max-w-full min-w-0 overflow-x-hidden">
        <div className="w-full max-w-full min-w-0 flex flex-wrap justify-center gap-2">
          {[...Array(10)].map((_, i) => (
            <div
              key={i}
              className="w-12 h-12 shrink-0 bg-slate-100 dark:bg-zinc-800 rounded-xl animate-pulse"
            ></div>
          ))}
        </div>
        <div className="text-center">
          <CountdownTimer endsAt={nextRoundStartsAt ?? 0} />
        </div>
      </div>
    );
  }

  const letters = round.words.join("").split("");

  return (
    <div className="w-full max-w-full min-w-0 overflow-x-hidden flex flex-col items-center gap-2">
      <div className="w-full max-w-full min-w-0 flex flex-wrap justify-center gap-2">
        {letters.map((letter, i) => (
          <div
            key={i}
            className="
              w-8 h-8 text-sm
              sm:w-10 sm:h-10 sm:text-base
              md:w-12 md:h-12 md:text-2xl
              flex items-center justify-center
              rounded-md
              bg-slate-200 dark:bg-zinc-800
              font-bold"
          >
            {letter.toUpperCase()}
          </div>
        ))}
      </div>

      <div className="text-sm text-slate-500 dark:text-slate-400 min-w-0">
        There are {round.validWordsCount} possible valid words
      </div>
      <CountdownTimer endsAt={round.endsAt} />
    </div>
  );
}
