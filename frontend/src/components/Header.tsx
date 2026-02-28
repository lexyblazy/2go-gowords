import { useContext } from "react";
import { ThemeContext } from "../theme/ThemeProvider";

export default function Header() {
  const { theme, toggle } = useContext(ThemeContext);

  return (
    <div className="flex justify-between items-center p-4">
      <h1 className="font-bold">GoWords</h1>
      <button
        onClick={toggle}
        className="p-2 rounded-lg bg-slate-200 dark:bg-zinc-800"
      >
        {theme === "dark" ? "☀️" : "🌙"}
      </button>
    </div>
  );
}
