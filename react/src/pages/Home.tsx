import { useState } from "react";
import { useAuth0, type GetTokenSilentlyOptions, type IdToken } from "@auth0/auth0-react";
import { useEffect } from "react";


export const HomePage = () => {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [idTokenClaims, setIdTokenClaims] = useState<IdToken | undefined>(undefined);
  const { getAccessTokenSilently, getIdTokenClaims, loginWithRedirect } = useAuth0();

  useEffect(() => {
    const getAccessToken = async () => {

      try {
        const options: GetTokenSilentlyOptions = {
          authorizationParams: {
            audience: "https://api.trade.gustavobrunetto.com",
          }
        }
        const accessToken = await getAccessTokenSilently(options);
        setAccessToken(accessToken);
        const claims = await getIdTokenClaims();
        const response = await fetch("http://localhost:6060/api/messages/admin", {
          headers: {
            "content-type": "application/json",
            Authorization: `Bearer ${accessToken}`,
          },
        });
        if (!response.ok) {
          throw new Error("Failed to fetch protected resource");
        }
        const data = await response.json();
        console.log("Protected resource data:", data);
        setIdTokenClaims(claims);
      } catch (error: any) {
        console.error("Error getting access token:", error);

        // Handle consent required error
        if (error.error === 'consent_required' || error.error === 'interaction_required') {
          // Redirect to login with consent
          loginWithRedirect({
            authorizationParams: {
              audience: "https://api.trade.gustavobrunetto.com",
              prompt: "consent"
            }
          });
        }

        return null;
      }
    };

    getAccessToken();
  }, [getAccessTokenSilently, getIdTokenClaims, loginWithRedirect]);

  return (
    <div className="page-layout">
      <div className="page-layout__content">
        <h1>Welcome to the Home Page</h1>
        <p>This is a protected route. You must be logged in to see this content.</p>
        {accessToken ? (
          <p>Your access token is: {accessToken}</p>
        ) : (
          <p>You are not authenticated.</p>
        )}
      </div>
      <div>
        <h2>ID Token Claims</h2>
        {idTokenClaims ? (
          <pre>{JSON.stringify(idTokenClaims, null, 2)}</pre>
        ) : (
          <p>No ID token claims available.</p>
        )}
      </div>
    </div>
  );
};

export default HomePage;