export namespace Theme {
  const LIGHT = "light";
  const DARK = "dark";

  export const isLight = (): boolean => {
    const current = localStorage.getItem("theme") || LIGHT;
    return current === LIGHT;
  };

  export const getAll = (): string[] => {
    return [LIGHT, DARK];
  };
}
