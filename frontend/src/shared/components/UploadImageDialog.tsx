import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Box } from "@mui/material";
import { ChangeEvent, useState } from "react";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";

type Props = {
  open: boolean;
  onClose: () => void;
  onUpload: (base64: string) => void;
  title: string;
};

export const UploadImageDialog = ({ open, onClose, onUpload, title }: Props) => {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = () => {
        const base64 = reader.result as string;
        setPreviewUrl(base64);
      };
      reader.readAsDataURL(file);
    }
  };

  const handleUpload = () => {
    if (previewUrl) {
      onUpload(previewUrl);
      setPreviewUrl(null);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Box display="flex" flexDirection="column" alignItems="center" gap={2} py={2}>
          {previewUrl ? (
            <img
              src={previewUrl}
              alt="Preview"
              style={{ maxWidth: "100%", maxHeight: "300px", objectFit: "contain" }}
            />
          ) : (
            <Button component="label" variant="outlined" startIcon={<CloudUploadIcon />} sx={{ marginTop: 1 }}>
              Choose File
              <input type="file" hidden accept="image/*" onChange={handleFileChange} />
            </Button>
          )}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleUpload} disabled={!previewUrl} variant="contained">
          Upload
        </Button>
      </DialogActions>
    </Dialog>
  );
};
