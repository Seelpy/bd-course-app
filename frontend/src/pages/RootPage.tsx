import { Navigate } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";

export const RootPage = () => {
  return <Navigate to={AppRoute.NotFound} replace />;
};
