<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
  import ConfirmModal from "src/component/modal/ConfirmModal.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import UIkit from "uikit";
  import {
    ApplyUserConfig,
    DefaultUserConfig,
    UserConfig,
  } from "wailsjs/go/main/App";
  import { domain } from "wailsjs/go/models";
  import clone from "clone";
  import { ColumnProvider } from "src/lib/column/ColumnProvider";
  import { SAMPLE_TEAM } from "src/lib/rating/RatingColorFactory";
  import { ModalElementID } from "src/component/modal/ModalElementID";

  export let inputUserConfig: domain.UserConfig;

  const dispatch = createEventDispatcher();

  const reset = async () => {
    try {
      const defaultConfig = await DefaultUserConfig();

      inputUserConfig.font_size = defaultConfig.font_size;
      inputUserConfig.displays = defaultConfig.displays;
      inputUserConfig.custom_color = defaultConfig.custom_color;
      inputUserConfig.custom_digit = defaultConfig.custom_digit;

      await ApplyUserConfig(inputUserConfig);

      const latest = await UserConfig();
      storedUserConfig.set(latest);
      dispatch("UpdateSuccess");
    } catch (error) {
      inputUserConfig = clone($storedUserConfig);
      dispatch("Failure", { message: error });
    }
  };

  $: displays = ColumnProvider.displayItems();
</script>

<ConfirmModal message="表示設定をリセットしますか？" on:Confirmed={reset} />

<div class="uk-padding-small">
  <div>UIサイズ</div>
  <select
    class="uk-select uk-form-width-small"
    bind:value={inputUserConfig.font_size}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.FONT_SIZES.toArray() as fs}
      <option selected={fs.key === $storedUserConfig.font_size} value={fs.key}
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
                bind:checked={inputUserConfig.displays.ship[display.shipKey]}
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
                bind:checked={inputUserConfig.displays.overall[
                  display.overallKey
                ]}
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
                bind:value={inputUserConfig.custom_digit[display.digitKey]}
                on:change={() => dispatch("Change")}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit ===
                      inputUserConfig.custom_digit[display.digitKey]}
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
              bind:value={inputUserConfig.custom_color.skill.text[sl.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputUserConfig.custom_color.skill.background[sl.key]}
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
              bind:value={inputUserConfig.custom_color.tier.own[tg.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputUserConfig.custom_color.tier.other[tg.key]}
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
              bind:value={inputUserConfig.custom_color.ship_type.own[st.key]}
              on:input={() => dispatch("Change")}
            />
          </td>

          <td>
            <input
              class="uk-input"
              type="color"
              bind:value={inputUserConfig.custom_color.ship_type.other[st.key]}
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
    bind:value={inputUserConfig.custom_color.player_name}
    on:change={() => dispatch("Change")}
  >
    {#each DispName.PLAYER_NAME_COLORS.toArray() as pnc}
      <option
        selected={pnc.key === $storedUserConfig.custom_color.player_name}
        value={pnc.key}>{pnc.value}</option
      >
    {/each}
  </select>
</div>

<div class="uk-padding-small">
  <div>プレビュー</div>
  <StatisticsTable teams={[SAMPLE_TEAM]} userConfig={inputUserConfig} />
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