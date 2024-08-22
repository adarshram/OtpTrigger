import { useEffect, useState } from "react";
import "./App.css";
import { auth, googleProvider, authLogout } from "./BaseConfig";
import { signInWithPopup } from "firebase/auth";
import DisplayOtps from "./pages/DisplayOtps";
type User = {
  uid: string | null;
  displayName: string | null;
  email: string | null;
};
function App() {
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  const signInWithGoogle = async () => {
    try {
      await signInWithPopup(auth, googleProvider);
    } catch (error) {
      console.error("Error signing in with Google:", error);
    }
  };
  const logoutFromFirebase = async () => {
    try {
      await authLogout();
    } catch (error) {
      console.error("Error logging out:", error);
    }
  };
  useEffect(() => {
    auth.onAuthStateChanged((user) => {
      setCurrentUser(user);
    });
    return () => {
      //auth.onAuthStateChanged(null);
    };
  }, []);

  return (
    <>
      {!currentUser || !currentUser.uid ? (
        <button onClick={signInWithGoogle}>Sign in with Google</button>
      ) : (
        <>
          <button onClick={logoutFromFirebase}>Logout</button>
          Welcome{" "}
          {currentUser?.displayName
            ? currentUser?.displayName
            : currentUser?.email ?? "No Name"}
          <DisplayOtps />
        </>
      )}
    </>
  );
}

export default App;
