<script lang="ts">
import { Notifier } from "src/lib/Notifier";
import { storedConfig, storedInstallPathError } from "src/stores";
import {
  SelectDirectory,
  StartWatching,
  UpdateInstallPath,
} from "wailsjs/go/main/App";

let isLoading = false;

$: inputConfig = $storedConfig;

const clickSelectDirectory = async () => {
  try {
    const path = await SelectDirectory();
    if (!path) return;
    inputConfig.install_path = path;
  } catch (error) {
    Notifier.failure(error);
  }
};

const clickApply = async () => {
  try {
    isLoading = true;
    await UpdateInstallPath(inputConfig.install_path);
    storedInstallPathError.set("");
    Notifier.success("設定を更新しました");
    StartWatching();
  } catch (error) {
    storedInstallPathError.set(error as string);
  } finally {
    isLoading = false;
  }
};
</script>

<div>
  <div class="flex">
    <input
      class="input"
      type="text"
      placeholder="World of Warshipsインストールフォルダ"
      bind:value={inputConfig.install_path}
    />
    <button class="btn btn-neutral" on:click={clickSelectDirectory}
      >フォルダ選択</button
    >
  </div>
  <span>ゲームクライアントの実行ファイルがあるフォルダを選択してください。</span
  >
  {#if $storedInstallPathError}
    <div>
      <i class="bi bi-warning">{$storedInstallPathError}</i>
    </div>
  {/if}
</div>

<div>
  <div class="flex">
    <button class="btn btn-primary" disabled={isLoading} on:click={clickApply}>
      {#if isLoading}
        <span class="loading loading-spinner"></span>
      {:else}
        保存
      {/if}
    </button>
  </div>
</div>
