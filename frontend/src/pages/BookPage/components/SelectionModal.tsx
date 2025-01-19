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
  Menu,
  MenuItem,
} from "@mui/material";
import CheckIcon from "@mui/icons-material/Check";
import { useState, useMemo } from "react";
import { Add, MoreVert, Edit, Delete } from "@mui/icons-material";

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
  onEdit?: (item: Item) => void;
  onDelete?: (item: Item) => void;
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
  onEdit,
  onDelete,
}: Props) => {
  const [search, setSearch] = useState("");
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedItem, setSelectedItem] = useState<Item | null>(null);

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

  const handleMenuClose = () => {
    setAnchorEl(null);
    setSelectedItem(null);
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
                    deleteIcon={<MoreVert />}
                    onDelete={(e: React.MouseEvent<HTMLElement>) => {
                      setAnchorEl(e.currentTarget);
                      setSelectedItem(item);
                    }}
                  />
                ))}
              </Box>
              <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={handleMenuClose}>
                {onEdit && (
                  <MenuItem
                    onClick={() => {
                      if (selectedItem) onEdit(selectedItem);
                      handleMenuClose();
                    }}
                  >
                    <Edit sx={{ mr: 1 }} /> Edit
                  </MenuItem>
                )}
                {onDelete && (
                  <MenuItem
                    onClick={() => {
                      if (selectedItem) onDelete(selectedItem);
                      handleMenuClose();
                    }}
                    sx={{ color: "error.main" }}
                  >
                    <Delete sx={{ mr: 1 }} /> Delete
                  </MenuItem>
                )}
              </Menu>
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
