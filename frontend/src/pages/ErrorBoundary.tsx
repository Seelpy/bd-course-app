import { Box, Button, Typography } from "@mui/material";

export default function ErrorBoundary() {
  return (
    <Box
      display="flex"
      justifyContent="center"
      alignItems="center"
      minHeight="100vh"
      flexDirection="column"
      sx={{ textAlign: "center", p: 2 }}
    >
      <Typography variant="h5" gutterBottom>
        Страница не найдена
      </Typography>
      <Typography variant="body1" gutterBottom>
        К сожалению, страница, которую вы пытаетесь открыть, не существует.
      </Typography>
      <Button
        variant="contained"
        color="primary"
        onClick={() => {
          window.history.back();
        }}
        sx={{ mt: 2 }}
      >
        Вернуться назад
      </Button>
    </Box>
  );
}
