import {
  Dialog,
  DialogTitle,
  DialogContent,
  Box,
  Chip,
  CircularProgress,
  TextField,
  Button,
  Stack,
} from "@mui/material";
import CheckIcon from "@mui/icons-material/Check";
import { useState, useMemo } from "react";
import { Add } from "@mui/icons-material";

type Item = {
  id: string;
  name: string;
};

type Props = {
  open: boolean;
  onClose: () => void;
  title: string;
  items: Item[];
  selectedIds: string[];
  onSelect: (id: string) => void;
  onCreate?: (name: string) => void;
  loading?: boolean;
  searchPlaceholder?: string;
  createButtonText?: string;
  forceCreateNewButton?: boolean;
};

export const SelectionModal = ({
  open,
  onClose,
  title,
  items,
  selectedIds,
  onSelect,
  onCreate,
  loading,
  searchPlaceholder = "Search...",
  createButtonText = "Create New",
  forceCreateNewButton = false,
}: Props) => {
  const [search, setSearch] = useState("");

  const filteredItems = useMemo(
    () => items.filter((item) => item.name.toLowerCase().includes(search.toLowerCase())),
    [items, search],
  );

  const showCreateButton =
    ((search.length > 0 && !filteredItems.some((item) => search === item.name)) || forceCreateNewButton) &&
    onCreate &&
    !loading;

  const handleCreate = () => {
    if (onCreate) {
      onCreate(search);
      setSearch("");
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>
        <Stack spacing={2}>
          <TextField
            fullWidth
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
            }}
            placeholder={searchPlaceholder}
            size="small"
            sx={{ mt: 1 }}
          />
          {loading ? (
            <Box display="flex" justifyContent="center" p={2}>
              <CircularProgress />
            </Box>
          ) : (
            <>
              <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                {filteredItems.map((item) => (
                  <Chip
                    key={item.id}
                    label={item.name}
                    onClick={() => {
                      onSelect(item.id);
                    }}
                    icon={selectedIds.includes(item.id) ? <CheckIcon /> : undefined}
                    color={selectedIds.includes(item.id) ? "primary" : "default"}
                  />
                ))}
              </Box>
              {showCreateButton && (
                <Button startIcon={<Add />} variant="outlined" onClick={handleCreate}>
                  {createButtonText} {!forceCreateNewButton && search}
                </Button>
              )}
            </>
          )}
        </Stack>
      </DialogContent>
    </Dialog>
  );
};
