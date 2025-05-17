<script lang="ts">
import UkIcon from "src/component/common/uikit/UkIcon.svelte";
import UkTab from "src/component/common/uikit/UkTab.svelte";
import { storedInstallPathError } from "src/stores";
import AlertPlayer from "./internal/AlertPlayer.svelte";
import Display from "./internal/Display.svelte";
import Other from "./internal/Other.svelte";
import Required from "./internal/Required.svelte";
import TeamSummary from "./internal/TeamSummary.svelte";

const CONFIG_MENU_ID = "config-menu-id";
</script>

<div class="uk-padding-small uk-grid">
  <div class="uk-width-auto@m">
    <UkTab clazz="uk-tab-left" id={CONFIG_MENU_ID}>
      <li>
        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#">
          必須設定
          {#if $storedInstallPathError}
            <span class="uk-text-warning uk-text-small">
              <UkIcon name="warning" />
            </span>
          {/if}
        </a>
      </li>
      {#each ["表示設定", "チームサマリー設定", "プレイヤーリスト設定", "その他設定"] as menu}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li><a href="#">{menu}</a></li>
      {/each}
    </UkTab>
  </div>
  <div class="uk-width-expand@m">
    <ul id={CONFIG_MENU_ID} class="uk-switcher">
      <li>
        <Required />
      </li>
      <li>
        <Display />
      </li>
      <li>
        <TeamSummary />
      </li>
      <li>
        <AlertPlayer
          on:AddAlertPlayer
          on:EditAlertPlayer
          on:RemoveAlertPlayer
        />
      </li>
      <li>
        <Other />
      </li>
    </ul>
  </div>
</div>
