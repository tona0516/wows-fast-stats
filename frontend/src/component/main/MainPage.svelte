<script lang="ts">
  import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
  import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
  import {
    storedBattle,
    storedConfig,
    storedInstallPathError,
    storedSummary,
  } from "src/stores";
  import Menu from "./internal/Menu.svelte";
  import Summary from "./internal/Summary.svelte";
  import Ofuse from "./internal/Ofuse.svelte";
  import UkSpinner from "../common/uikit/UkSpinner.svelte";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";
  import { LogInfo } from "wailsjs/go/main/App";
  import { data } from "wailsjs/go/models";
  import { format, fromUnixTime } from "date-fns";
  import { Screenshot } from "src/lib/Screenshot";

  const MAIN_PAGE_ID = "mainpage";

  let menu: Menu | undefined;
  let isLoading = false;
  let isScreenshotting = false;

  // Note: Promiseがガベージコレクションによって解放されてしまうため保持する
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
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
      isScreenshotting = true;
      const isSuccess = await Screenshot.manual(
        MAIN_PAGE_ID,
        deriveFileName($storedBattle!.meta),
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
      isScreenshotting = true;
      await Screenshot.auto(MAIN_PAGE_ID, deriveFileName($storedBattle!.meta));
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
<div
  id={MAIN_PAGE_ID}
  class="uk-padding-small uk-light uk-background-secondary"
>
  <div class="uk-margin-small uk-flex uk-flex-center">
    <Menu
      bind:this={menu}
      {isScreenshotting}
      on:ManualScreenshot={() => manualScreenshot()}
    />
  </div>

  <div class="uk-margin-small">
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}
      {@const meta = $storedBattle.meta}
      {@const config = $storedConfig}

      <div class="uk-flex uk-flex-center">
        <StatisticsTable
          {teams}
          {config}
          on:EditAlertPlayer
          on:RemoveAlertPlayer
        />
      </div>

      <div class="uk-flex uk-flex-center">
        <BattleMeta {meta} />
      </div>

      {#if $storedSummary}
        {@const summary = $storedSummary}
        <div class="uk-flex uk-flex-center">
          <Summary {summary} />
        </div>
      {/if}
    {:else}
      <p class="uk-text-center">
        {#if $storedInstallPathError}
          戦闘中ではありません。開始時に自動的にリロードします。
        {:else}
          設定画面から初期設定を行ってください。
        {/if}
      </p>
    {/if}
  </div>

  <div class="uk-margin-small">
    <Ofuse />
  </div>

  {#if isLoading}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</div>
