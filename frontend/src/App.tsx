import { useEffect, useReducer } from "react";
import { reducer, initialState } from "./state/reducer";
import { initSocket } from "./socket/socket";
import Controls from "./components/Controls";
import JoinScreen from "./screens/JoinScreen";
import GameScreen from "./screens/GameScreen";

export default function App() {
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    initSocket(dispatch);
  }, []);

  return (
    <div className="min-h-dvh w-full max-w-full min-w-0 overflow-x-hidden flex flex-col">
      <Controls />
      <main className="flex-1 min-w-0 flex flex-col min-h-0">
        {state.joinedRoom ? (
          <GameScreen
            round={state.round}
            feedItems={state.feed}
            nextRoundStartsAt={state.nextRoundStartsAt}
            playerId={state.playerId}
          />
        ) : (
          <JoinScreen joinError={state.joinError} />
        )}
      </main>
    </div>
  );
}
