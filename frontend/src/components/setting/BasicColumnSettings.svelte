<script lang="ts">
  import { storedDisplayPref } from "@libs/stores";
  import type { PlayerNameColorType } from "@libs/types";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const PLAYER_NAME_COLORS: Readonly<Map<PlayerNameColorType, string>> =
    new Map<PlayerNameColorType, string>([
      ["shipPR", "艦成績のPR"],
      ["overallPR", "総合成績のPR"],
      ["threatLevel", "戦力評価"],
      ["none", "なし"],
    ]);

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div class="bg-base-100 shadow-xl rounded-xl p-4">
  <table class="table">
    <tbody>
      <tr class="hover:bg-base-100">
        <td class="font-bold">プレイヤー名</td>
        <td class="text-right">
          <div>
            <span>クラン国籍の国旗を表示する（クラン説明から言語検出）</span>
            <input
              class="toggle toggle-success"
              type="checkbox"
              bind:checked={$storedDisplayPref.player.enableNationFlag}
              on:change={onChangePref}
            />
          </div>
          <div class="mt-2">
            <span>背景色タイプ</span>
            <select
              class="select select-sm select-bordered"
              bind:value={$storedDisplayPref.player.colorType}
              on:change={onChangePref}
            >
              {#each PLAYER_NAME_COLORS as color}
                <option
                  selected={color[0] === $storedDisplayPref.player.colorType}
                  value={color[0]}>{color[1]}</option
                >
              {/each}
            </select>
          </div>
        </td>
      </tr>
      <tr class="hover:bg-base-100">
        <td class="text-bold">艦名</td>
        <td class="text-right">
          <span>国旗を表示する</span>
          <input
            class="toggle toggle-success"
            type="checkbox"
            bind:checked={$storedDisplayPref.warship.enableNationFlag}
            on:change={onChangePref}
          />
        </td>
      </tr>
    </tbody>
  </table>
</div>
