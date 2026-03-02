import { useContext, useState } from "react";
import { ThemeContext } from "../theme/ThemeProvider";
import { soundManager } from "../lib/sound";

export default function Controls() {
  const { theme, toggle } = useContext(ThemeContext);
  const [soundOn, setSoundOn] = useState(false);

  const toggleSound = () => {
    const enabled = soundManager.toggle();
    setSoundOn(enabled);
  };

  return (
    <div className="fixed top-4 right-4 flex flex-col gap-2 z-10">
      <button
        onClick={toggleSound}
        className="p-2 rounded-lg bg-slate-200 dark:bg-zinc-800 opacity-70 hover:opacity-100"
      >
        {soundOn ? "🔊" : "🔇"}
      </button> 
      <button
        onClick={toggle}
        className="p-2 rounded-lg bg-slate-200 dark:bg-zinc-800"
      >
        {theme === "dark" ? "☀️" : "🌙"}
      </button>
    </div>
  );
}
