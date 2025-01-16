import { Navigate } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";

export const RootPage = () => {
  // TODO: Implement the root page instead of redirecting to the 404 page
  return <Navigate to={AppRoute.NotFound} replace />;
};
