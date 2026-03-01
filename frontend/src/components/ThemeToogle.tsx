import { useContext } from "react";
import { ThemeContext } from "../theme/ThemeProvider";

export default function ThemeToggle() {
  const { theme, toggle } = useContext(ThemeContext);

  return (
    <button
      onClick={toggle}
      className="p-2 rounded-lg bg-slate-200 dark:bg-zinc-800 fixed top-4 right-4 z-10"
    >
      {theme === "dark" ? "☀️" : "🌙"}
    </button>
  );
}
