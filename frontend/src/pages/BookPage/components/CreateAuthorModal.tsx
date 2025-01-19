import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button, Stack } from "@mui/material";
import { useState } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
  onCreate: (author: { firstName: string; secondName: string; middleName?: string; nickName?: string }) => void;
};

export const CreateAuthorModal = ({ open, onClose, onCreate }: Props) => {
  const [form, setForm] = useState({
    firstName: "",
    secondName: "",
    middleName: "",
    nickName: "",
  });

  const handleCreate = () => {
    if (!form.firstName || !form.secondName) {
      return;
    }

    onCreate({
      firstName: form.firstName,
      secondName: form.secondName,
      ...(form.middleName && { middleName: form.middleName }),
      ...(form.nickName && { nickName: form.nickName }),
    });

    setForm({ firstName: "", secondName: "", middleName: "", nickName: "" });
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Create Author</DialogTitle>
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
        <Button onClick={handleCreate} variant="contained" disabled={!form.firstName || !form.secondName}>
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
};
