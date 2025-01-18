import { useEffect, useState } from "react";
import {
  Container,
  Paper,
  Typography,
  TextField,
  Button,
  Box,
  Grid2,
  Card,
  CardContent,
  IconButton,
  Divider,
} from "@mui/material";
import DeleteIcon from "@mui/icons-material/Delete";
import { bookApi } from "@api/book";
import { verifyBookRequestApi } from "@api/verifyBookRequest";
import { useSnackbar } from "notistack";
import { VerifyBookRequest } from "@shared/types/verifyBookRequest";
import { useUserStore } from "@shared/stores/userStore";
import { useShallow } from "zustand/shallow";
import { UserRole } from "@shared/types/user";

export function RequestsPage() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [requests, setRequests] = useState<VerifyBookRequest[]>([]);
  const { enqueueSnackbar } = useSnackbar();

  const { userInfo } = useUserStore(
    useShallow((state) => ({
      userInfo: state.userInfo,
    })),
  );

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

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    bookApi
      .createBook({
        title,
        description,
      })
      .then(() => {
        enqueueSnackbar("Book verification request created successfully", {
          variant: "success",
        });
        setTitle("");
        setDescription("");
        loadRequests();
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleAcceptRequest = (requestId: string, accept: boolean) => {
    verifyBookRequestApi
      .acceptVerifyBookRequest({
        verifyBookRequestId: requestId,
        accept,
      })
      .then(() => {
        enqueueSnackbar(`Request ${accept ? "accepted" : "declined"} successfully`, {
          variant: "success",
        });
        loadRequests();
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  const handleDeleteRequest = (requestId: string) => {
    verifyBookRequestApi
      .deleteVerifyBookRequest({
        verifyBookRequestId: requestId,
      })
      .then(() => {
        enqueueSnackbar("Request deleted successfully", { variant: "success" });
        loadRequests();
      })
      .catch((error: Error) => {
        enqueueSnackbar(error.message, { variant: "error" });
      });
  };

  return (
    <Container>
      {userInfo && (
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
      )}

      <Paper elevation={3} sx={{ p: 3, borderRadius: 4, width: "100%", mt: 2 }}>
        <Typography variant="h5" gutterBottom>
          Verification Requests
        </Typography>
        <Grid2 container spacing={2}>
          {requests.map((request) => (
            <Grid2 size={{ xs: 12 }} key={request.verifyBookRequestId}>
              <Card elevation={2}>
                <Box position="relative">
                  {request.isVerified === undefined &&
                    (userInfo?.role === UserRole.Admin || request.translatorId === userInfo?.id) && (
                      <IconButton
                        size="small"
                        sx={{
                          position: "absolute",
                          right: 8,
                          top: 8,
                        }}
                        onClick={() => {
                          handleDeleteRequest(request.verifyBookRequestId);
                        }}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    )}
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {request.book.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {request.book.description}
                    </Typography>
                    <Divider sx={{ marginY: 1 }} />
                    <Typography variant="body2" color="text.secondary">
                      Status:{" "}
                      {request.isVerified === undefined ? (
                        "Pending"
                      ) : request.isVerified ? (
                        <Typography component="span" variant="body2" color="success.main">
                          Verified
                        </Typography>
                      ) : (
                        <Typography component="span" variant="body2" color="error.main">
                          Declined
                        </Typography>
                      )}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Sent: {new Date(request.sendDateMilli).toLocaleDateString()}
                    </Typography>
                  </CardContent>
                  {userInfo?.role === UserRole.Admin && request.isVerified === undefined && (
                    <Box
                      sx={{
                        position: "absolute",
                        bottom: 8,
                        right: 8,
                      }}
                    >
                      <Button
                        color="error"
                        onClick={() => {
                          handleAcceptRequest(request.verifyBookRequestId, false);
                        }}
                      >
                        Decline
                      </Button>
                      <Button
                        color="success"
                        onClick={() => {
                          handleAcceptRequest(request.verifyBookRequestId, true);
                        }}
                      >
                        Accept
                      </Button>
                    </Box>
                  )}
                </Box>
              </Card>
            </Grid2>
          ))}
        </Grid2>
      </Paper>
    </Container>
  );
}
