import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Box, Button, Container, TextField, Typography } from "@mui/material";
import { AuthForm } from "@shared/types/auth";
import { authApi } from "@api/auth";
import { AppRoute } from "@shared/constants/routes";
import { useSnackbar } from "notistack";
import { useUserStore } from "@shared/stores/userStore";
import { useShallow } from "zustand/shallow";

export const LoginPage = () => {
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const { setUserInfo } = useUserStore(
    useShallow((state) => ({
      setUserInfo: state.setUserInfo,
    })),
  );

  const [form, setForm] = useState<AuthForm>({ login: "", password: "" });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    authApi
      .login(form)
      .then(() => {
        enqueueSnackbar("Login successful", { variant: "success" });
        setUserInfo({ id: "test", name: form.login, email: "mock@gmail.com" });
        navigate(AppRoute.Root);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  return (
    <Container maxWidth="xs">
      <Box
        sx={{
          marginTop: 8,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
        }}
      >
        <Typography component="h1" variant="h5">
          Sign in
        </Typography>
        <Box component="form" onSubmit={handleSubmit} sx={{ mt: 1 }}>
          <TextField
            margin="normal"
            required
            fullWidth
            label="Login"
            value={form.login}
            onChange={(e) => {
              setForm({ ...form, login: e.target.value });
            }}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="Password"
            type="password"
            value={form.password}
            onChange={(e) => {
              setForm({ ...form, password: e.target.value });
            }}
          />
          <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }}>
            Sign In
          </Button>
          <Button
            fullWidth
            onClick={() => {
              navigate(AppRoute.Register);
            }}
          >
            Don&apos;t have an account? Sign Up
          </Button>
        </Box>
      </Box>
    </Container>
  );
};
