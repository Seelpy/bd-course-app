import { useState } from "react";
import {
  AppBar,
  Toolbar,
  IconButton,
  Typography,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Divider,
  ListItemIcon,
} from "@mui/material";
import { Menu as MenuIcon, Person, Login, Logout, Close } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";
import ThemeButton from "../ThemeButton/ThemeButton";
import { User } from "@shared/types/user";
import { AppRoute } from "@shared/constants/routes";

type MenuMobileProps = {
  userInfo: User | null;
  handleLogout: () => void;
  menuItems: {
    text: string;
    icon: React.ReactNode;
    onClick: () => void;
  }[];
};

export const MenuMobile = ({ userInfo, handleLogout, menuItems }: MenuMobileProps) => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const navigate = useNavigate();

  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <IconButton
            color="inherit"
            edge="start"
            onClick={() => {
              setDrawerOpen(true);
            }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            NovelRead
          </Typography>
        </Toolbar>
      </AppBar>

      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={() => {
          setDrawerOpen(false);
        }}
      >
        <List dense={true} sx={{ width: 250 }}>
          <ListItem>
            <ThemeButton />
            <Typography variant="h6" component="div" sx={{ flexGrow: 1, marginLeft: 1 }}>
              NovelRead
            </Typography>
            <IconButton
              sx={{ marginLeft: "auto" }}
              onClick={() => {
                setDrawerOpen(false);
              }}
            >
              <Close />
            </IconButton>
          </ListItem>

          {userInfo ? (
            <ListItem>
              <ListItemButton
                onClick={() => {
                  navigate("/profile");
                  setDrawerOpen(false);
                }}
              >
                <ListItemIcon>
                  <Person />
                </ListItemIcon>
                <ListItemText primary={userInfo.login} />
              </ListItemButton>
            </ListItem>
          ) : (
            <ListItem>
              <ListItemButton
                onClick={() => {
                  navigate(AppRoute.Login);
                  setDrawerOpen(false);
                }}
              >
                <ListItemIcon>
                  <Login />
                </ListItemIcon>
                <ListItemText primary="Login" />
              </ListItemButton>
            </ListItem>
          )}

          {menuItems.map((item, index) => (
            <ListItem key={index}>
              <ListItemButton
                onClick={() => {
                  item.onClick();
                  setDrawerOpen(false);
                }}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.text} />
              </ListItemButton>
            </ListItem>
          ))}

          {userInfo && (
            <>
              <Divider />
              <ListItem>
                <ListItemButton onClick={handleLogout}>
                  <ListItemIcon>
                    <Logout />
                  </ListItemIcon>
                  <ListItemText primary="Logout" />
                </ListItemButton>
              </ListItem>
            </>
          )}
        </List>
      </Drawer>
    </>
  );
};
