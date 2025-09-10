import { useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
} from "@/components/ui/card";
import { LogIn } from "lucide-react";

const GOOGLE_CLIENT_ID = import.meta.env.VITE_GOOGLE_CLIENT_ID || "";
const GOOGLE_REDIRECT_URI = import.meta.env.VITE_GOOGLE_REDIRECT_URI || "";

export default function GoogleLoginPage() {
  const todoAuthToken = localStorage.getItem("todo-auth-token");

  useEffect(() => {
    if (todoAuthToken) {
      window.location.href = "/";
    }
  }, [todoAuthToken]);

  const handleGoogleLogin = () => {
    // Construct Google OAuth URL for authorization code flow
    // Use environment variable for redirect URI or fallback to current origin
    const redirectUri =
      GOOGLE_REDIRECT_URI || window.location.origin + "/google";
    const scope = "openid email profile";
    const responseType = "code";
    const state = Math.random().toString(36).substring(2, 15);

    const googleAuthUrl =
      `https://accounts.google.com/o/oauth2/v2/auth?` +
      `client_id=${GOOGLE_CLIENT_ID}&` +
      `redirect_uri=${encodeURIComponent(redirectUri)}&` +
      `scope=${encodeURIComponent(scope)}&` +
      `response_type=${responseType}&` +
      `state=${state}`;

    // Store state in sessionStorage for verification
    sessionStorage.setItem("oauth_state", state);

    // Redirect to Google OAuth
    window.location.href = googleAuthUrl;
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 via-white to-slate-50 flex items-center justify-center p-6">
      <div className="w-full max-w-md">
        {/* Header Section */}
        <div className="text-center mb-12">
          <div className="flex justify-center mb-8">
            <div className="w-16 h-16 bg-slate-900 rounded-2xl flex items-center justify-center shadow-sm">
              <LogIn className="h-8 w-8 text-white" />
            </div>
          </div>
          <h1 className="text-3xl font-light text-slate-900 mb-3 tracking-tight">
            Welcome
          </h1>
          <p className="text-slate-500 font-light">
            Sign in to access your tasks
          </p>
        </div>

        {/* Login Card */}
        <Card className="border-0 shadow-sm bg-white/70 backdrop-blur-sm">
          <CardContent className="p-8">
            <Button
              variant="outline"
              className="w-full h-14 text-slate-700 border-slate-200 hover:bg-slate-50 hover:border-slate-300 font-medium transition-all duration-200 shadow-sm hover:shadow-md"
              onClick={handleGoogleLogin}
            >
              <svg className="w-5 h-5 mr-3" viewBox="0 0 24 24">
                <path
                  fill="#4285F4"
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                />
                <path
                  fill="#34A853"
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                />
                <path
                  fill="#FBBC05"
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                />
                <path
                  fill="#EA4335"
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                />
              </svg>
              Continue with Google
            </Button>
          </CardContent>
        </Card>

        {/* Footer */}
        <div className="text-center mt-8">
          <p className="text-xs text-slate-400 font-light">
            Secure authentication powered by Google
          </p>
        </div>
      </div>
    </div>
  );
}
