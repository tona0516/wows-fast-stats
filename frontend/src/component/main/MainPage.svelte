<script lang="ts">
import { format, fromUnixTime } from "date-fns";
import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
import { FetchProxy } from "src/lib/FetchProxy";
import { Notifier } from "src/lib/Notifier";
import { Screenshot } from "src/lib/Screenshot";
import {
  storedBattle,
  storedConfig,
  storedInstallPathError,
  storedSummary,
} from "src/stores";
import { LogInfo } from "wailsjs/go/main/App";
import type { data } from "wailsjs/go/models";
import Menu from "./internal/Menu.svelte";
import Ofuse from "./internal/Ofuse.svelte";
import Summary from "./internal/Summary.svelte";

const MAIN_PAGE_ID = "mainpage";

let menu: Menu | undefined;
let isLoading = false;
let isScreenshotting = false;

// Note: Promiseがガベージコレクションによって解放されてしまうため保持する
let autoScreenshotPromise: Promise<void>;

export const fetchBattle = async () => {
  try {
    isLoading = true;

    const start = new Date().getTime();
    await FetchProxy.getBattle();
    const elapsed = (new Date().getTime() - start) / 1000;

    Notifier.success(`データ取得完了: ${elapsed.toFixed(1)}秒`);

    if ($storedConfig.save_screenshot) {
      autoScreenshotPromise = autoScreenshot();
    }

    LogInfo("fetch success", { "duration(s)": elapsed.toFixed(1) });
  } catch (error) {
    Notifier.failure(error);
  } finally {
    isLoading = false;
  }
};

const manualScreenshot = async () => {
  try {
    const meta = $storedBattle?.meta;
    if (!meta) {
      return;
    }

    isScreenshotting = true;
    const isSuccess = await Screenshot.manual(
      MAIN_PAGE_ID,
      deriveFileName(meta),
    );

    if (isSuccess) {
      Notifier.success("スクリーンショットを保存しました");
    }
  } catch (error) {
    Notifier.failure("スクリーンショットに失敗しました", 10000);
  } finally {
    isScreenshotting = false;
  }
};

const autoScreenshot = async () => {
  try {
    const meta = $storedBattle?.meta;
    if (!meta) {
      return;
    }
    isScreenshotting = true;
    await Screenshot.auto(MAIN_PAGE_ID, deriveFileName(meta));
  } catch (error) {
    Notifier.failure("スクリーンショットに失敗しました", 10000);
  } finally {
    isScreenshotting = false;
  }
};

const deriveFileName = (meta: data.Meta): string => {
  const items = [
    format(fromUnixTime(meta.unixtime), "yyyy-MM-dd-HH-mm-ss"),
    meta.own_ship.replaceAll(" ", "-"),
    meta.arena,
    meta.type,
  ];

  return `${items.join("_")}`;
};
</script>

<!-- Note: Use the same color as that of body.  -->
<div id={MAIN_PAGE_ID}>
  <div class="flex">
    <Menu
      bind:this={menu}
      {isScreenshotting}
      on:ManualScreenshot={() => manualScreenshot()}
    />
  </div>

  <div>
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}
      {@const meta = $storedBattle.meta}
      {@const config = $storedConfig}

      <div class="flex">
        <StatisticsTable
          {teams}
          {config}
          on:EditAlertPlayer
          on:RemoveAlertPlayer
        />
      </div>

      <div class="flex">
        <BattleMeta {meta} />
      </div>

      {#if $storedSummary}
        {@const summary = $storedSummary}
        <div class="flex">
          <Summary {summary} />
        </div>
      {/if}
    {:else}
      <p>
        {#if $storedInstallPathError}
          設定画面から初期設定を行ってください。
        {:else}
          戦闘中ではありません。開始時に自動的にリロードします。
        {/if}
      </p>
    {/if}
  </div>

  <div>
    <Ofuse />
  </div>

  {#if isLoading}
    <div class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-20">
      <span class="loading loading-spinner"></span>
    </div>
  {/if}
</div>
