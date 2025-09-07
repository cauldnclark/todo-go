import { createBrowserRouter } from "react-router";
import TodoListPage from "./routes/home";
import GoogleLoginPage from "./routes/login";
import { GoogleCallback } from "./routes/callback";

const router = createBrowserRouter([
  {
    path: "/",
    element: <TodoListPage />,
  },
  {
    path: "/login",
    element: <GoogleLoginPage />,
  },
  {
    path: "/google",
    element: <GoogleCallback />,
  },
]);

export default router;
