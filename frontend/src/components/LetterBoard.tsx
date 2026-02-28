import type { RoundState } from "../state/types";

export default function LetterBoard({
  round,
}: {
  round: RoundState | undefined;
}) {
  if (!round) {
    return (
      <div className="text-center text-slate-500 dark:text-slate-400">
        Loading...
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
        {round.validWordsCount} possible valid words
      </div>
    </div>
  );
}
