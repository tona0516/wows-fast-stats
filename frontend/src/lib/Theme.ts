export namespace Theme {
  const lighterThemes = ["light", "nord", "garden"];
  const darkerThemes = ["dark", "night", "black"];

  export const getAll = (): string[] => {
    return [...lighterThemes, ...darkerThemes];
  };

  export const isLighter = (): boolean => {
    const current = localStorage.getItem("theme") || "light";
    return lighterThemes.includes(current);
  };
}
