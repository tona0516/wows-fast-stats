<script lang="ts">
  import { storedConfig, storedRequiredConfigError } from "src/stores";
  import { UserConfig, ApplyUserConfig } from "wailsjs/go/main/App";
  import Required from "./internal/Required.svelte";
  import Other from "./internal/Other.svelte";
  import TeamSummary from "./internal/TeamSummary.svelte";
  import Display from "./internal/Display.svelte";
  import AlertPlayer from "./internal/AlertPlayer.svelte";
  import AppInfo from "./internal/AppInfo.svelte";
  import Logging from "./internal/Logging.svelte";
  import { createEventDispatcher } from "svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkTab from "src/component/common/uikit/UkTab.svelte";
  import clone from "clone";

  const dispatch = createEventDispatcher();
  const CONFIG_MENU_ID = "config-menu-id";

  let inputConfig = clone($storedConfig);
  storedConfig.subscribe((it) => {
    inputConfig = clone(it);
  });

  const silentApply = async () => {
    try {
      await ApplyUserConfig(inputConfig);
      const latest = await UserConfig();
      storedConfig.set(latest);
    } catch (error) {
      inputConfig = clone($storedConfig);
      dispatch("Failure", { message: error });
    }
  };
</script>

<div class="uk-padding-small uk-grid">
  <div class="uk-width-auto@m">
    <UkTab clazz="uk-tab-left" id={CONFIG_MENU_ID}>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          必須設定
          {#if !$storedRequiredConfigError.valid}
            <span class="uk-text-warning uk-text-small">
              <UkIcon name="warning" />
            </span>
          {/if}
        </a>
      </li>
      {#each ["表示設定", "チームサマリー設定", "プレイヤーリスト設定", "その他設定", "ログ", "アプリ情報"] as menu}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li><a href="#">{menu}</a></li>
      {/each}
    </UkTab>
  </div>
  <div class="uk-width-expand@m">
    <ul id={CONFIG_MENU_ID} class="uk-switcher">
      <li>
        <Required
          {inputConfig}
          on:UpdateSuccess={() => dispatch("UpdateRequired")}
          on:Failure
        />
      </li>
      <li>
        <Display {inputConfig} on:UpdateSuccess on:Change={silentApply} />
      </li>
      <li>
        <TeamSummary {inputConfig} on:UpdateSuccess on:Failure />
      </li>
      <li>
        <AlertPlayer
          on:AddAlertPlayer
          on:EditAlertPlayer
          on:RemoveAlertPlayer
        />
      </li>
      <li>
        <Other {inputConfig} on:Change={silentApply} on:Failure />
      </li>
      <li>
        <Logging />
      </li>
      <li>
        <AppInfo />
      </li>
    </ul>
  </div>
</div>
