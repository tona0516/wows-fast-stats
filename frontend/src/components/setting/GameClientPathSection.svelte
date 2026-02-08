<script lang="ts">
  import Section from "@components/commons/Section.svelte";
  import {
    showToast,
    storedGameClientPath,
    storedGameClientPathError,
  } from "@libs/stores";
  import { SelectGameClientPath } from "@wails/go/main/App";

  const onClickSelect = async () => {
    try {
      const selectedPath = await SelectGameClientPath();
      storedGameClientPath.set(selectedPath);
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

<Section title="ゲームクライアントパス設定" badgeText="必須">
  <div class="flex flex-col gap-2">
    <div class="flex items-center">
      <i class="bi bi-info-circle"></i>
      <p class="ml-1 text-sm">
        WorldOfWarships.exeがあるフォルダを選択してください
      </p>
    </div>
    <div class="flex items-center gap-2">
      <input
        type="text"
        class="rounded-sm border border-base-300 bg-base-200 shadow-sm p-2 flex-1 overflow-x-auto"
        value={$storedGameClientPath}
        disabled
      />
      <button class="btn btn-primary shrink-0" on:click={onClickSelect}>
        フォルダ選択
      </button>
    </div>
    {#if $storedGameClientPathError}
      <div class="alert alert-error alert-outline">
        <span>{$storedGameClientPathError}</span>
      </div>
    {/if}
  </div>
</Section>
