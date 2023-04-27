<script lang="ts">
  import {
    UserConfig,
    TempArenaInfoHash,
    Battle,
    SaveScreenshot,
  } from "../wailsjs/go/main/App.js";
  import type { vo } from "wailsjs/go/models.js";
  import Notification from "./Notification.svelte";
  import ConfigPage from "./ConfigPage.svelte";
  import MainPage from "./MainPage.svelte";
  import domtoimage from "dom-to-image";
  import AppInfo from "./AppInfo.svelte";
  import Tab, { Icon, Label } from "@smui/tab";
  import TabBar from "@smui/tab-bar";
  import LinearProgress from "@smui/linear-progress";

  type ScreenshotType = "auto" | "manual";
  type Page = "main" | "config" | "appinfo";
  let tabs = [
    {
      icon: "home",
      label: "ホーム",
      on_click: () => (currentPage = "main"),
    },
    {
      icon: "settings",
      label: "設定",
      on_click: () => (currentPage = "config"),
    },
    {
      icon: "info",
      label: "アプリ情報",
      on_click: () => (currentPage = "appinfo"),
    },
  ];
  let active = tabs[0];

  let currentPage: Page = "main";

  let loadState: LoadState;
  let latestHash: string;
  let battle: vo.Battle;
  let config: vo.UserConfig;

  let notification: Notification;

  setInterval(looper, 1000);

  function saveScreenshot(type: ScreenshotType) {
    domtoimage
      .toPng(document.getElementById("mainpage"))
      .then((dataUrl) => {
        const date = battle.meta.date.replaceAll(":", "-").replaceAll(" ", "-");
        const ownShip = battle.meta.own_ship.replaceAll(" ", "-");
        const filename = `${date}_${ownShip}_${battle.meta.arena}_${battle.meta.type}.png`;
        const base64Data = dataUrl.split(",")[1];
        if (type === "auto") {
          return SaveScreenshot(filename, base64Data, false);
        }
        if (type === "manual") {
          return SaveScreenshot(filename, base64Data, true);
        }
      })
      .then(() => {
        if (type === "manual") {
          notification.showToast(
            "スクリーンショットを保存しました。",
            "success"
          );
        }
      })
      .catch((error) => {
        notification.showToast(error, "error");
      });
  }

  async function looper() {
    try {
      config = await UserConfig();
      notification.removeToastWithKey("need_config");
    } catch (error) {
      notification.showToastWithKey(
        "未設定の状態のため開始できません。「設定」から入力してください。",
        "info",
        "need_config"
      );
      return;
    }

    if (loadState === "error" || loadState === "fetching") {
      notification.removeToastWithKey("not_in_battle");
      return;
    }

    let hash: string;
    try {
      hash = await TempArenaInfoHash();
    } catch (error) {
      loadState = "standby";
      notification.showToastWithKey(
        "戦闘中ではありません。開始時に自動的にリロードします。",
        "info",
        "not_in_battle"
      );
      return;
    }

    if (hash === latestHash) {
      return;
    }

    loadState = "fetching";
    try {
      const start = new Date().getTime();
      battle = await Battle();
      latestHash = hash;
      loadState = "standby";

      const elapsed = (new Date().getTime() - start) / 1000;
      notification.showToast(`データ取得完了: ${elapsed}秒`, "success");
    } catch (error) {
      loadState = "error";
      notification.showToast(error, "error");
    }

    if (config.save_screenshot) {
      await saveScreenshot("auto");
    }
  }
</script>

<main>
  <div style="font-size: {config?.font_size || 'medium'};">
    {#if loadState === "fetching"}
      <LinearProgress indeterminate />
    {/if}

    <TabBar {tabs} let:tab bind:active>
      <Tab {tab} on:click={() => tab.on_click()}>
        <Icon class="material-icons">{tab.icon}</Icon>
        <Label>{tab.label}</Label>
      </Tab>
    </TabBar>
    {#if currentPage === "main"}
      <div id="mainpage">
        <MainPage
          bind:loadState
          bind:latestHash
          bind:battle
          bind:config
          on:onClickScreenshot={() => saveScreenshot("manual")}
        />
      </div>
    {/if}

    {#if currentPage === "config"}
      <ConfigPage
        on:SuccessToast={(event) =>
          notification.showToast(event.detail.message, "success")}
        on:ErrorToast={(event) =>
          notification.showToast(event.detail.message, "error")}
      />
    {/if}

    {#if currentPage === "appinfo"}
      <AppInfo />
    {/if}

    <Notification bind:this={notification} />
  </div>
</main>
