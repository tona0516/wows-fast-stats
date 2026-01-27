import App from "src/App.svelte";
import "src/style.css";

const app = new App({
  // biome-ignore lint/style/noNonNullAssertion: HTML guaranteed
  target: document.getElementById("app")!,
});

export default app;
