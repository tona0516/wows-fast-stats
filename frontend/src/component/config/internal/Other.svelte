<script lang="ts">
import { Notifier } from "src/lib/Notifier";
import { storedConfig } from "src/stores";
import { OpenDirectory, UpdateUserConfig } from "wailsjs/go/main/App";

$: inputConfig = $storedConfig;

const openDirectory = async (path: string) => {
  try {
    await OpenDirectory(path);
  } catch (error) {
    Notifier.failure(error);
  }
};

const change = async () => {
  try {
    await UpdateUserConfig(inputConfig);
  } catch (error) {
    inputConfig = $storedConfig;
    Notifier.failure(error);
  }
};
</script>

<div>
  <div>
    <label
      ><input
        class="checkbox"
        type="checkbox"
        bind:checked={inputConfig.save_temp_arena_info}
        on:change={change}
      /> 【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する</label
    >
    <div>
      <!-- svelte-ignore a11y-invalid-attribute -->
      <a href="#" on:click={() => openDirectory("temp_arena_info")}>
        <i class="bi bi-folder">保存フォルダを開く</i>
      </a>
    </div>
  </div>

  <div>
    <label
      ><input
        class="checkbox"
        type="checkbox"
        bind:checked={inputConfig.send_report}
        on:change={change}
      /> アプリ改善のためのデータ送信を許可する</label
    >
    <div>
      <ul>
        <li>アプリ統計情報</li>
        <li>プレイヤー名</li>
        <li>エラーログ</li>
        <li>設定値(config/user.json)</li>
        <li>戦闘情報(tempArenaInfo.json)</li>
      </ul>
    </div>
  </div>

  <div>
    <label
      ><input
        class="checkbox"
        type="checkbox"
        bind:checked={inputConfig.notify_updatable}
        on:change={change}
      /> 新しいバージョンがある場合に通知する</label
    >
  </div>
</div>
