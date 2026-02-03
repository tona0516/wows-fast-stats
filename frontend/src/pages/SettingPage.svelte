<script lang="ts">
  import GameClientPathSection from "@components/setting/GameClientPathSection.svelte";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import GlobalDisplaySection from "@components/setting/GlobalDisplaySection.svelte";
  import ColumnDisplaySection from "@components/setting/ColumnDisplaySection.svelte";
  import ModalCommon from "@components/modals/ModalCommon.svelte";
  import { DEFAULT_DISPLAY_PREF } from "@libs/DisplayPref";
  import { showToast, storedDisplayPref } from "@libs/stores";
  import { SaveDisplayPref } from "@wails/go/main/App";

  onMount(async () => {
    themeChange(false);
  });

  let isConfirmOpen = false;
  let resetToken = 0;

  const openConfirm = () => {
    isConfirmOpen = true;
  };

  const closeConfirm = () => {
    isConfirmOpen = false;
  };

  const onResetDisplayPref = async () => {
    const nextPref = JSON.parse(JSON.stringify(DEFAULT_DISPLAY_PREF));
    storedDisplayPref.set(nextPref);
    await SaveDisplayPref(JSON.stringify(nextPref, null, 2));
    showToast("表示設定をデフォルトに戻しました");
    resetToken += 1;
    closeConfirm();
  };
</script>

<div class="min-h-full w-full bg-base-200/40 p-4">
  <div class="mx-auto flex w-full max-w-3xl flex-col gap-4">
    <GameClientPathSection />
    <GlobalDisplaySection />
    <ColumnDisplaySection {resetToken} />
    <div class="flex justify-end">
      <button class="btn btn-outline" on:click={openConfirm} type="button">
        表示設定をデフォルトに戻す
      </button>
    </div>
  </div>
</div>

{#if isConfirmOpen}
  <ModalCommon close={closeConfirm}>
    <div class="space-y-4">
      <div class="text-lg font-bold">確認</div>
      <p class="text-sm text-base-content/80">
        表示設定をデフォルトに戻します。よろしいですか？
      </p>
      <div class="flex justify-end gap-2">
        <button class="btn btn-ghost" type="button" on:click={closeConfirm}>
          キャンセル
        </button>
        <button
          class="btn btn-error"
          type="button"
          on:click={onResetDisplayPref}
        >
          戻す
        </button>
      </div>
    </div>
  </ModalCommon>
{/if}
