<script lang="ts">
import { createEventDispatcher } from "svelte";
import { WindowReloadApp } from "../wailsjs/runtime/runtime";
import type { Page } from "./Page";
import type { vo } from "wailsjs/go/models";

const dispatch = createEventDispatcher();

type Func = "reload" | "screenshot";
type NavigationMenu = Page | Func;

export let config: vo.UserConfig;
export let currentPage: Page = "main";
export let battle: vo.Battle;

function onClickMenu(menu: NavigationMenu) {
  switch (menu) {
    case "main":
      currentPage = "main";
      break;
    case "config":
      currentPage = "config";
      break;
    case "appinfo":
      currentPage = "appinfo";
      break;
    case "reload":
      WindowReloadApp();
      break;
    case "screenshot":
      dispatch("onScreenshot");
      break;
    default:
      break;
  }
}
</script>

<nav class="navbar navbar-expand-sm navbar-light bg-light">
  <div class="container-fluid">
    <button
      class="navbar-toggler"
      type="button"
      data-bs-toggle="collapse"
      data-bs-target="#navbarNavAltMarkup"
      aria-controls="navbarNavAltMarkup"
      aria-expanded="false"
      aria-label="Toggle navigation"
      style="font-size: {config?.font_size || 'medium'};"
    >
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
      <div class="navbar-nav">
        <button
          type="button"
          class="btn btn-sm btn-outline-secondary m-1 {currentPage === 'main' &&
            'active'}"
          title="ホーム"
          style="font-size: {config?.font_size || 'medium'};"
          on:click="{() => onClickMenu('main')}"
        >
          <i class="bi bi-house"></i>
          ホーム
        </button>
        <button
          type="button"
          class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
            'config' && 'active'}"
          title="設定"
          style="font-size: {config?.font_size || 'medium'};"
          on:click="{() => onClickMenu('config')}"
        >
          <i class="bi bi-gear"></i>
          設定
        </button>
        <button
          type="button"
          class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
            'appinfo' && 'active'}"
          title="アプリ情報"
          style="font-size: {config?.font_size || 'medium'};"
          on:click="{() => onClickMenu('appinfo')}"
        >
          <i class="bi bi-info-circle"></i>
          アプリ情報
        </button>
        {#if currentPage == "main"}
          <button
            type="button"
            class="btn btn-sm btn-outline-success m-1"
            title="リロード"
            style="font-size: {config?.font_size || 'medium'};"
            on:click="{() => onClickMenu('reload')}"
          >
            <i class="bi bi-arrow-clockwise"></i>
            リロード
          </button>

          <button
            type="button"
            class="btn btn-sm btn-outline-success m-1"
            title="スクリーンショット"
            disabled="{battle === undefined}"
            style="font-size: {config?.font_size || 'medium'};"
            on:click="{() => onClickMenu('screenshot')}"
          >
            <i class="bi bi-camera"></i>
            スクリーンショット
          </button>
        {/if}
      </div>
    </div>
  </div>
</nav>
