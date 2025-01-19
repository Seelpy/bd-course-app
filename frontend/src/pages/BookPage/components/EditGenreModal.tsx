import { Dialog, DialogTitle, DialogContent, DialogActions, TextField, Button } from "@mui/material";
import { useState, useEffect } from "react";
import { Genre } from "@shared/types/genre";

type Props = {
  open: boolean;
  onClose: () => void;
  genre: Genre;
  onEdit: (genreId: string, name: string) => void;
};

export const EditGenreModal = ({ open, onClose, genre, onEdit }: Props) => {
  const [name, setName] = useState(genre.name);

  useEffect(() => {
    setName(genre.name);
  }, [genre]);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Edit Genre</DialogTitle>
      <DialogContent>
        <TextField
          fullWidth
          label="Name"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
          }}
          sx={{ mt: 2 }}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={() => {
            onEdit(genre.id, name);
          }}
          variant="contained"
          disabled={!name.trim()}
        >
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
};
