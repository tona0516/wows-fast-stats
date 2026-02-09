<script lang="ts">
  import { STATS_COLUMN_INFO } from "@libs/constants";
  import { storedDisplayPref } from "@libs/stores";
  import type { PlayerNameColorType, StatsKey } from "@libs/types";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const PLAYER_NAME_COLORS: Readonly<Map<PlayerNameColorType, string>> =
    new Map<PlayerNameColorType, string>([
      ["shipPR", "艦成績のPR"],
      ["overallPR", "総合成績のPR"],
      ["threatLevel", "戦力評価"],
      ["none", "なし"],
    ]);
  const STATS_KEYS = Object.keys(STATS_COLUMN_INFO) as readonly StatsKey[];

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<table class="table">
  <tbody>
    <tr>
      <td class="font-semibold text-nowrap">プレイヤー名</td>
      <td></td>
      <td>
        <div class="py-1">
          <label class="label text-base-content">
            <div class="flex flex-col text-nowrap">
              <span>クラン国籍の国旗を表示</span>
              <span>(クラン説明から言語検出)</span>
            </div>
            <input
              class="toggle toggle-success"
              type="checkbox"
              bind:checked={$storedDisplayPref.player.enableNationFlag}
              on:change={onChangePref}
            />
          </label>
        </div>
        <div class="py-1 flex items-center gap-2 whitespace-nowrap">
          <div>背景色タイプ</div>
          <div>
            <select
              class="select select-sm"
              bind:value={$storedDisplayPref.player.colorType}
              on:change={onChangePref}
            >
              {#each PLAYER_NAME_COLORS as color}
                <option value={color[0]}>{color[1]}</option>
              {/each}
            </select>
          </div>
        </div>
      </td>
    </tr>

    <tr>
      <td class="font-semibold text-nowrap">艦名</td>
      <td></td>
      <td>
        <div class="py-1 text-nowrap">
          <label class="label text-base-content">
            国旗を表示
            <input
              class="toggle toggle-success"
              type="checkbox"
              bind:checked={$storedDisplayPref.warship.enableNationFlag}
              on:change={onChangePref}
            />
          </label>
        </div>
      </td>
    </tr>

    {#each STATS_KEYS as statsKey}
      {@const info = STATS_COLUMN_INFO[statsKey]}
      <tr>
        <td class="font-semibold text-nowrap">{info.minName}</td>
        <td class="whitespace-break-spaces">{info.discription}</td>
        <td>
          {#if "digit" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].digit === "number"}
            <div class="py-1 flex items-center gap-2 text-nowrap">
              <div>小数部桁数</div>
              <div>
                <select
                  class="select select-sm"
                  bind:value={$storedDisplayPref[statsKey].digit}
                  on:change={onChangePref}
                >
                  {#each [0, 1, 2] as digit}
                    <option value={digit}>{digit}</option>
                  {/each}
                </select>
              </div>
            </div>
          {/if}
        </td>
      </tr>
    {/each}
  </tbody>
</table>
