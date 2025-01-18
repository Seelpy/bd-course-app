import { useEffect, useState } from "react";
import { Container, Paper, Typography, TextField, Button, Box, Grid2, Card, CardContent } from "@mui/material";
import { bookApi } from "@api/book";
import { verifyBookRequestApi } from "@api/verifyBookRequest";
import { useSnackbar } from "notistack";
import { VerifyBookRequest } from "@shared/types/verifyBookRequest";

export function RequestsPage() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [requests, setRequests] = useState<VerifyBookRequest[]>([]);
  const { enqueueSnackbar } = useSnackbar();

  const loadRequests = () => {
    verifyBookRequestApi
      .listVerifyBookRequests()
      .then((data) => {
        setRequests(data.verifyBookRequests);
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  useEffect(() => {
    loadRequests();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      await bookApi.createBook({
        title,
        description,
      });
      enqueueSnackbar("Book verification request created successfully", {
        variant: "success",
      });
      setTitle("");
      setDescription("");
      loadRequests();
    } catch (error) {
      enqueueSnackbar((error as Error).message, { variant: "error" });
    }
  };

  return (
    <Container>
      <Paper elevation={3} sx={{ p: 3, borderRadius: 4, width: "100%" }}>
        <Typography variant="h5" gutterBottom>
          Missing book? Create a request here
        </Typography>
        <Box component="form" onSubmit={handleSubmit}>
          <Grid2 container spacing={2}>
            <Grid2 size={{ xs: 12 }}>
              <TextField
                fullWidth
                label="Title"
                value={title}
                onChange={(e) => {
                  setTitle(e.target.value);
                }}
                required
              />
            </Grid2>
            <Grid2 size={{ xs: 12 }}>
              <TextField
                fullWidth
                label="Description"
                multiline
                rows={4}
                value={description}
                onChange={(e) => {
                  setDescription(e.target.value);
                }}
                required
              />
            </Grid2>
            <Grid2>
              <Button variant="contained" type="submit">
                Create Request
              </Button>
            </Grid2>
          </Grid2>
        </Box>
      </Paper>

      <Paper elevation={3} sx={{ p: 3, borderRadius: 4, width: "100%", mt: 2 }}>
        <Typography variant="h5" gutterBottom>
          Verification Requests
        </Typography>
        <Grid2 container spacing={2}>
          {requests.map((request) => (
            <Grid2 size={{ xs: 12 }} key={request.verifyBookRequestId}>
              <Card elevation={2}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    Book ID: {request.bookId}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Translator ID: {request.translatorId}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Status: {request.isVerified ? "Verified" : "Pending"}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Sent: {new Date(request.sendDateMilli).toLocaleDateString()}
                  </Typography>
                </CardContent>
              </Card>
            </Grid2>
          ))}
        </Grid2>
      </Paper>
    </Container>
  );
}
