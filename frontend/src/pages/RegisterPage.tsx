import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Box, Button, Container, TextField, Typography } from "@mui/material";
import { userApi } from "@api/user";
import { CreateUser } from "@shared/types/user";
import { AppRoute } from "@shared/constants/routes";
import { useSnackbar } from "notistack";

type RegisterFormData = CreateUser & { confirmPassword: string };

export const RegisterPage = () => {
  const navigate = useNavigate();
  const { enqueueSnackbar } = useSnackbar();

  const [form, setForm] = useState<RegisterFormData>({
    login: "",
    password: "",
    aboutMe: "",
    confirmPassword: "",
  });
  const [passwordsMismatch, setPasswordsMismatch] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const { confirmPassword, ...createUserData } = form;

    if (createUserData.password !== confirmPassword) {
      setPasswordsMismatch(true);
      return;
    }

    userApi
      .createUser(createUserData)
      .then(() => {
        enqueueSnackbar("Register successful", { variant: "success" });
        navigate(AppRoute.Login);
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
          Sign up
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
            error={!!passwordsMismatch}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            label="Confirm Password"
            type="password"
            value={form.confirmPassword}
            onChange={(e) => {
              setForm({ ...form, confirmPassword: e.target.value });
            }}
            error={!!passwordsMismatch}
            helperText={passwordsMismatch && "Passwords do not match"}
          />
          <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }}>
            Sign Up
          </Button>
          <Button
            fullWidth
            onClick={() => {
              navigate(AppRoute.Login);
            }}
          >
            Already have an account? Sign in
          </Button>
        </Box>
      </Box>
    </Container>
  );
};
