<script lang="ts">
  import clone from "clone";
  import {
    storedUserConfig,
    storedDefaultUserConfig,
    storedLogs,
  } from "src/stores";
  import { UserConfig, ApplyUserConfig } from "wailsjs/go/main/App";
  import Required from "./tab/Required.svelte";
  import Other from "./tab/Other.svelte";
  import TeamSummary from "./tab/TeamSummary.svelte";
  import Display from "./tab/Display.svelte";
  import AlertPlayer from "./tab/AlertPlayer.svelte";
  import AppInfo from "./tab/AppInfo.svelte";
  import Logging from "./tab/Logging.svelte";
  import { createEventDispatcher } from "svelte";

  export let isSatisfiedRequired: boolean;

  const dispatch = createEventDispatcher();

  let inputUserConfig = clone($storedUserConfig);
  storedUserConfig.subscribe((it) => {
    inputUserConfig = clone(it);
  });

  const silentApply = async () => {
    try {
      await ApplyUserConfig(inputUserConfig);
      const latest = await UserConfig();
      storedUserConfig.set(latest);
    } catch (error) {
      inputUserConfig = clone($storedUserConfig);
      dispatch("Failure", { message: error });
    }
  };
</script>

<div class="uk-padding-small">
  <div uk-grid>
    <div class="uk-width-auto@m">
      <ul
        class="uk-tab-left"
        uk-tab="connect: #component-tab-left; animation: uk-animation-fade"
      >
        <li>
          <a href="#">
            必須設定
            {#if !isSatisfiedRequired}
              <span
                class="uk-text-warning uk-text-small"
                uk-icon="icon: warning"
              ></span>
            {/if}
          </a>
        </li>
        {#each ["表示設定", "チームサマリー設定", "プレイヤーリスト設定", "その他設定", "ログ", "アプリ情報"] as menu}
          <!-- svelte-ignore a11y-invalid-attribute -->
          <li><a href="#">{menu}</a></li>
        {/each}
      </ul>
    </div>
    <div class="uk-width-expand@m">
      <ul id="component-tab-left" class="uk-switcher">
        <li>
          <Required {inputUserConfig} on:UpdateSuccess on:Failure />
        </li>
        <li>
          <Display
            {inputUserConfig}
            defaultUserConfig={$storedDefaultUserConfig}
            on:UpdateSuccess
            on:Change={silentApply}
          />
        </li>
        <li>
          <TeamSummary {inputUserConfig} on:UpdateSuccess on:Failure />
        </li>
        <li>
          <AlertPlayer
            on:AddAlertPlayer
            on:EditAlertPlayer
            on:RemoveAlertPlayer
          />
        </li>
        <li>
          <Other {inputUserConfig} on:Change={silentApply} on:Failure />
        </li>
        <li>
          <Logging logs={$storedLogs} />
        </li>
        <li>
          <AppInfo />
        </li>
      </ul>
    </div>
  </div>
</div>
