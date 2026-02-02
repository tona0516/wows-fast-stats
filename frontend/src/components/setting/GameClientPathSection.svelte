<script lang="ts">
  import {
    showToast,
    storedGameClientPath,
    storedGameClientPathError,
  } from "@libs/stores";
  import { SelectGameClientPath } from "@wails/go/main/App";

  const onClickSelect = async () => {
    try {
      await SelectGameClientPath();
      storedGameClientPathError.set("");

      showToast("ゲームクライアントパスを設定しました");
    } catch (error) {
      const errorString = error as string;
      if (errorString.includes("C103")) {
        return;
      }

      storedGameClientPathError.set(errorString);
    }
  };
</script>

<div class="bg-base-100 shadow-xl rounded-xl p-4">
  <div class="flex items-center mb-2">
    <span class="text-xl font-bold">ゲームクライアントパス設定</span>
    <span class="ml-2 badge badge-outline badge-error">必須</span>
  </div>
  <p class="text-sm text-gray-500 mb-2">
    WorldOfWarships.exeが存在するフォルダを選択してください
  </p>
  {#if $storedGameClientPath}
    <div class="stats shadow mb-2">
      <div class="stat">
        <div class="stat-title">パス</div>
        <div class="stat-value text-lg break-all">
          {$storedGameClientPath}
        </div>
      </div>
    </div>
  {/if}
  {#if $storedGameClientPathError}
    <div class="alert alert-error mb-2">
      <span>{$storedGameClientPathError}</span>
    </div>
  {/if}
  <div class="flex justify-end">
    <button class="btn btn-primary" on:click={onClickSelect}>
      フォルダ選択
    </button>
  </div>
</div>
