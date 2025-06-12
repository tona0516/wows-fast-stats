<script lang="ts">
import { storedConfig } from "src/stores";
import AppInfo from "./internal/AppInfo.svelte";
import ColorDescription from "./internal/DispColor.svelte";
import Logging from "./internal/Logging.svelte";

const tabName = "tab_info";
</script>

<div class="tabs tabs-lift">
  <input
    type="radio"
    name={tabName}
    class="tab"
    aria-label="表示色について"
    checked={true}
  />
  <div class="tab-content bg-base-100 border-base-300 p-6">
    <ColorDescription config={$storedConfig} />
  </div>

  <input type="radio" name={tabName} class="tab" aria-label="ログ" />
  <div class="tab-content bg-base-100 border-base-300 p-6">
    <Logging />
  </div>

  <input type="radio" name={tabName} class="tab" aria-label="アプリ情報" />
  <div class="tab-content bg-base-100 border-base-300 p-6">
    <AppInfo />
  </div>
</div>
