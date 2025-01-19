import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Stack } from "@mui/material";
import { useState, useEffect } from "react";
import { Author } from "@shared/types/author";

type Props = {
  open: boolean;
  onClose: () => void;
  author: Author;
  onEdit: (
    authorId: string,
    data: { firstName: string; secondName: string; middleName?: string; nickName?: string },
  ) => void;
};

export const EditAuthorModal = ({ open, onClose, author, onEdit }: Props) => {
  const [form, setForm] = useState({
    firstName: author.firstName,
    secondName: author.secondName,
    middleName: author.middleName ?? "",
    nickName: author.nickname ?? "",
  });

  useEffect(() => {
    setForm({
      firstName: author.firstName,
      secondName: author.secondName,
      middleName: author.middleName ?? "",
      nickName: author.nickname ?? "",
    });
  }, [author]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Edit Author</DialogTitle>
      <DialogContent>
        <Stack spacing={2} sx={{ mt: 1 }}>
          <TextField
            label="First Name"
            required
            value={form.firstName}
            onChange={(e) => {
              setForm((prev) => ({ ...prev, firstName: e.target.value }));
            }}
          />
          <TextField
            label="Last Name"
            required
            value={form.secondName}
            onChange={(e) => {
              setForm((prev) => ({ ...prev, secondName: e.target.value }));
            }}
          />
          <TextField
            label="Middle Name"
            value={form.middleName}
            onChange={(e) => {
              setForm((prev) => ({ ...prev, middleName: e.target.value }));
            }}
          />
          <TextField
            label="Nickname"
            value={form.nickName}
            onChange={(e) => {
              setForm((prev) => ({ ...prev, nickName: e.target.value }));
            }}
          />
        </Stack>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={() => {
            onEdit(author.id, form);
          }}
          variant="contained"
          disabled={!form.firstName || !form.secondName}
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};
