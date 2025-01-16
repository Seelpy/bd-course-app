import { IconButton, styled } from "@mui/material";
import { DarkMode, LightMode } from "@mui/icons-material";
import { useTheme } from "@shared/hooks/useTheme";
import { usePreferencesStore } from "@shared/stores/preferencesStore";
import { Theme } from "@shared/types/theme";
import { useShallow } from "zustand/shallow";

const IconContainer = styled("div")`
  position: relative;
  width: 24px;
  height: 24px;
`;

const IconWrapper = styled("div")<{ $isActive: boolean }>`
  position: absolute;
  transition: all 0.3s ease-in-out;
  transform-origin: bottom right;
  opacity: ${({ $isActive }) => ($isActive ? 1 : 0)};

  &.sun {
    transform: ${({ $isActive }) =>
      $isActive ? "rotate(0deg) translate(0, 0)" : "rotate(-90deg) translate(0, -100%)"};
  }

  &.moon {
    transform: ${({ $isActive }) =>
      $isActive ? "rotate(0deg) translate(0, 0)" : "rotate(-90deg) translate(0, -100%)"};
  }
`;

const AutoIndicator = styled("span")`
  position: absolute;
  top: -7px;
  left: -7px;
  font-size: 12px;
  font-weight: bold;
  color: inherit;
`;

const ThemeButton = () => {
  const { mode } = useTheme();
  const { theme, changeTheme } = usePreferencesStore(
    useShallow((state) => ({
      theme: state.theme,
      changeTheme: state.changeTheme,
    })),
  );

  const handleChange = () => {
    const themes: Theme[] = ["light", "dark", "auto"];
    const currentIndex = themes.indexOf(theme);
    const nextTheme = themes[(currentIndex + 1) % themes.length];
    changeTheme(nextTheme);
  };

  return (
    <IconButton onClick={handleChange} color="inherit" aria-label="toggle theme">
      <IconContainer>
        <IconWrapper $isActive={mode === "light"} className="sun">
          <LightMode sx={{ color: "#FFB800" }} />
        </IconWrapper>
        <IconWrapper $isActive={mode === "dark"} className="moon">
          <DarkMode sx={{ color: "#90CAF9" }} />
        </IconWrapper>
        {theme === "auto" && <AutoIndicator>A</AutoIndicator>}
      </IconContainer>
    </IconButton>
  );
};

export default ThemeButton;
