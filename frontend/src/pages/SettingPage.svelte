<script lang="ts">
  import GameClientPathSection from "@components/setting/GameClientPathSection.svelte";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import GlobalDisplaySection from "@components/setting/GlobalDisplaySection.svelte";
  import ModalCommon from "@components/modals/ModalCommon.svelte";
  import { DEFAULT_DISPLAY_PREF } from "@libs/DisplayPref";
  import { showToast, storedDisplayPref } from "@libs/stores";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import Section from "@components/commons/Section.svelte";
  import SubSection from "@components/commons/SubSection.svelte";
  import ColumnSettingTable from "@components/setting/ColumnSettingTable.svelte";
  import StatsColumnOrderSettings from "@components/setting/StatsColumnOrderSettings.svelte";

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

<div class="p-4 mx-auto max-w-3xl flex flex-col gap-4">
  <GameClientPathSection />
  <GlobalDisplaySection />

  {#key resetToken}
    <Section title="カラム別表示設定">
      <div class="flex flex-col gap-4">
        <SubSection title="詳細設定">
          <ColumnSettingTable />
        </SubSection>

        <StatsColumnOrderSettings />
      </div>
    </Section>
  {/key}

  <div class="flex justify-end">
    <button class="btn btn-outline" on:click={openConfirm} type="button">
      表示設定をデフォルトに戻す
    </button>
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
