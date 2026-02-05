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
  <div class="flex items-center">
    <span class="text-xl font-bold">ゲームクライアントパス設定</span>
    <span class="ml-2 badge badge-outline badge-error whitespace-nowrap"
      >必須</span
    >
  </div>
  <p class="text-sm text-gray-500 mt-2">
    WorldOfWarships.exeが存在するフォルダを選択してください
  </p>
  {#if $storedGameClientPath}
    <div class="mb-2 flex items-center justify-between gap-2">
      <div
        class="rounded-xl border border-base-300 bg-base-200 shadow-sm p-2 flex-1 overflow-x-auto"
      >
        {$storedGameClientPath}
      </div>
      <button class="btn btn-primary shrink-0" on:click={onClickSelect}>
        フォルダ選択
      </button>
    </div>
  {/if}
  {#if $storedGameClientPathError}
    <div class="alert alert-error mt-2">
      <span>{$storedGameClientPathError}</span>
    </div>
  {/if}
  {#if !$storedGameClientPath}
    <button class="btn btn-primary" on:click={onClickSelect}>
      フォルダ選択
    </button>
  {/if}
</div>
