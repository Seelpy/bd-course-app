import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Stack } from "@mui/material";
import { useState } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
  onSave: (data: { title?: string; description?: string }) => void;
  title?: string;
  description?: string;
  mode: "title" | "description";
};

export const EditBookDialog = ({ open, onClose, onSave, title = "", description = "", mode }: Props) => {
  const [value, setValue] = useState(mode === "title" ? title : description);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{mode === "title" ? "Edit Title" : "Edit Description"}</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          {mode === "title" ? (
            <TextField
              label="Title"
              fullWidth
              value={value}
              onChange={(e) => {
                setValue(e.target.value);
              }}
              required
            />
          ) : (
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              value={value}
              onChange={(e) => {
                setValue(e.target.value);
              }}
              required
            />
          )}
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={() => {
            onSave(mode === "title" ? { title: value } : { description: value });
            onClose();
          }}
          variant="contained"
          disabled={!value.trim()}
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};
