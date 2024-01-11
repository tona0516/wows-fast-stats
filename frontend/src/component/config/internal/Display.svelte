<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
  import ConfirmModal from "src/component/modal/ConfirmModal.svelte";
  import { storedConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import UIkit from "uikit";
  import { ApplyUserConfig, DefaultUserConfig } from "wailsjs/go/main/App";
  import { model } from "wailsjs/go/models";
  import clone from "clone";
  import { SAMPLE_TEAM } from "src/lib/rating/RatingColorFactory";
  import { ModalElementID } from "src/component/modal/ModalElementID";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";

  export let inputConfig: model.UserConfig;

  const dispatch = createEventDispatcher();

  const reset = async () => {
    try {
      const defaultConfig = await DefaultUserConfig();

      inputConfig.font_size = defaultConfig.font_size;
      inputConfig.displays = defaultConfig.displays;
      inputConfig.custom_color = defaultConfig.custom_color;
      inputConfig.custom_digit = defaultConfig.custom_digit;

      await ApplyUserConfig(inputConfig);
      await FetchProxy.getConfig();

      Notifier.success("設定を更新しました");
    } catch (error) {
      inputConfig = clone($storedConfig);
      Notifier.failure(error);
    }
  };

  $: displays = ColumnProvider.getDisplayableColumns();
</script>

<ConfirmModal message="表示設定をリセットしますか？" on:Confirmed={reset} />

<div class="uk-padding-small">
  <div>UIサイズ</div>
  <select
    class="uk-select uk-form-width-small"
    bind:value={inputConfig.font_size}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.FONT_SIZES.toArray() as fs}
      <option selected={fs.key === $storedConfig.font_size} value={fs.key}
        >{fs.value}</option
      >
    {/each}
  </select>
</div>

<div class="uk-padding-small">
  <div>表示項目</div>
  <table
    class="uk-table uk-table-shrink uk-table-divider uk-table-small uk-table-middle uk-text-nowrap"
  >
    <thead>
      <tr>
        {#each ["項目", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
          <th class="uk-text-center">{columns}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each displays as display}
        <tr>
          <td>
            {display.name}
          </td>

          {#if display.shipKey}
            <td class="uk-text-center">
              <input
                class="uk-checkbox"
                type="checkbox"
                bind:checked={inputConfig.displays.ship[display.shipKey]}
                on:change={() => dispatch("Change")}
              />
            </td>
          {:else}
            <td></td>
          {/if}

          {#if display.overallKey}
            <td class="uk-text-center">
              <input
                class="uk-checkbox"
                type="checkbox"
                bind:checked={inputConfig.displays.overall[display.overallKey]}
                on:change={() => dispatch("Change")}
              />
            </td>
          {:else}
            <td></td>
          {/if}

          {#if display.digitKey}
            <td class="uk-text-center">
              <select
                class="uk-select uk-form-small uk-form-width-xsmall"
                bind:value={inputConfig.custom_digit[display.digitKey]}
                on:change={() => dispatch("Change")}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      inputConfig.custom_digit[display.digitKey]}
                    value={digit}>{digit}</option
                  >
                {/each}
              </select>
            </td>
          {:else}
            <td></td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<div class="uk-padding-small">
  <div>各種カラー</div>
  <table
    class="uk-table uk-width-medium uk-table-divider uk-table-middle uk-text-nowrap"
  >
    <thead>
      <tr>
        {#each ["スキル", "文字色", "背景色"] as column}
          <th class="uk-text-center">{column}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each DispName.SKILL_LEVELS.toArray() as sl}
        <tr>
          <td>{sl.value}</td>
          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.skill.text[sl.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.skill.background[sl.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>

  <table
    class="uk-table uk-width-medium uk-table-divider uk-table-middle uk-text-nowrap"
  >
    <thead>
      <tr>
        {#each ["Tier", "使用艦", "非使用艦"] as column}
          <th class="uk-text-center">{column}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each DispName.TIER_GROUPS.toArray() as tg}
        <tr>
          <td>{tg.value}</td>
          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.tier.own[tg.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.tier.other[tg.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>

  <table
    class="uk-table uk-width-medium uk-table-divider uk-table-middle uk-text-nowrap"
  >
    <thead>
      <tr>
        {#each ["艦種", "使用艦", "非使用艦"] as column}
          <th class="uk-text-center">{column}</th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each DispName.SHIP_TYPES.toArray() as st}
        <tr>
          <td>{st.value}</td>
          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.ship_type.own[st.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputConfig.custom_color.ship_type.other[st.key]}
              on:input={() => dispatch("Change")}
            />
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<div class="uk-padding-small">
  <div>プレイヤー名の背景色</div>
  <select
    class="uk-select uk-form-width-medium"
    bind:value={inputConfig.custom_color.player_name}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.PLAYER_NAME_COLORS.toArray() as pnc}
      <option
        selected={pnc.key === $storedConfig.custom_color.player_name}
        value={pnc.key}>{pnc.value}</option
      >
    {/each}
  </select>
</div>

<div class="uk-padding-small">
  <div>プレビュー</div>
  <StatisticsTable teams={[SAMPLE_TEAM]} config={inputConfig} />
</div>

<div class="uk-padding-small">
  <button
    class="uk-button uk-button-danger uk-text-nowrap"
    on:click={() => {
      const elem = document.getElementById(ModalElementID.CONFIRM);
      if (elem) {
        UIkit.modal(elem).show();
      }
    }}>リセット</button
  >
</div>
