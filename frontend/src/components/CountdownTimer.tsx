import { useEffect, useState } from "react";
import { soundManager } from "../lib/sound";

export default function CountdownTimer({ endsAt }: { endsAt: number }) {
  const [now, setNow] = useState(Date.now());

  useEffect(() => {
    const interval = setInterval(() => {
      setNow(Date.now());
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const remaining = Math.max(0, Math.floor((endsAt - now) / 1000));

    if (remaining <= 10 && remaining > 0) {
      soundManager.play("beep");
    } else {
      soundManager.stop("beep");
    }
  }, [endsAt, now]);

  function formatTime(seconds: number) {
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes.toString().padStart(2, "0")}:${remainingSeconds.toString().padStart(2, "0")}`;
  }

  const remaining = Math.max(0, Math.floor((endsAt - now) / 1000));

  return remaining > 0 ? (
    <div
      className={
        remaining <= 10
          ? "text-red-500 font-bold animate-pulse"
          : "text-slate-500 font-semibold dark:text-slate-400"
      }
    >
      {formatTime(remaining)}
    </div>
  ) : null;
}
