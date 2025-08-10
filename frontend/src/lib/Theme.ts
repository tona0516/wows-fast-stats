export namespace Theme {
  const lighterThemes = ["light", "nord", "garden"];
  const darkerThemes = ["dark", "night", "black"];

  export const getAll = (): string[] => {
    return [...lighterThemes, ...darkerThemes];
  };

  export const isLighter = (): boolean => {
    const current = getCurrent();
    return lighterThemes.includes(current);
  };

  const getCurrent = (): string => {
    return localStorage.getItem("theme") || "light";
  };
}
