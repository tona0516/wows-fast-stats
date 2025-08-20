// Note: see style.css
const LIGHTER_THEMES: string[] = ["light", "nord", "garden"] as const;
const DARKER_THEMES: string[] = ["dark", "night", "black"] as const;
const THEMES: string[] = [...LIGHTER_THEMES, ...DARKER_THEMES] as const;

export namespace Theme {
  export const isLighter = (): boolean => {
    const current = getCurrent();
    return LIGHTER_THEMES.includes(current);
  };

  export const getCurrent = (): string => {
    return localStorage.getItem("theme") || "light";
  };

  export const getAll = (): string[] => {
    return THEMES;
  };
}
