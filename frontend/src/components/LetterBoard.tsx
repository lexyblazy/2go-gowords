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
      <>
        <div className="flex flex-wrap justify-center gap-3">
          {[...Array(10)].map((_, i) => (
            <div
              key={i}
              className="w-12 h-12 bg-slate-100 dark:bg-zinc-800 rounded-xl animate-pulse"
            ></div>
          ))}
        </div>
        <div className="text-center">

        <CountdownTimer endsAt={nextRoundStartsAt ?? 0} />
        </div>
      </>
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
      <CountdownTimer endsAt={round.endsAt} />
    </div>
  );
}
