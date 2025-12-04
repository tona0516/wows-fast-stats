// Note: see style.css
export namespace Theme {
  const LIGHT = "light";
  const DARK = "dark";

  export const isLighter = (): boolean => {
    return getCurrent() === LIGHT;
  };

  export const getCurrent = (): string => {
    return localStorage.getItem("theme") || LIGHT;
  };

  export const getAll = (): string[] => {
    return [LIGHT, DARK];
  };
}
