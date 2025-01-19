import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, Rating } from "@mui/material";
import { useState } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
  currentRating?: number;
  onRate: (value: number) => void;
  onRemove: () => void;
};

export const RatingModal = ({ open, onClose, currentRating = 0, onRate, onRemove }: Props) => {
  const [rating, setRating] = useState(currentRating);

  const handleRate = () => {
    onRate(rating);
    onClose();
  };

  const handleRemove = () => {
    onRemove();
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>Rate this book</DialogTitle>
      <DialogContent>
        <Box display="flex" alignItems="center" justifyContent="center" py={2}>
          <Rating
            value={rating}
            onChange={(_, value) => {
              setRating(value ?? 0);
            }}
            size="large"
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        {currentRating === rating && currentRating !== 0 ? (
          <Button onClick={handleRemove} color="error">
            Remove Rating
          </Button>
        ) : (
          <Button onClick={handleRate} variant="contained">
            Rate
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};
