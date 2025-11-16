import { createContext, useContext, useEffect, useState } from 'react';

type Theme = 'light' | 'dark' | 'system';

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  defaultAppearance?: 'one';
  storageKey?: string;
};

type ThemeProviderState = {
  theme: Theme;
  appearance: 'one';
  setTheme: (theme: Theme) => void;
  setAppearance: (appearance: 'one') => void;
};

const initialState: ThemeProviderState = {
  theme: 'system',
  appearance: 'one',
  setTheme: () => null,
  setAppearance: () => null,
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = 'system',
  defaultAppearance = 'one',
  storageKey = 'one-theme',
  ...props
}: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme
  );
  const [appearance, setAppearance] = useState<'one'>(
    () => (localStorage.getItem('one-appearance') as 'one') || defaultAppearance
  );

  useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove('light', 'dark');

    if (theme === 'system') {
      const systemTheme = window.matchMedia('(prefers-color-scheme: dark)')
        .matches
        ? `dark`
        : `light`;

      root.classList.add(systemTheme);

      return;
    }

    root.classList.add(`${theme}`);
  }, [theme, appearance]);

  const value = {
    theme,
    setTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setTheme(theme);
    },
    appearance,
    setAppearance: (appearance: 'one') => {
      localStorage.setItem('one-appearance', appearance);
      setAppearance(appearance);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error('useTheme must be used within a ThemeProvider');

  return context;
};
