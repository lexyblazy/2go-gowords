import React, { useState, useRef, useEffect } from "react";

interface Props {
  isRoundActive: boolean;
  sendWord: (word: string) => void;
}

export default React.memo(function WordInputBar({
  isRoundActive,
  sendWord,
}: Props) {
  const [word, setWord] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isRoundActive) {
      inputRef.current?.focus();
    }
  }, [isRoundActive]);

  const handleSubmit = () => {
    const trimmed = word.trim().toLowerCase();
    if (!trimmed || !isRoundActive) {
      return;
    }

    sendWord(trimmed);
    setWord("");
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleSubmit();
    }
  };

  return (
    <div
      className="
        border-t
        border-slate-200
        dark:border-zinc-700
        p-3
        bg-white
        dark:bg-zinc-900
      "
    >
      <div className="flex gap-2">
        <input
          ref={inputRef}
          value={word}
          onChange={(e) => setWord(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={!isRoundActive}
          placeholder={
            isRoundActive ? "Type a word..." : "Waiting for next round..."
          }
          autoCapitalize="none"
          autoCorrect="off"
          spellCheck={false}
          className="
            flex-1 rounded-md
            border border-slate-300
            bg-slate-100 focus:outline-none
            focus:ring-2 focus:ring-blue-500 focus:border-transparent
            dark:bg-zinc-800 dark:border-zinc-600 dark:focus:ring-blue-400
            disabled:opacity-50 px-3 py-2 text-sm pb-safe"
        />

        <button
          onClick={handleSubmit}
          onMouseDown={(e) => e.preventDefault()}
          disabled={!isRoundActive}
          className="
          px-5 py-3
          rounded-md
          font-semibold
          text-white
          bg-blue-600
          active:scale-95
          transition-transform
          duration-100
          shadow-md
          dark:bg-blue-500
          "
        >
          Send
        </button>
      </div>
    </div>
  );
});
