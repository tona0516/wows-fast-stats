<script lang="ts">
  import { storedConfig } from "src/stores";
  import AppInfo from "./internal/AppInfo.svelte";
  import UkTab from "src/component/common/uikit/UkTab.svelte";
  import ColorDescription from "./internal/DispColor.svelte";

  const INFO_MENU_ID = "info-menu-id";
</script>

<div class="uk-padding-small uk-grid">
  <div class="uk-width-auto@m">
    <UkTab clazz="uk-tab-left" id={INFO_MENU_ID}>
      {#each ["表示色について", "アプリ情報"] as menu}
        <!-- svelte-ignore a11y-invalid-attribute -->
        <li><a href="#">{menu}</a></li>
      {/each}
    </UkTab>
  </div>
  <div class="uk-width-expand@m">
    <ul id={INFO_MENU_ID} class="uk-switcher">
      <li>
        <ColorDescription config={$storedConfig} />
      </li>
      <li>
        <AppInfo />
      </li>
    </ul>
  </div>
</div>
