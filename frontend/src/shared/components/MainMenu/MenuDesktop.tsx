import { useState } from "react";
import { AppBar, Box, Toolbar, Typography, Button, Menu, MenuItem, Avatar } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { AppRoute } from "@shared/constants/routes";
import { User } from "@shared/types/user";
import ThemeButton from "../ThemeButton/ThemeButton";

type MenuDesktopProps = {
  userInfo: User | null;
  handleLogout: () => void;
  menuItems: {
    text: string;
    icon: React.ReactNode;
    onClick: () => void;
  }[];
};

export const MenuDesktop = ({ userInfo, handleLogout, menuItems }: MenuDesktopProps) => {
  const [userMenuAnchor, setUserMenuAnchor] = useState<null | HTMLElement>(null);
  const navigate = useNavigate();

  return (
    <AppBar position="sticky">
      <Toolbar>
        <Box sx={{ display: "flex", alignItems: "center" }}>
          <ThemeButton />
          <Typography
            variant="h6"
            component="div"
            sx={{ cursor: "pointer", marginLeft: 1 }}
            onClick={() => {
              navigate(AppRoute.Root);
            }}
          >
            NovelRead
          </Typography>
        </Box>

        <Box
          sx={{
            position: "absolute",
            left: "50%",
            transform: "translateX(-50%)",
            display: "flex",
            gap: 2,
          }}
        >
          {menuItems.map((item, index) => (
            <Button color="inherit" key={index} startIcon={item.icon} onClick={item.onClick}>
              {item.text}
            </Button>
          ))}
        </Box>

        <Box sx={{ marginLeft: "auto" }}>
          {userInfo ? (
            <>
              <Button
                variant="text"
                color="inherit"
                startIcon={<Avatar src={userInfo.avatar} sx={{ width: 32, height: 32 }} />}
                onClick={(e) => {
                  setUserMenuAnchor(e.currentTarget);
                }}
              >
                {userInfo.login}
              </Button>
              <Menu
                anchorEl={userMenuAnchor}
                open={Boolean(userMenuAnchor)}
                onClose={() => {
                  setUserMenuAnchor(null);
                }}
              >
                <MenuItem
                  onClick={() => {
                    navigate(`/profile/${userInfo.id}`);
                    setUserMenuAnchor(null);
                  }}
                >
                  Go to Profile
                </MenuItem>
                <MenuItem onClick={handleLogout}>Logout</MenuItem>
              </Menu>
            </>
          ) : (
            <Button
              color="inherit"
              onClick={() => {
                navigate(AppRoute.Login);
              }}
            >
              Login
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
};
