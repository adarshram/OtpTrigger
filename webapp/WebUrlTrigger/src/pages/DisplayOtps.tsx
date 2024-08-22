import { useEffect, useState } from "react";
import { auth } from "../BaseConfig";

import RefreshIcon from "@mui/icons-material/Refresh";
import { Alert, CircularProgress, Box, Button } from "@mui/material";
const apiUrl = import.meta.env.VITE_API_URL_BASE;
const fetchOtps = async () => {
  const user = auth.currentUser;
  if (user) {
    const idToken = await user.getIdToken();

    const fetchPromise = fetch(`${apiUrl}/otps`, {
      mode: "cors",
      method: "GET",
      headers: {
        Authorization: `Bearer ${idToken}`,
        "Content-Type": "application/json",
      },
    });
    return fetchPromise
      .then((response) => {
        try {
          return response.json();
        } catch (e) {
          throw new Error("Failed to parse response");
        }
      })
      .then((data) => {
        if (data.Success && data.Otps && data.Otps.length > 0) {
        }
        return data;
      });
  }
};

function DisplayOtps() {
  const [otps, setOtps] = useState([]);
  const [fetched, setFetched] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!fetched) {
      const fetchFunction = async () => {
        setError(null);
        try {
          const data = await fetchOtps();
          if (data.success && data.otps && data.otps.length > 0) {
            setOtps(data.otps);
            setFetched(true);
            return;
          }
          setOtps([]);
        } catch (e) {
          setOtps([]);
          setError("No otps found");
        }
        setFetched(true);
      };
      if (!fetched) {
        fetchFunction();
      }
    }
  }, [fetched]);

  useEffect(() => {
    const interval = setInterval(() => {
      setFetched(false);
    }, 15000);
    return () => clearInterval(interval);
  }, []);

  const fetchOtpOnClick = async () => {
    setFetched(false);
  };
  const getOtpFromText = (inputString: string): string => {
    let code = copyOtp(inputString);
    if (typeof code === "boolean") {
      return "";
    }
    return code;
  };

  const copyOtp = (inputString: string): string | boolean => {
    //933413 is the OTP to login to your ICICI direct account. It is valid for 5 mins. OTPs are SECRET. DO NOT disclose it to anyone. ICICI direct NEVER asks for OTP
    let codeRegex = /\s?(\d{6})(\s|\.)/;
    let codeMatch = inputString.match(codeRegex);

    if (codeMatch) {
      return codeMatch[1];
    }
    codeRegex = /\s(\d{5})(\s|\.)/;
    codeMatch = inputString.match(codeRegex);
    if (codeMatch) {
      return codeMatch[1];
    }

    codeRegex = /\s(\d{4})(\s|\.)/;
    codeMatch = inputString.match(codeRegex);
    if (codeMatch) {
      return codeMatch[1];
    }
    return false;
  };
  return (
    <div>
      <h1>Display OTPs</h1>
      <div>
        {error && <Alert severity="error">{error}</Alert>}
        {!fetched && !otps.length && (
          <div>
            <Alert severity="warning">No OTPs found</Alert>
          </div>
        )}
        {otps.map((otp, index) => {
          const canCopyOtp = copyOtp(otp);
          return (
            <Box
              key={index}
              sx={{
                userSelect: "text",
              }}
            >
              {otp}
              {canCopyOtp !== false && (
                <Button
                  onClick={() => {
                    navigator.clipboard.writeText(getOtpFromText(otp));
                  }}
                >
                  Copy {copyOtp(otp)}
                </Button>
              )}
            </Box>
          );
        })}
        <button onClick={fetchOtpOnClick}>
          <RefreshIcon />
        </button>
        {!fetched && (
          <div>
            <CircularProgress />
          </div>
        )}
      </div>
    </div>
  );
}
export default DisplayOtps;
