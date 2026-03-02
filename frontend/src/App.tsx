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
    <>
      <Controls />
      {state.joinedRoom ? (
        <GameScreen state={state} />
      ) : (
        <JoinScreen state={state} />
      )}
    </>
  );
}
