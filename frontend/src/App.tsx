import { useEffect, useReducer } from "react";
import { reducer, initialState } from "./state/reducer";
import { initSocket } from "./socket/socket";
import JoinScreen from "./screens/JoinScreen";
import GameScreen from "./screens/GameScreen";

export default function App() {
  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    initSocket(dispatch);
  }, []);

  if (!state.joinedRoom) {
    return <JoinScreen state={state} />;
  }

  return <GameScreen state={state} />;
}
